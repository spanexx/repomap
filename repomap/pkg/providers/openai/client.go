package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/spanexx/agents-cli/repomap/pkg/adapter"
	"github.com/spanexx/agents-cli/repomap/pkg/tools"
)

const (
	BASE_URL      = "https://api.openai.com/v1"
	DEFAULT_MODEL = "gpt-4o"
	MAX_RETRIES   = 3
	BASE_DELAY_MS = 2000
)

type Provider struct {
	apiKey       string
	model        string
	systemPrompt string
}

func New(ctx context.Context, apiKey string) (*Provider, error) {
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}
	return &Provider{
		apiKey: apiKey,
		model:  DEFAULT_MODEL,
	}, nil
}

func (p *Provider) Name() string {
	return "openai"
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

// Request structures
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
	Tools    []Tool    `json:"tools,omitempty"`
}

type Message struct {
	Role       string     `json:"role"`
	Content    any        `json:"content"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
}

type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
}

type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

type Function struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type ChatResponse struct {
	Choices []Choice `json:"choices"`
	Usage   *Usage   `json:"usage,omitempty"`
	Error   *Error   `json:"error,omitempty"`
}

type Choice struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Error struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    any    `json:"code"`
}

// ContentPart represents a part of multimodal content for OpenAI
type ContentPart struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

type ImageURL struct {
	URL string `json:"url"`
}

// GenerateStream implements streaming response (simulated).
func (p *Provider) GenerateStream(prompt string, attachments []adapter.Attachment, tokens chan<- string) error {
	resp, err := p.Generate(prompt, attachments)
	if err != nil {
		return err
	}
	tokens <- resp
	return nil
}

func (p *Provider) Generate(prompt string, attachments []adapter.Attachment) (string, error) {
	ctx := context.Background()

	buildAPITools := func() []Tool {
		toolDefs := tools.GetToolDefinitionsForActiveSession()
		apiTools := make([]Tool, 0, len(toolDefs))
		for _, td := range toolDefs {
			fn, ok := td["function"].(map[string]interface{})
			if !ok {
				continue
			}
			name, _ := fn["name"].(string)
			desc, _ := fn["description"].(string)
			params, _ := fn["parameters"].(map[string]interface{})
			apiTools = append(apiTools, Tool{
				Type: "function",
				Function: Function{
					Name:        name,
					Description: desc,
					Parameters:  params,
				},
			})
		}
		return apiTools
	}

	// Build content - use multimodal format if attachments exist
	var content any
	if len(attachments) > 0 {
		var parts []ContentPart
		// Add attachment parts first
		for _, att := range attachments {
			switch att.Type {
			case adapter.AttachmentTypeImage:
				parts = append(parts, ContentPart{
					Type: "image_url",
					ImageURL: &ImageURL{
						URL: fmt.Sprintf("data:%s;base64,%s", att.MimeType, att.Data),
					},
				})
			case adapter.AttachmentTypeText:
				parts = append(parts, ContentPart{
					Type: "text",
					Text: fmt.Sprintf("--- File: %s ---\n%s\n--- End File ---", att.Name, att.Data),
				})
			}
		}
		// Add the user prompt
		parts = append(parts, ContentPart{Type: "text", Text: prompt})
		content = parts
	} else {
		content = prompt
	}

	messages := []Message{
		{Role: "user", Content: content},
	}
	if p.systemPrompt != "" {
		messages = append([]Message{{Role: "system", Content: p.systemPrompt}}, messages...)
	}

	for {
		apiTools := buildAPITools()
		result, err := p.callAPI(ctx, messages, apiTools)
		if err != nil {
			return "", err
		}

		if result.Error != nil {
			return "", fmt.Errorf("API error: %s", result.Error.Message)
		}

		if len(result.Choices) == 0 {
			return "", fmt.Errorf("no choices in response")
		}

		choice := result.Choices[0]
		messages = append(messages, choice.Message)

		if len(choice.Message.ToolCalls) == 0 {
			if content, ok := choice.Message.Content.(string); ok {
				return content, nil
			}
			return "", nil
		}

		for _, tc := range choice.Message.ToolCalls {
			var args map[string]interface{}
			json.Unmarshal([]byte(tc.Function.Arguments), &args)

			fmt.Println(tools.FormatToolCall("OpenAI", tc.Function.Name, args))
			toolResult := tools.SafeExecute(tc.Function.Name, args)

			messages = append(messages, Message{
				Role:       "tool",
				Content:    toolResult,
				ToolCallID: tc.ID,
			})
		}
	}
}

func (p *Provider) callAPI(ctx context.Context, messages []Message, apiTools []Tool) (*ChatResponse, error) {
	url := BASE_URL + "/chat/completions"

	reqBody := ChatRequest{
		Model:    p.model,
		Messages: messages,
		Stream:   false,
	}

	if len(apiTools) > 0 {
		reqBody.Tools = apiTools
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	var lastErr error
	for attempt := 0; attempt <= MAX_RETRIES; attempt++ {
		if attempt > 0 {
			delay := time.Duration(BASE_DELAY_MS*(1<<(attempt-1))) * time.Millisecond
			fmt.Printf("[Rate limited, retrying in %v...]\n", delay)
			time.Sleep(delay)
		}

		req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBody))
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+p.apiKey)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		if resp.StatusCode == 429 {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			lastErr = fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
			continue
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
		}

		var result ChatResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			resp.Body.Close()
			return nil, err
		}
		resp.Body.Close()

		return &result, nil
	}

	return nil, lastErr
}
