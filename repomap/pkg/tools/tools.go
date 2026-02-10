package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ToolDefinition map[string]interface{}

var activeTools []ToolDefinition

// RegisterTool adds a tool definition to the registry.
func RegisterTool(name, description string, parameters map[string]interface{}, handler func(map[string]interface{}) (string, error)) {
	def := ToolDefinition{
		"type": "function",
		"function": map[string]interface{}{
			"name":        name,
			"description": description,
			"parameters":  parameters,
		},
	}
	activeTools = append(activeTools, def)
	toolHandlers[name] = handler
}

var toolHandlers = make(map[string]func(map[string]interface{}) (string, error))

// GetToolDefinitionsForActiveSession returns all registered tools.
func GetToolDefinitionsForActiveSession() []ToolDefinition {
	return activeTools
}

// FormatToolCall formats a tool call for logging/display.
func FormatToolCall(provider, name string, args map[string]interface{}) string {
	argsJSON, _ := json.Marshal(args)
	return fmt.Sprintf("[%s] Using tool %s with args: %s", provider, name, string(argsJSON))
}

// SafeExecute runs the tool handler safely.
func SafeExecute(name string, args map[string]interface{}) string {
	handler, ok := toolHandlers[name]
	if !ok {
		return fmt.Sprintf("Error: Tool not found: %s", name)
	}

	result, err := handler(args)
	if err != nil {
		return fmt.Sprintf("Error executing tool %s: %v", name, err)
	}
	return result
}

// Core Filesystem Tools

func init() {
	// ReadFile Tool
	RegisterTool("read_file", "Read contents of a file", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"path": map[string]interface{}{"type": "string", "description": "Path to the file"},
		},
		"required": []string{"path"},
	}, func(args map[string]interface{}) (string, error) {
		path, ok := args["path"].(string)
		if !ok {
			return "", fmt.Errorf("invalid path argument")
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		return string(data), nil
	})

	// WriteFile Tool
	RegisterTool("write_file", "Write content to a file", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"path":    map[string]interface{}{"type": "string", "description": "Path to the file"},
			"content": map[string]interface{}{"type": "string", "description": "Content to write"},
		},
		"required": []string{"path", "content"},
	}, func(args map[string]interface{}) (string, error) {
		path, ok := args["path"].(string)
		if !ok {
			return "", fmt.Errorf("invalid path argument")
		}
		content, ok := args["content"].(string)
		if !ok {
			return "", fmt.Errorf("invalid content argument")
		}

		// Ensure directory exists
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return "", err
		}

		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return "", err
		}
		return fmt.Sprintf(" successfully wrote to %s", path), nil
	})

	// ListDir Tool
	RegisterTool("list_dir", "List files in a directory", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"path": map[string]interface{}{"type": "string", "description": "Directory path"},
		},
		"required": []string{"path"},
	}, func(args map[string]interface{}) (string, error) {
		path, ok := args["path"].(string)
		if !ok {
			return "", fmt.Errorf("invalid path argument")
		}
		entries, err := os.ReadDir(path)
		if err != nil {
			return "", err
		}
		var names []string
		for _, e := range entries {
			names = append(names, e.Name())
		}
		return strings.Join(names, "\n"), nil
	})
}
