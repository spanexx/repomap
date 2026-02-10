package generic

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"llm-adapter/pkg/adapter"
	"net/http"
)

// GenerateStream sends a prompt and streams tokens to the channel.
func (p *Provider) GenerateStream(prompt string, attachments []adapter.Attachment, tokens chan<- string) error {
	ctx := context.Background()

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

	messages := []Message{
		{Role: "user", Content: content},
	}

	return p.callAPIStream(ctx, messages, tokens)
}

func (p *Provider) callAPIStream(ctx context.Context, messages []Message, tokens chan<- string) error {
	url := p.baseURL + "/chat/completions"

	reqBody := ChatRequest{
		Model:    p.model,
		Messages: messages,
		Stream:   true,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	if p.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.apiKey)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		line = bytes.TrimSpace(line)
		if !bytes.HasPrefix(line, []byte("data: ")) {
			continue
		}

		data := bytes.TrimPrefix(line, []byte("data: "))
		if string(data) == "[DONE]" {
			break
		}

		var chunk struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
		}

		if err := json.Unmarshal(data, &chunk); err != nil {
			// Skip malformed chunks
			continue
		}

		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			tokens <- chunk.Choices[0].Delta.Content
		}
	}

	return nil
}
