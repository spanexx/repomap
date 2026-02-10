package copilot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/spanexx/agents-cli/repomap/pkg/adapter"
	"github.com/spanexx/agents-cli/repomap/pkg/tools"
)

const (
	DEFAULT_MODEL = "gpt-4o"
	MAX_RETRIES   = 3
	BASE_DELAY_MS = 2000
)

type Provider struct {
	accessToken  string
	baseURL      string
	model        string
	systemPrompt string
}

func New(ctx context.Context) (*Provider, error) {
	return &Provider{
		model: DEFAULT_MODEL,
	}, nil
}

func (p *Provider) SetAccessToken(token string) {
	p.accessToken = token
	// Extract base URL from token
	p.baseURL = GetBaseURLFromToken(token)
}

func (p *Provider) Init(ctx context.Context) error {
	if p.accessToken == "" {
		return fmt.Errorf("access token is required")
	}
	if p.baseURL == "" {
		p.baseURL = "https://api.individual.githubcopilot.com"
	}
	return nil
}

func (p *Provider) Name() string {
	return "copilot"
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

// OpenAI-compatible request structures
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
	Tools    []Tool    `json:"tools,omitempty"`
}

type Message struct {
	Role       string     `json:"role"`
	Content    any        `json:"content"` // string or []ContentPart
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
}

type ContentPart struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
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

// Response structures
type ChatResponse struct {
	Choices []Choice `json:"choices"`
	Usage   *Usage   `json:"usage,omitempty"`
}

type StreamChunk struct {
	Choices []StreamChoice `json:"choices"`
	Usage   *Usage         `json:"usage,omitempty"`
}

type Choice struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type StreamChoice struct {
	Delta        Delta  `json:"delta"`
	FinishReason string `json:"finish_reason"`
}

type Delta struct {
	Role      string     `json:"role,omitempty"`
	Content   string     `json:"content,omitempty"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// GenerateStream implements streaming response (simulated).
func (c *Provider) GenerateStream(prompt string, attachments []adapter.Attachment, tokens chan<- string) error {
	resp, err := c.Generate(prompt, attachments)
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

	// Build prompt with text file attachments prepended
	fullPrompt := prompt
	for _, att := range attachments {
		if att.Type == adapter.AttachmentTypeText {
			fullPrompt = fmt.Sprintf("--- File: %s ---\n%s\n--- End File ---\n\n%s", att.Name, att.Data, fullPrompt)
		} else if att.Type == adapter.AttachmentTypeImage {
			fullPrompt = fmt.Sprintf("[Image attached: %s]\n\n%s", att.Name, fullPrompt)
		}
	}

	// Initial conversation
	messages := []Message{
		{Role: "user", Content: fullPrompt},
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

		if len(result.Choices) == 0 {
			return "", fmt.Errorf("no choices in response")
		}

		choice := result.Choices[0]
		messages = append(messages, choice.Message)

		// Check for tool calls
		if len(choice.Message.ToolCalls) == 0 {
			// No tool calls, return text response
			if content, ok := choice.Message.Content.(string); ok {
				return content, nil
			}
			return "", nil
		}

		// Execute tool calls
		for _, tc := range choice.Message.ToolCalls {
			var args map[string]interface{}
			json.Unmarshal([]byte(tc.Function.Arguments), &args)

			fmt.Println(tools.FormatToolCall("Copilot", tc.Function.Name, args))
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
	url := p.baseURL + "/chat/completions"

	reqBody := ChatRequest{
		Model:    p.model,
		Messages: messages,
		Stream:   false, // Non-streaming for simplicity
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

		// Set headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+p.accessToken)

		// Copilot-specific headers
		for k, v := range COPILOT_HEADERS {
			req.Header.Set(k, v)
		}

		// Determine X-Initiator based on last message role
		isAgentCall := len(messages) > 0 && messages[len(messages)-1].Role != "user"
		if isAgentCall {
			req.Header.Set("X-Initiator", "agent")
		} else {
			req.Header.Set("X-Initiator", "user")
		}
		req.Header.Set("Openai-Intent", "conversation-edits")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		if resp.StatusCode == 429 {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			lastErr = fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
			delay := time.Duration(BASE_DELAY_MS*(1<<(attempt+1))) * time.Millisecond
			time.Sleep(delay)
			continue
		}

		if resp.StatusCode == http.StatusUnauthorized {
			resp.Body.Close()
			fmt.Println("🔄 Token expired, refreshing...")
			if newToken, err := p.Login(ctx); err == nil {
				p.SetAccessToken(newToken.AccessToken)
				continue // Retry with new token
			}
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
