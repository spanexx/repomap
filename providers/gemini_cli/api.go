package gemini_cli

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"llm-adapter/pkg/adapter"
	"llm-adapter/pkg/tools"
)

// Generate sends a prompt and returns the response, handling tool calls internally.
func (p *Provider) Generate(prompt string, attachments []adapter.Attachment) (string, error) {
	ctx := context.Background()

	// Build tool declarations from shared definitions
	toolDefs := tools.GetToolDefinitionsForActiveSession()
	var funcDecls []FunctionDeclaration
	for _, td := range toolDefs {
		fn, ok := td["function"].(map[string]interface{})
		if !ok {
			continue
		}
		name, _ := fn["name"].(string)
		desc, _ := fn["description"].(string)
		params, _ := fn["parameters"].(map[string]interface{})
		funcDecls = append(funcDecls, FunctionDeclaration{
			Name:        name,
			Description: desc,
			Parameters:  params,
		})
	}

	// Build parts with attachments first
	var parts []Part
	for _, att := range attachments {
		if att.Type == adapter.AttachmentTypeText {
			parts = append(parts, Part{Text: fmt.Sprintf("--- File: %s ---\n%s\n--- End File ---", att.Name, att.Data)})
		} else if att.Type == adapter.AttachmentTypeImage {
			parts = append(parts, Part{
				InlineData: &Blob{
					MimeType: att.MimeType,
					Data:     att.Data,
				},
			})
		}
	}
	parts = append(parts, Part{Text: prompt})

	// Initial conversation
	contents := []Content{
		{Role: "user", Parts: parts},
	}

	for {
		result, err := p.callAPI(ctx, contents, funcDecls)
		if err != nil {
			return "", err
		}

		if len(result.Candidates) == 0 {
			return "", fmt.Errorf("no candidates in response")
		}

		candidate := result.Candidates[0]
		if candidate.Content != nil {
			contents = append(contents, *candidate.Content)
		}

		// Check for function calls
		var functionCalls []FunctionCall
		if candidate.Content != nil {
			for _, part := range candidate.Content.Parts {
				if part.FunctionCall != nil {
					functionCalls = append(functionCalls, *part.FunctionCall)
				}
			}
		}

		if len(functionCalls) == 0 {
			// No function calls, extract text response
			if candidate.Content != nil {
				for _, part := range candidate.Content.Parts {
					if part.Text != "" {
						return part.Text, nil
					}
				}
			}
			return "", nil
		}

		// Execute function calls and build response
		var responseParts []Part
		for _, fc := range functionCalls {
			fmt.Println(tools.FormatToolCall("Gemini CLI", fc.Name, fc.Args))
			toolResult := tools.SafeExecute(fc.Name, fc.Args)
			responseParts = append(responseParts, Part{
				FunctionResponse: &FunctionResponse{
					Name:     fc.Name,
					Response: map[string]interface{}{"result": toolResult},
				},
			})
		}
		contents = append(contents, Content{Role: "user", Parts: responseParts})
	}
}

// GenerateStream sends a prompt and streams tokens to the channel.
func (p *Provider) GenerateStream(prompt string, attachments []adapter.Attachment, tokens chan<- string) error {
	ctx := context.Background()

	// 1. Initial Prompt Construction
	// Build parts with attachments
	var parts []Part
	for _, att := range attachments {
		if att.Type == adapter.AttachmentTypeText {
			parts = append(parts, Part{Text: fmt.Sprintf("--- File: %s ---\n%s\n--- End File ---", att.Name, att.Data)})
		} else if att.Type == adapter.AttachmentTypeImage {
			parts = append(parts, Part{
				InlineData: &Blob{
					MimeType: att.MimeType,
					Data:     att.Data,
				},
			})
		}
	}
	parts = append(parts, Part{Text: prompt})

	contents := []Content{
		{Role: "user", Parts: parts},
	}

	buildFuncDecls := func() []FunctionDeclaration {
		toolDefs := tools.GetToolDefinitionsForActiveSession()
		funcDecls := make([]FunctionDeclaration, 0, len(toolDefs))
		for _, td := range toolDefs {
			fn, ok := td["function"].(map[string]interface{})
			if !ok {
				continue
			}
			name, _ := fn["name"].(string)
			desc, _ := fn["description"].(string)
			params, _ := fn["parameters"].(map[string]interface{})
			funcDecls = append(funcDecls, FunctionDeclaration{
				Name:        name,
				Description: desc,
				Parameters:  params,
			})
		}
		return funcDecls
	}

	// 2. Conversation Loop (for tool use)
	for {
		funcDecls := buildFuncDecls()
		// Prepare Request Body
		genReq := GenerateContentRequest{
			Contents: contents,
			GenerationConfig: &GenerationConfig{
				MaxOutputTokens: 8192,
			},
		}

		if p.systemPrompt != "" {
			genReq.SystemInstruction = &SystemInstruction{
				Parts: []Part{{Text: p.systemPrompt}},
			}
		}

		reqBody := CloudCodeAssistRequest{
			Project: p.projectID,
			Model:   p.model,
			Request: genReq,
		}

		if len(funcDecls) > 0 {
			reqBody.Request.Tools = []Tool{{FunctionDeclarations: funcDecls}}
		}

		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			return err
		}

		var lastErr error
		var streamResponse *ResponseData // Accumulator for full response of this turn
		success := false

		// 3. API Call with Refresh & Retry Logic
		url := fmt.Sprintf("%s/v1internal:streamGenerateContent?alt=sse", DEFAULT_ENDPOINT)

		for attempt := 0; attempt <= MAX_RETRIES; attempt++ {
			if attempt > 0 {
				// Exponential backoff
				delay := time.Duration(BASE_DELAY_MS*(1<<(attempt-1))) * time.Millisecond
				time.Sleep(delay)
			}

			req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBody))
			if err != nil {
				lastErr = err
				break
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "text/event-stream")
			req.Header.Set("Authorization", "Bearer "+p.accessToken)
			for k, v := range GEMINI_CLI_HEADERS {
				req.Header.Set(k, v)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				lastErr = err
				continue
			}

			// Handle 401 Unauthorized - Refresh Token
			if resp.StatusCode == http.StatusUnauthorized {
				resp.Body.Close()
				fmt.Println("🔄 Gemini CLI Token expired (401), refreshing...")
				if newToken, err := p.Login(ctx); err == nil {
					fmt.Println("✅ Token refreshed successfully")
					p.accessToken = newToken.Access
					// Re-marshal body if needed (project ID shouldn't change here usually, but good practice)
					// reqBody.Project = p.projectID
					// jsonBody, _ = json.Marshal(reqBody)
					continue // Retry with new token
				} else {
					fmt.Printf("⚠️ Token refresh failed: %v\n", err)
					lastErr = fmt.Errorf("token expired and refresh failed: %v", err)
					break // Fatal auth error
				}
			}

			if resp.StatusCode == 429 {
				body, _ := io.ReadAll(resp.Body)
				resp.Body.Close()
				lastErr = fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
				continue // Retry on rate limit
			}

			if resp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				resp.Body.Close()
				lastErr = fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
				break // Don't retry non-429/401 errors
			}

			// Parse SSE Stream
			streamResponse, err = p.parseAndStreamSSE(resp.Body, tokens)
			resp.Body.Close()
			if err != nil {
				lastErr = err
				break
			}

			success = true
			break
		}

		if !success {
			return lastErr
		}

		// 4. Handle Tool Calls from Response
		var functionCalls []FunctionCall
		var assistantParts []Part

		if len(streamResponse.Candidates) > 0 {
			candidate := streamResponse.Candidates[0]
			if candidate.Content != nil {
				assistantParts = candidate.Content.Parts
				for _, part := range assistantParts {
					if part.FunctionCall != nil {
						functionCalls = append(functionCalls, *part.FunctionCall)
					}
				}
			}
		}

		// Append assistant response to history
		contents = append(contents, Content{Role: "model", Parts: assistantParts})

		if len(functionCalls) == 0 {
			// Done, no tools called
			return nil
		}

		// 5. Execute Tools
		var responseParts []Part
		for _, fc := range functionCalls {
			fmt.Println(tools.FormatToolCall("Gemini CLI", fc.Name, fc.Args))
			toolResult := tools.SafeExecute(fc.Name, fc.Args)
			responseParts = append(responseParts, Part{
				FunctionResponse: &FunctionResponse{
					Name:     fc.Name,
					Response: map[string]interface{}{"result": toolResult},
				},
			})
		}

		// Append tool results to history
		contents = append(contents, Content{Role: "user", Parts: responseParts})
		// Loop continues to generate model response to tool output
	}
}

// parseAndStreamSSE reads scanning SSE lines, streaming text to tokens chan, and accumulating full response
func (p *Provider) parseAndStreamSSE(body io.Reader, tokens chan<- string) (*ResponseData, error) {
	reader := bufio.NewReader(body)
	var finalResult ResponseData

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "data:") {
			continue
		}

		jsonStr := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if jsonStr == "" {
			continue
		}

		var chunk CloudCodeAssistResponseChunk
		if err := json.Unmarshal([]byte(jsonStr), &chunk); err != nil {
			continue
		}

		if chunk.Response == nil {
			continue
		}

		// Accumulate candidates
		if len(chunk.Response.Candidates) > 0 {
			candidate := chunk.Response.Candidates[0]
			if candidate.Content != nil {
				// Initialize structure if empty
				if len(finalResult.Candidates) == 0 {
					finalResult.Candidates = append(finalResult.Candidates, Candidate{
						Content: &Content{
							Role:  candidate.Content.Role,
							Parts: []Part{},
						},
					})
				}

				// Stream text parts and accumulate all parts
				for _, part := range candidate.Content.Parts {
					if part.Text != "" && tokens != nil {
						tokens <- part.Text
					}
				}

				// Warning: Simple accumulation here. Use a more robust merging strategy
				// for partial tool calls if Antigravity splits them across chunks.
				// For now, assuming parts are complete or additive list elements.
				finalResult.Candidates[0].Content.Parts = append(
					finalResult.Candidates[0].Content.Parts,
					candidate.Content.Parts...,
				)

				if candidate.FinishReason != "" {
					finalResult.Candidates[0].FinishReason = candidate.FinishReason
				}
			}
		}

		// Update usage
		if chunk.Response.UsageMetadata != nil {
			finalResult.UsageMetadata = chunk.Response.UsageMetadata
		}
	}

	return &finalResult, nil
}

func (p *Provider) callAPIStream(ctx context.Context, contents []Content, funcDecls []FunctionDeclaration, tokens chan<- string) error {
	// Function deprecated in favor of loop inside GenerateStream, but kept for interface compatibility if needed
	// or internal single-turn calls. For now, it's just a wrapper around the logic inside GenerateStream's loop body
	// but without the loop. To avoid duplication, we could extract the loop body, but GenerateStream is the main entry.
	// We'll leave this as a basic implementation or just return error as it's not main path.
	// Actually, let's just use the same logic as callAPI but with streaming.
	return fmt.Errorf("callAPIStream is deprecated, utilize GenerateStream")
}

func (p *Provider) callAPI(ctx context.Context, contents []Content, funcDecls []FunctionDeclaration) (*ResponseData, error) {
	// Use Cloud Code Assist endpoint (same as pi-ai SDK)
	url := fmt.Sprintf("%s/v1internal:streamGenerateContent?alt=sse", DEFAULT_ENDPOINT)

	// Build request body (same structure as pi-ai SDK)
	reqBody := CloudCodeAssistRequest{
		Project: p.projectID,
		Model:   p.model,
		Request: GenerateContentRequest{
			Contents: contents,
			GenerationConfig: &GenerationConfig{
				MaxOutputTokens: 8192,
			},
		},
	}

	if p.systemPrompt != "" {
		reqBody.Request.SystemInstruction = &SystemInstruction{
			Parts: []Part{{Text: p.systemPrompt}},
		}
	}

	if len(funcDecls) > 0 {
		reqBody.Request.Tools = []Tool{{FunctionDeclarations: funcDecls}}
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	var lastErr error
	for attempt := 0; attempt <= MAX_RETRIES; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			delay := time.Duration(BASE_DELAY_MS*(1<<(attempt-1))) * time.Millisecond
			fmt.Printf("[Rate limited, retrying in %v...]\n", delay)
			time.Sleep(delay)
		}

		req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBody))
		if err != nil {
			return nil, err
		}

		// Set headers exactly like pi-ai SDK
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "text/event-stream")
		req.Header.Set("Authorization", "Bearer "+p.accessToken)
		for k, v := range GEMINI_CLI_HEADERS {
			req.Header.Set(k, v)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		// Handle 401 Unauthorized - Refresh Token
		if resp.StatusCode == http.StatusUnauthorized {
			resp.Body.Close()
			fmt.Println("🔄 Gemini CLI Token expired (401), refreshing...")
			if newToken, err := p.Login(ctx); err == nil {
				p.accessToken = newToken.Access
				continue // Retry with new token
			} else {
				lastErr = fmt.Errorf("token expired and refresh failed: %v", err)
				break
			}
		}

		if resp.StatusCode == 429 {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			lastErr = fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
			continue // Retry on rate limit
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
		}

		// Parse SSE stream
		result, err := p.parseSSEStream(resp.Body)
		resp.Body.Close()
		return result, err
	}

	return nil, lastErr
}

func (p *Provider) parseSSEStream(body io.Reader) (*ResponseData, error) {
	reader := bufio.NewReader(body)
	var result ResponseData

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "data:") {
			continue
		}

		jsonStr := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if jsonStr == "" {
			continue
		}

		var chunk CloudCodeAssistResponseChunk
		if err := json.Unmarshal([]byte(jsonStr), &chunk); err != nil {
			continue
		}

		if chunk.Response == nil {
			continue
		}

		// Accumulate candidates
		if len(chunk.Response.Candidates) > 0 {
			candidate := chunk.Response.Candidates[0]
			if candidate.Content != nil {
				if len(result.Candidates) == 0 {
					result.Candidates = append(result.Candidates, Candidate{
						Content: &Content{
							Role:  candidate.Content.Role,
							Parts: []Part{},
						},
					})
				}
				// Append parts
				result.Candidates[0].Content.Parts = append(
					result.Candidates[0].Content.Parts,
					candidate.Content.Parts...,
				)
				if candidate.FinishReason != "" {
					result.Candidates[0].FinishReason = candidate.FinishReason
				}
			}
		}

		// Update usage
		if chunk.Response.UsageMetadata != nil {
			result.UsageMetadata = chunk.Response.UsageMetadata
		}
	}

	return &result, nil
}
