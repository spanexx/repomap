package providers

import "strings"

type Kind string

const (
	KindToolLoop    Kind = "tool_loop"
	KindExternalCLI Kind = "external_cli"
	KindRemoteAPI   Kind = "remote_api"
	KindUnknown     Kind = "unknown"
)

func NormalizeProviderName(name string) string {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return ""
	}
	parts := strings.SplitN(trimmed, ":", 2)
	return strings.ToLower(strings.TrimSpace(parts[0]))
}

func KindForProvider(name string) Kind {
	p := NormalizeProviderName(name)
	switch p {
	case "claude", "opencode-cli", "qwen-cli", "qodercli", "gemini-shell":
		return KindExternalCLI
	case "gemini", "openai", "anthropic", "openrouter", "gateway", "venice", "xiaomi", "huggingface", "minimax", "moonshot", "zai":
		return KindRemoteAPI
	case "ollama", "lmstudio", "qwen":
		return KindToolLoop
	case "gemini-cli", "antigravity":
		// These are CLI-backed but are not currently treated as external shell-capable agents.
		// Keep them out of the external CLI list for now to avoid over-blocking.
		return KindToolLoop
	default:
		return KindUnknown
	}
}

func IsExternalCLIProvider(name string) bool {
	return KindForProvider(name) == KindExternalCLI
}
