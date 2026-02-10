package qwen

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"llm-adapter/pkg/adapter"
	"llm-adapter/pkg/tools"
	"net/http"
	"os"
)

// GenerateStream sends a prompt and streams tokens to the channel.
// GenerateStream sends a prompt and streams tokens to the channel.
func (p *Provider) GenerateStream(prompt string, attachments []adapter.Attachment, tokens chan<- string) error {
	var content any
	var hasImages bool

	// Check if we have images (logic identical to Generate)
	for _, att := range attachments {
		if att.Type == adapter.AttachmentTypeImage {
			hasImages = true
			break
		}
	}

	if hasImages {
		var parts []ContentPart
		for _, att := range attachments {
			if att.Type == adapter.AttachmentTypeText {
				parts = append(parts, ContentPart{
					Type: "text",
					Text: fmt.Sprintf("--- File: %s ---\n%s\n--- End File ---", att.Name, att.Data),
				})
			}
		}
		parts = append(parts, ContentPart{
			Type: "text",
			Text: prompt,
		})
		for _, att := range attachments {
			if att.Type == adapter.AttachmentTypeImage {
				parts = append(parts, ContentPart{
					Type: "image_url",
					ImageURL: &ImageURL{
						URL: fmt.Sprintf("data:%s;base64,%s", att.MimeType, att.Data),
					},
				})
			}
		}
		content = parts
	} else {
		fullPrompt := prompt
		for _, att := range attachments {
			if att.Type == adapter.AttachmentTypeText {
				fullPrompt = fmt.Sprintf("--- File: %s ---\n%s\n--- End File ---\n\n%s", att.Name, att.Data, fullPrompt)
			}
		}
		content = fullPrompt
	}

	// Add System Prompt if set and not already present
	// Since we append to conversation, we should set it if conversation is empty or if we want to force it.
	// Simple approach: If p.systemPrompt is set, ensure it's the first message if possible.
	// But p.conversation persists. If we append every time, it's weird.
	// However, usually system prompt is static for the session.
	// Let's assume we just want to ensure it's there.
	// Better: Use the provider's systemPrompt for the current turn.
	// But Qwen API expects messages.
	// Strategy: Prepend System prompt to the messages sent *in this request*,
	// OR append to p.conversation if it's the start.
	if len(p.conversation) == 0 && p.systemPrompt != "" {
		p.conversation = append(p.conversation, Message{Role: "system", Content: p.systemPrompt})
	}

	// Append the new user message to the conversation history
	p.conversation = append(p.conversation, Message{Role: "user", Content: content})

	for {
		rb := reqBody{
			Model:    p.model,
			Messages: p.conversation,
			Tools:    nil, // Start with no tools, logic below will set it
		}

		// Create a separate struct for streaming request as it has 'stream: true'
		type streamReqBody struct {
			reqBody
			Stream bool `json:"stream"`
		}

		if os.Getenv("LLM_ADAPTER_DISABLE_TOOLS") != "1" {
			rb.Tools = anySlice(tools.GetToolDefinitionsForActiveSession())
		}

		srb := streamReqBody{
			reqBody: rb,
			Stream:  true,
		}

		jsonData, err := json.Marshal(srb)
		if err != nil {
			return err
		}

		req, err := http.NewRequest("POST", p.baseURL, bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+p.accessToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		if resp.StatusCode == http.StatusUnauthorized {
			resp.Body.Close() // Close before retrying
			fmt.Println("🔄 Token expired, refreshing...")
			if newToken, err := p.Login(context.Background()); err == nil {
				p.accessToken = newToken.AccessToken
				req.Header.Set("Authorization", "Bearer "+p.accessToken)
				// Re-create body/request just in case? No, reset body reader
				// Actually we need to recreate the request because Body is consumed/closed
				req, _ = http.NewRequest("POST", p.baseURL, bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+p.accessToken)

				client = &http.Client{} // Reuse or new client
				resp, err = client.Do(req)
				if err != nil {
					return err
				}
			}
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return fmt.Errorf("qwen api error: %d %s", resp.StatusCode, string(body))
		}

		reader := bufio.NewReader(resp.Body)

		var fullContent string
		var toolCallBuffer []ToolCall
		// Map simple index to partial tool call
		partialToolCalls := make(map[int]*ToolCall)

		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				resp.Body.Close()
				return err
			}

			line = bytes.TrimSpace(line)
			if len(line) == 0 {
				continue
			}

			if !bytes.HasPrefix(line, []byte("data: ")) {
				continue
			}

			data := bytes.TrimPrefix(line, []byte("data: "))
			if string(data) == "[DONE]" {
				break
			}

			// Define structure for streaming response with tool calls
			var chunk struct {
				Choices []struct {
					Delta struct {
						Content   string `json:"content"`
						ToolCalls []struct {
							Index    int    `json:"index"`
							ID       string `json:"id"`
							Type     string `json:"type"`
							Function struct {
								Name      string `json:"name"`
								Arguments string `json:"arguments"`
							} `json:"function"`
						} `json:"tool_calls"`
					} `json:"delta"`
				} `json:"choices"`
			}

			if err := json.Unmarshal(data, &chunk); err != nil {
				continue
			}

			if len(chunk.Choices) > 0 {
				delta := chunk.Choices[0].Delta

				// Handle Content
				if delta.Content != "" {
					fullContent += delta.Content
					tokens <- delta.Content
				}

				// Handle Tool Calls
				for _, tc := range delta.ToolCalls {
					if _, exists := partialToolCalls[tc.Index]; !exists {
						partialToolCalls[tc.Index] = &ToolCall{
							ID:   tc.ID,
							Type: tc.Type,
							Function: Function{
								Name:      tc.Function.Name,
								Arguments: "",
							},
						}
					}

					// Accumulate data
					if tc.Function.Name != "" {
						partialToolCalls[tc.Index].Function.Name = tc.Function.Name
					}
					if tc.Function.Arguments != "" {
						partialToolCalls[tc.Index].Function.Arguments += tc.Function.Arguments
					}
					if tc.ID != "" {
						partialToolCalls[tc.Index].ID = tc.ID
					}
					if tc.Type != "" {
						partialToolCalls[tc.Index].Type = tc.Type
					}
				}
			}
		}

		resp.Body.Close()

		// Convert partial tool calls to final slice
		for i := 0; ; i++ {
			if tc, ok := partialToolCalls[i]; ok {
				// Fix missing defaults if any
				if tc.Type == "" {
					tc.Type = "function"
				}
				toolCallBuffer = append(toolCallBuffer, *tc)
			} else {
				break
			}
		}

		// Append assistant message
		msg := Message{Role: "assistant", Content: fullContent}
		if len(toolCallBuffer) > 0 {
			msg.ToolCalls = toolCallBuffer
		}
		p.conversation = append(p.conversation, msg)

		// Check if we need to execute tools
		if len(toolCallBuffer) == 0 {
			return nil // Done, just text response
		}

		// Execute tools
		for _, call := range toolCallBuffer {
			var args map[string]interface{}
			if err := json.Unmarshal([]byte(call.Function.Arguments), &args); err != nil {
				// Handle invalid JSON in arguments, maybe log token error?
				errRes := fmt.Sprintf("Error parsing arguments: %v", err)
				p.conversation = append(p.conversation, Message{
					Role:       "tool",
					Content:    errRes,
					ToolCallID: call.ID,
				})
				continue
			}

			// Execute tool - this emits event to server automatically via SafeExecute
			res := tools.SafeExecute(call.Function.Name, args)

			// Append result
			p.conversation = append(p.conversation, Message{
				Role:       "tool",
				Content:    res,
				ToolCallID: call.ID,
			})
		}
		// LOOP continues to next turn to generate response to tool outputs
	}
}
