package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"llm-adapter/pkg/adapter"
	"llm-adapter/pkg/tools"
	"net/http"
	"os"
	"time"
)

const (
	BASE_URL      = "https://api.anthropic.com/v1"
	DEFAULT_MODEL = "claude-3-5-sonnet-20241022"
	MAX_RETRIES   = 3
	BASE_DELAY_MS = 2000
	API_VERSION   = "2023-06-01"
)

type Provider struct {
	apiKey       string
	model        string
	systemPrompt string
}

func New(ctx context.Context, apiKey string) (*Provider, error) {
	if apiKey == "" {
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY environment variable not set")
	}
	return &Provider{
		apiKey: apiKey,
		model:  DEFAULT_MODEL,
	}, nil
}

func (p *Provider) Name() string {
	return "anthropic"
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

// Anthropic Request structures
type MessagesRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	System    string    `json:"system,omitempty"`
	MaxTokens int       `json:"max_tokens"`
	Tools     []Tool    `json:"tools,omitempty"`
}

type Message struct {
	Role    string         `json:"role"`
	Content []ContentBlock `json:"content"`
}

// ImageSource represents a base64-encoded image for Anthropic API
type ImageSource struct {
	Type      string `json:"type"`       // "base64"
	MediaType string `json:"media_type"` // e.g., "image/png"
	Data      string `json:"data"`       // base64 string
}

type ContentBlock struct {
	Type      string                 `json:"type"`
	Text      string                 `json:"text,omitempty"`
	Source    *ImageSource           `json:"source,omitempty"` // for images
	ID        string                 `json:"id,omitempty"`
	Name      string                 `json:"name,omitempty"`
	Input     map[string]interface{} `json:"input,omitempty"`
	ToolUseID string                 `json:"tool_use_id,omitempty"`
	Content   string                 `json:"content,omitempty"`
}

type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"input_schema"`
}

// Anthropic Response structures
type MessagesResponse struct {
	ID           string         `json:"id"`
	Type         string         `json:"type"`
	Role         string         `json:"role"`
	Content      []ContentBlock `json:"content"`
	StopReason   string         `json:"stop_reason"`
	StopSequence string         `json:"stop_sequence"`
	Usage        Usage          `json:"usage"`
	Error        *Error         `json:"error,omitempty"`
}

type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type Error struct {
	Type    string `json:"type"`
	Message string `json:"message"`
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
				Name:        name,
				Description: desc,
				InputSchema: params,
			})
		}
		return apiTools
	}

	// Build content blocks from attachments and prompt
	var content []ContentBlock

	// Add attachment content blocks first
	for _, att := range attachments {
		switch att.Type {
		case adapter.AttachmentTypeImage:
			content = append(content, ContentBlock{
				Type: "image",
				Source: &ImageSource{
					Type:      "base64",
					MediaType: att.MimeType,
					Data:      att.Data,
				},
			})
		case adapter.AttachmentTypeText:
			content = append(content, ContentBlock{
				Type: "text",
				Text: fmt.Sprintf("--- File: %s ---\n%s\n--- End File ---", att.Name, att.Data),
			})
		}
	}

	// Add the user prompt
	content = append(content, ContentBlock{Type: "text", Text: prompt})

	messages := []Message{
		{
			Role:    "user",
			Content: content,
		},
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

		// Append assistant response to history
		messages = append(messages, Message{
			Role:    "assistant",
			Content: result.Content,
		})

		// Check for tool usage
		var toolResults []ContentBlock
		hasToolUse := false

		for _, block := range result.Content {
			if block.Type == "tool_use" {
				hasToolUse = true
				fmt.Println(tools.FormatToolCall("Anthropic", block.Name, block.Input))

				toolResult := tools.SafeExecute(block.Name, block.Input)

				toolResults = append(toolResults, ContentBlock{
					Type:      "tool_result",
					ToolUseID: block.ID,
					Content:   toolResult,
				})
			}
		}

		if !hasToolUse {
			// Find text content
			for _, block := range result.Content {
				if block.Type == "text" {
					return block.Text, nil
				}
			}
			return "", nil
		}

		// Append tool results and continue
		messages = append(messages, Message{
			Role:    "user",
			Content: toolResults,
		})
	}
}

func (p *Provider) callAPI(ctx context.Context, messages []Message, apiTools []Tool) (*MessagesResponse, error) {
	url := BASE_URL + "/messages"

	reqBody := MessagesRequest{
		Model:     p.model,
		Messages:  messages,
		System:    p.systemPrompt,
		MaxTokens: 4096,
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
		req.Header.Set("x-api-key", p.apiKey)
		req.Header.Set("anthropic-version", API_VERSION)

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

		var result MessagesResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			resp.Body.Close()
			return nil, err
		}
		resp.Body.Close()

		return &result, nil
	}

	return nil, lastErr
}
