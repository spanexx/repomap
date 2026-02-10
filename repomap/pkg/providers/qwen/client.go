package qwen

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spanexx/agents-cli/repomap/pkg/adapter"
	"github.com/spanexx/agents-cli/repomap/pkg/tools"
)

type Provider struct {
	accessToken  string
	model        string
	baseURL      string
	conversation []Message
	systemPrompt string
}

func New(accessToken string) *Provider {
	return &Provider{
		accessToken:  accessToken,
		model:        "qwen2.5-coder-32b-instruct",
		baseURL:      "https://portal.qwen.ai/v1/chat/completions",
		conversation: []Message{},
	}
}

func (p *Provider) Name() string {
	return "qwen"
}

// SetModel sets the model to use for this provider.
func (p *Provider) SetModel(model string) {
	if model != "" {
		p.model = model
	}
}

func (p *Provider) SetSystemPrompt(prompt string) {
	p.systemPrompt = prompt
}

func (p *Provider) SetAccessToken(token string) {
	p.accessToken = token
}

type Message struct {
	Role       string     `json:"role"`
	Content    any        `json:"content"` // Changed from string to any to support []ContentPart
	ToolCallID string     `json:"tool_call_id,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
}

type ContentPart struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

type ImageURL struct {
	URL string `json:"url"`
}

type ToolCall struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

type Function struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type reqBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Tools    []any     `json:"tools,omitempty"`
}

// Generate sends a prompt to Qwen and maintains conversation history across multiple calls.
// Supports tool calling with automatic tool execution and response feedback.
func (p *Provider) Generate(prompt string, attachments []adapter.Attachment) (string, error) {
	var content any
	var hasImages bool

	// Check if we have images
	for _, att := range attachments {
		if att.Type == adapter.AttachmentTypeImage {
			hasImages = true
			break
		}
	}

	if hasImages {
		var parts []ContentPart

		// Add text attachments first
		for _, att := range attachments {
			if att.Type == adapter.AttachmentTypeText {
				parts = append(parts, ContentPart{
					Type: "text",
					Text: fmt.Sprintf("--- File: %s ---\n%s\n--- End File ---", att.Name, att.Data),
				})
			}
		}

		// Add prompt
		parts = append(parts, ContentPart{
			Type: "text",
			Text: prompt,
		})

		// Add images
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
		// Text only mode
		fullPrompt := prompt
		for _, att := range attachments {
			if att.Type == adapter.AttachmentTypeText {
				fullPrompt = fmt.Sprintf("--- File: %s ---\n%s\n--- End File ---\n\n%s", att.Name, att.Data, fullPrompt)
			}
		}
		content = fullPrompt
	}

	// Append the new user message to the conversation history
	if len(p.conversation) == 0 && p.systemPrompt != "" {
		p.conversation = append(p.conversation, Message{Role: "system", Content: p.systemPrompt})
	}
	p.conversation = append(p.conversation, Message{Role: "user", Content: content})

	for {
		rb := reqBody{
			Model:    p.model,
			Messages: p.conversation,
			Tools:    nil,
		}

		if os.Getenv("LLM_ADAPTER_DISABLE_TOOLS") != "1" {
			rb.Tools = anySlice(tools.GetToolDefinitionsForActiveSession())
		}

		jsonData, _ := json.Marshal(rb)
		req, _ := http.NewRequest("POST", p.baseURL, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+p.accessToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == http.StatusUnauthorized {
			// Token might be expired, try to refresh
			fmt.Println("🔄 Token expired, refreshing...")
			if _, err := p.Login(context.Background()); err == nil {
				// Login successful, p.accessToken should be updated?
				// Wait, Login returns token but doesn't strictly set p.accessToken in all paths in the current code?
				// Check Login implementation: It returns token.
				// We need to ensure p.accessToken is updated.
				// Let's assume Login or a wrapper does it, or update it here.
				// We'll read the token from file to be safe or use the return.
				if tok, err := p.LoadTokenFromFile(); err == nil {
					p.accessToken = tok.Access
					// Retry the request
					req.Header.Set("Authorization", "Bearer "+p.accessToken)
					resp, err = client.Do(req)
					if err == nil {
						body, _ = io.ReadAll(resp.Body)
						resp.Body.Close()
					}
				}
			}
		}

		if resp != nil && resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("qwen api error: %d %s", resp.StatusCode, string(body))
		}

		var qwenResp struct {
			Choices []struct {
				Message Message `json:"message"`
			} `json:"choices"`
		}
		if err := json.Unmarshal(body, &qwenResp); err != nil {
			return "", err
		}

		if len(qwenResp.Choices) == 0 {
			return "", fmt.Errorf("no response from qwen")
		}

		msg := qwenResp.Choices[0].Message
		p.conversation = append(p.conversation, msg)

		if len(msg.ToolCalls) == 0 {
			// Handle content being string or nil (though for response it's usually string, strictly mapped)
			// But since we changed Content to 'any', we need to assert
			if strContent, ok := msg.Content.(string); ok {
				return strContent, nil
			}
			return fmt.Sprintf("%v", msg.Content), nil
		}

		for _, call := range msg.ToolCalls {
			var args map[string]interface{}
			json.Unmarshal([]byte(call.Function.Arguments), &args)

			fmt.Printf("\n%s\n", tools.FormatToolCall("Qwen", call.Function.Name, args))
			res := tools.SafeExecute(call.Function.Name, args)

			p.conversation = append(p.conversation, Message{
				Role:       "tool",
				Content:    res,
				ToolCallID: call.ID,
			})
		}
	}
}

func anySlice(m []tools.ToolDefinition) []any {
	res := make([]any, len(m))
	for i, v := range m {
		res[i] = v
	}
	return res
}
