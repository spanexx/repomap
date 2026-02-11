package analysis

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spanexx/agents-cli/repomap/internal/output"
	"github.com/spanexx/agents-cli/repomap/pkg/adapter"
)

const (
	BatchSize = 10
	MaxTokens = 2000
)

// AssignIntent populates the Intent field of FileNodes.
// If a provider is supplied, it uses the LLM for high-importance files.
// Otherwise, it falls back to static heuristics.
func AssignIntent(files []*output.FileNode, provider adapter.Provider) {
	if provider == nil {
		log.Println("No provider configured, using static heuristics for intent analysis.")
		assignStatic(files)
		return
	}

	// Filter for high/medium importance files
	var targetFiles []*output.FileNode
	for _, f := range files {
		// Only analyze source code, skip simple configs/docs if needed
		if f.Importance == "high" || f.Importance == "medium" {
			targetFiles = append(targetFiles, f)
		} else {
			// Low importance gets static immediately
			f.Intent = defineIntent(f.Path)
		}
	}

	log.Printf("Analyzing intent for %d priority files using LLM...", len(targetFiles))

	// Process in batches
	for i := 0; i < len(targetFiles); i += BatchSize {
		end := i + BatchSize
		if end > len(targetFiles) {
			end = len(targetFiles)
		}
		batch := targetFiles[i:end]

		if err := analyzeBatch(batch, provider); err != nil {
			log.Printf("Batch analysis failed (falling back to static): %v", err)
			assignStatic(batch)
		} else {
			// Check for missing intents in batch and fill with static
			for _, f := range batch {
				if f.Intent == "" || f.Intent == "Unknown" {
					f.Intent = defineIntent(f.Path)
				}
			}
		}

		// Rate limit protection
		time.Sleep(200 * time.Millisecond)
	}

	log.Printf("Intent analysis complete.")
}

func assignStatic(files []*output.FileNode) {
	for _, f := range files {
		if f.Intent == "" {
			f.Intent = defineIntent(f.Path)
		}
	}
}

func analyzeBatch(batch []*output.FileNode, provider adapter.Provider) error {
	// Prepare Prompt
	var fileList strings.Builder
	for _, f := range batch {
		defs := ""
		if len(f.Definitions) > 0 {
			// Truncate definitions to save tokens
			count := len(f.Definitions)
			if count > 5 {
				defs = strings.Join(f.Definitions[:5], "; ") + fmt.Sprintf("... (+%d more)", count-5)
			} else {
				defs = strings.Join(f.Definitions, "; ")
			}
		}
		fileList.WriteString(fmt.Sprintf("- Path: %s\n  Definitions: %s\n", f.Path, defs))
	}

	prompt := fmt.Sprintf(`Analyze the following Go files and determine their specific architectural intent (3-5 words).
Files:
%s

Return ONLY a JSON map where the key is the file path and the value is the intent description (e.g., "User Authentication", "Database Connection", "API Handler", "Business Logic Core").
JSON:`, fileList.String())

	// Call LLM
	// We use a temporary system prompt for this specific task
	// provider.SetSystemPrompt("You are a code analysis engine. Output valid JSON only.")
	// Note: Changing system prompt might affect global state if provider is shared/concurrent.
	// Ideally passes system prompt in Generate, but interface uses SetSystemPrompt.
	// We'll assume sequential execution for analysis phase.
	provider.SetSystemPrompt("You are a code analysis engine. Output valid JSON only.")

	resp, err := provider.Generate(prompt, nil)
	if err != nil {
		return err
	}

	// Clean Code Block
	resp = strings.TrimSpace(resp)
	if strings.HasPrefix(resp, "```json") {
		resp = strings.TrimPrefix(resp, "```json")
		resp = strings.TrimSuffix(resp, "```")
	} else if strings.HasPrefix(resp, "```") {
		resp = strings.TrimPrefix(resp, "```")
		resp = strings.TrimSuffix(resp, "```")
	}

	// Parse JSON
	results := make(map[string]string)
	if err := json.Unmarshal([]byte(resp), &results); err != nil {
		return fmt.Errorf("invalid json regex: %v | response: %s", err, resp)
	}

	// Assign results
	for _, f := range batch {
		if intent, ok := results[f.Path]; ok {
			f.Intent = intent
		}
	}

	return nil
}

func defineIntent(path string) string {
	lowerPath := strings.ToLower(path)

	// Framework / Entry Points
	if strings.HasPrefix(lowerPath, "cmd/") || strings.HasPrefix(lowerPath, "main.go") {
		return "Entry Point"
	}

	// Domain / Core
	if strings.Contains(lowerPath, "domain/") || strings.Contains(lowerPath, "core/") || strings.Contains(lowerPath, "entity/") || strings.Contains(lowerPath, "model/") {
		return "Core Domain"
	}

	// Application / Usecase
	if strings.Contains(lowerPath, "usecase/") || strings.Contains(lowerPath, "application/") || strings.Contains(lowerPath, "service/") {
		return "Application Logic"
	}

	// Infrastructure / Adapters
	if strings.Contains(lowerPath, "infrastructure/") || strings.Contains(lowerPath, "infra/") || strings.Contains(lowerPath, "adapter/") || strings.Contains(lowerPath, "repository/") || strings.Contains(lowerPath, "handler/") || strings.Contains(lowerPath, "controller/") {
		return "Infrastructure"
	}

	// Configuration
	if strings.Contains(lowerPath, "config/") {
		return "Configuration"
	}

	// Utilities
	if strings.Contains(lowerPath, "util/") || strings.Contains(lowerPath, "common/") || strings.Contains(lowerPath, "lib/") || strings.Contains(lowerPath, "pkg/") {
		return "Utility"
	}

	// Tests
	if strings.HasSuffix(lowerPath, "_test.go") {
		return "Test"
	}

	// Default
	return "Implementation"
}
