package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spanexx/agents-cli/repomap/pkg/adapter"
	"github.com/spanexx/agents-cli/repomap/pkg/providers/gemini_cli"
	"github.com/spanexx/agents-cli/repomap/pkg/providers/qodercli"
	"github.com/spanexx/agents-cli/repomap/pkg/providers/qwen_cli"
)

type Server struct {
	Port     string
	Provider adapter.Provider
	PlanPath string
}

func New(port string, planPath string) *Server {
	return &Server{
		Port:     port,
		PlanPath: planPath,
	}
}

func (s *Server) Start() error {
	ctx := context.Background()

	providerName := os.Getenv("REPOMAP_PROVIDER")
	if providerName == "" {
		providerName = "gemini-cli"
	}
	fmt.Printf("Initializing provider(s): %s\n", providerName)

	providerNames := strings.Split(providerName, ",")
	var providers []adapter.Provider

	for _, name := range providerNames {
		name = strings.TrimSpace(name)
		var p adapter.Provider

		switch name {
		case "qwen-cli":
			p = qwen_cli.New()
		case "qodercli":
			p = qodercli.New()
		case "gemini-cli":
			gp, err := gemini_cli.New(ctx)
			if err != nil {
				fmt.Printf("Warning: failed to create gemini provider: %v\n", err)
				continue
			}
			if err := gp.Init(ctx); err != nil {
				fmt.Printf("Warning: Gemini provider init failed: %v\n", err)
			}
			p = gp
		default:
			fmt.Printf("Warning: unknown provider: %s\n", name)
			continue
		}
		if p != nil {
			providers = append(providers, p)
		}
	}

	if len(providers) == 0 {
		return fmt.Errorf("no valid providers found in: %s", providerName)
	}

	if len(providers) == 1 {
		s.Provider = providers[0]
	} else {
		s.Provider = adapter.NewFallbackProvider(providers, true)
		fmt.Printf("Fallback strategy enabled with: %v\n", providerNames)
	}

	http.HandleFunc("/api/plan", s.handlePlan)
	http.HandleFunc("/api/chat", s.handleChat)
	http.HandleFunc("/api/config", s.handleConfig)

	// Fix static path to point to visualizer-ui/dist
	uiPath := filepath.Join(filepath.Dir(s.PlanPath), "..", "..", "repomap", "visualizer-ui", "dist")
	// If that doesn't exist, try a relative path from the binary location or cwd
	if _, err := os.Stat(uiPath); os.IsNotExist(err) {
		uiPath = "./visualizer-ui/dist"
	}

	fmt.Printf("Serving UI from: %s\n", uiPath)
	http.Handle("/", http.FileServer(http.Dir(uiPath)))

	fmt.Printf("Starting server on :%s\n", s.Port)
	return http.ListenAndServe(":"+s.Port, nil)
}

func (s *Server) handlePlan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method == "GET" {
		data, err := os.ReadFile(s.PlanPath)
		if err != nil {
			if os.IsNotExist(err) {
				http.Error(w, "Plan not found", 404)
				return
			}
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
		return
	}

	if r.Method == "POST" {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if err := os.WriteFile(s.PlanPath, body, 0644); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Write([]byte("OK"))
		return
	}
}

type ChatRequest struct {
	Message string `json:"message"`
}

func (s *Server) handleChat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	// SSE Setup
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	tokens := make(chan string)

	// Tool-aware System Prompt
	systemPrompt := `You are an expert software architect and assistant for the 'repomap' project.
Your primary job is to help the user develop and refine their project plan.
You have access to the following tools:
- read_file: Read content (e.g., tasks.md)
- write_file: Update files (e.g., mark tasks as done)
- list_dir: Explore project structure

When the user asks to change the plan, use the tools to modify the relevant markdown files.
Always explain what you are doing.`

	s.Provider.SetSystemPrompt(systemPrompt)

	go func() {
		defer close(tokens)
		// Assuming GenerateStream blocks until done
		err := s.Provider.GenerateStream(req.Message, nil, tokens)
		if err != nil {
			tokens <- fmt.Sprintf("\n[Error: %v]", err)
		}
	}()

	for token := range tokens {
		fmt.Fprintf(w, "data: %s\n\n", token)
		w.(http.Flusher).Flush()
	}
}

func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	config := map[string]string{
		"provider": s.Provider.Name(),
		"plan":     s.PlanPath,
	}
	json.NewEncoder(w).Encode(config)
}
