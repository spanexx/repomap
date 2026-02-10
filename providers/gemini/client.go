package gemini

import (
	"context"
	"fmt"
	"llm-adapter/pkg/adapter"
	"llm-adapter/pkg/tools"
	"strings"

	"google.golang.org/genai"
)

type Provider struct {
	client       *genai.Client
	model        string
	systemPrompt string
}

func New(ctx context.Context, apiKey string) (*Provider, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
	}
	return &Provider{
		client: client,
		model:  "gemini-2.0-flash", // latest default
	}, nil
}

func (p *Provider) Name() string {
	return "gemini"
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

func (p *Provider) Generate(prompt string, attachments []adapter.Attachment) (string, error) {
	ctx := context.Background()

	// Build parts array for multimodal content
	var parts []*genai.Part

	// Add attachments first
	for _, att := range attachments {
		switch att.Type {
		case adapter.AttachmentTypeImage:
			parts = append(parts, genai.NewPartFromBytes([]byte(att.Data), att.MimeType))
		case adapter.AttachmentTypeText:
			parts = append(parts, genai.NewPartFromText(fmt.Sprintf("--- File: %s ---\n%s\n--- End File ---", att.Name, att.Data)))
		}
	}

	// Add the user prompt
	parts = append(parts, genai.NewPartFromText(prompt))

	messages := []*genai.Content{
		{Parts: parts, Role: "user"},
	}

	// Transform shared tool definitions to genai.Tool
	var toSchema func(v any) *genai.Schema
	toSchema = func(v any) *genai.Schema {
		vm, ok := v.(map[string]interface{})
		if !ok || vm == nil {
			return &genai.Schema{Type: "OBJECT"}
		}

		s := &genai.Schema{}
		if t, ok := vm["type"].(string); ok && strings.TrimSpace(t) != "" {
			s.Type = genai.Type(strings.ToUpper(strings.TrimSpace(t)))
		} else {
			s.Type = "OBJECT"
		}
		if d, ok := vm["description"].(string); ok {
			s.Description = d
		}

		if items, ok := vm["items"]; ok && items != nil {
			s.Items = toSchema(items)
		}

		if propsAny, ok := vm["properties"]; ok && propsAny != nil {
			if props, ok := propsAny.(map[string]interface{}); ok {
				out := make(map[string]*genai.Schema)
				for k, vv := range props {
					out[k] = toSchema(vv)
				}
				s.Properties = out
			}
		}

		if req, ok := vm["required"]; ok && req != nil {
			switch rr := req.(type) {
			case []string:
				s.Required = rr
			case []interface{}:
				var required []string
				for _, r := range rr {
					if rs, ok := r.(string); ok {
						required = append(required, rs)
					}
				}
				s.Required = required
			}
		}

		return s
	}
	buildConfig := func() *genai.GenerateContentConfig {
		toolDefs := tools.GetToolDefinitionsForActiveSession()
		funcDecls := make([]*genai.FunctionDeclaration, 0, len(toolDefs))
		for _, td := range toolDefs {
			f := td["function"].(map[string]interface{})

			// Map parameters
			params := f["parameters"].(map[string]interface{})
			props, _ := params["properties"].(map[string]interface{})

			genaiProps := make(map[string]*genai.Schema)
			for k, v := range props {
				genaiProps[k] = toSchema(v)
			}

			// Map required fields safely
			var required []string
			if req, ok := params["required"]; ok && req != nil {
				if reqSlice, ok := req.([]string); ok {
					required = reqSlice
				} else if reqSlice, ok := req.([]interface{}); ok {
					for _, r := range reqSlice {
						if rs, ok := r.(string); ok {
							required = append(required, rs)
						}
					}
				}
			}

			funcDecls = append(funcDecls, &genai.FunctionDeclaration{
				Name:        f["name"].(string),
				Description: f["description"].(string),
				Parameters: &genai.Schema{
					Type:       "OBJECT",
					Properties: genaiProps,
					Required:   required,
				},
			})
		}

		config := &genai.GenerateContentConfig{
			Tools: []*genai.Tool{{FunctionDeclarations: funcDecls}},
		}
		if p.systemPrompt != "" {
			config.SystemInstruction = genai.NewContentFromText(p.systemPrompt, "system")
		}
		return config
	}

	for {
		config := buildConfig()
		result, err := p.client.Models.GenerateContent(ctx, p.model, messages, config)
		if err != nil {
			return "", err
		}

		if len(result.Candidates) == 0 {
			return "", fmt.Errorf("no response from model")
		}

		responseContent := result.Candidates[0].Content
		var foundFunctionCall bool
		for _, part := range responseContent.Parts {
			if part.Text != "" {
				return part.Text, nil
			}

			if part.FunctionCall != nil {
				foundFunctionCall = true
				fmt.Println(tools.FormatToolCall("Gemini", part.FunctionCall.Name, part.FunctionCall.Args))

				res := tools.SafeExecute(part.FunctionCall.Name, part.FunctionCall.Args)

				messages = append(messages, responseContent)
				messages = append(messages, genai.NewContentFromFunctionResponse(
					part.FunctionCall.Name,
					map[string]any{"result": res},
					"user",
				))
				break
			}
		}

		if !foundFunctionCall {
			return "", fmt.Errorf("no text or function call in response")
		}
	}
}

func (p *Provider) GenerateStream(prompt string, attachments []adapter.Attachment, tokens chan<- string) error {
	ctx := context.Background()

	// Build parts array for multimodal content
	var parts []*genai.Part

	// Add attachments first
	for _, att := range attachments {
		switch att.Type {
		case adapter.AttachmentTypeImage:
			parts = append(parts, genai.NewPartFromBytes([]byte(att.Data), att.MimeType))
		case adapter.AttachmentTypeText:
			parts = append(parts, genai.NewPartFromText(fmt.Sprintf("--- File: %s ---\n%s\n--- End File ---", att.Name, att.Data)))
		}
	}

	// Add the user prompt
	parts = append(parts, genai.NewPartFromText(prompt))

	messages := []*genai.Content{
		{Parts: parts, Role: "user"},
	}

	// Transform shared tool definitions to genai.Tool
	var toSchema func(v any) *genai.Schema
	toSchema = func(v any) *genai.Schema {
		vm, ok := v.(map[string]interface{})
		if !ok || vm == nil {
			return &genai.Schema{Type: "OBJECT"}
		}

		s := &genai.Schema{}
		if t, ok := vm["type"].(string); ok && strings.TrimSpace(t) != "" {
			s.Type = genai.Type(strings.ToUpper(strings.TrimSpace(t)))
		} else {
			s.Type = "OBJECT"
		}
		if d, ok := vm["description"].(string); ok {
			s.Description = d
		}

		if items, ok := vm["items"]; ok && items != nil {
			s.Items = toSchema(items)
		}

		if propsAny, ok := vm["properties"]; ok && propsAny != nil {
			if props, ok := propsAny.(map[string]interface{}); ok {
				out := make(map[string]*genai.Schema)
				for k, vv := range props {
					out[k] = toSchema(vv)
				}
				s.Properties = out
			}
		}

		if req, ok := vm["required"]; ok && req != nil {
			switch rr := req.(type) {
			case []string:
				s.Required = rr
			case []interface{}:
				var required []string
				for _, r := range rr {
					if rs, ok := r.(string); ok {
						required = append(required, rs)
					}
				}
				s.Required = required
			}
		}

		return s
	}
	buildConfig := func() *genai.GenerateContentConfig {
		toolDefs := tools.GetToolDefinitionsForActiveSession()
		funcDecls := make([]*genai.FunctionDeclaration, 0, len(toolDefs))
		for _, td := range toolDefs {
			f := td["function"].(map[string]interface{})

			// Map parameters
			params := f["parameters"].(map[string]interface{})
			props, _ := params["properties"].(map[string]interface{})

			genaiProps := make(map[string]*genai.Schema)
			for k, v := range props {
				genaiProps[k] = toSchema(v)
			}

			// Map required fields safely
			var required []string
			if req, ok := params["required"]; ok && req != nil {
				if reqSlice, ok := req.([]string); ok {
					required = reqSlice
				} else if reqSlice, ok := req.([]interface{}); ok {
					for _, r := range reqSlice {
						if rs, ok := r.(string); ok {
							required = append(required, rs)
						}
					}
				}
			}

			funcDecls = append(funcDecls, &genai.FunctionDeclaration{
				Name:        f["name"].(string),
				Description: f["description"].(string),
				Parameters: &genai.Schema{
					Type:       "OBJECT",
					Properties: genaiProps,
					Required:   required,
				},
			})
		}

		config := &genai.GenerateContentConfig{
			Tools: []*genai.Tool{{FunctionDeclarations: funcDecls}},
		}
		if p.systemPrompt != "" {
			config.SystemInstruction = genai.NewContentFromText(p.systemPrompt, "system")
		}
		return config
	}

	// Tool loop (Simplified: Text only for now)
	for {
		emittedAny := false
		sawFuncCall := false
		config := buildConfig()

		// Go 1.23 iterator
		for resp, err := range p.client.Models.GenerateContentStream(ctx, p.model, messages, config) {
			if err != nil {
				return err
			}
			if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
				continue
			}

			content := resp.Candidates[0].Content
			for _, part := range content.Parts {
				if part.Text != "" {
					emittedAny = true
					tokens <- part.Text
				}

				if part.FunctionCall != nil {
					sawFuncCall = true
					fmt.Println(tools.FormatToolCall("Gemini", part.FunctionCall.Name, part.FunctionCall.Args))
					res := tools.SafeExecute(part.FunctionCall.Name, part.FunctionCall.Args)

					messages = append(messages, content)
					messages = append(messages, genai.NewContentFromFunctionResponse(
						part.FunctionCall.Name,
						map[string]any{"result": res},
						"user",
					))
					break
				}
			}

			if sawFuncCall {
				break
			}
		}

		if sawFuncCall {
			continue
		}

		if !emittedAny {
			result, err := p.client.Models.GenerateContent(ctx, p.model, messages, config)
			if err != nil {
				return err
			}
			if len(result.Candidates) > 0 && result.Candidates[0].Content != nil {
				for _, part := range result.Candidates[0].Content.Parts {
					if part.Text != "" {
						tokens <- part.Text
						emittedAny = true
						break
					}
				}
			}
		}

		return nil
	}

}
