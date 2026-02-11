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

	"github.com/spanexx/agents-cli/repomap/internal/session"
	"github.com/spanexx/agents-cli/repomap/pkg/adapter"
	"github.com/spanexx/agents-cli/repomap/pkg/providers/gemini_cli"
	"github.com/spanexx/agents-cli/repomap/pkg/providers/qodercli"
	"github.com/spanexx/agents-cli/repomap/pkg/providers/qwen_cli"
)

type Server struct {
	Port        string
	Provider    adapter.Provider
	PlanPath    string
	Sessions    *session.Manager
	RepoMapData interface{}
}

func New(port string, planPath string, repoMapData interface{}) *Server {
	// Initialize session manager in .repomap/sessions
	home, _ := os.UserHomeDir()
	sessionDir := filepath.Join(home, ".repomap", "sessions")
	sm, _ := session.NewManager(sessionDir)

	return &Server{
		Port:        port,
		PlanPath:    planPath,
		Sessions:    sm,
		RepoMapData: repoMapData,
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
	http.HandleFunc("/api/repomap", s.handleRepoMap)
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

func (s *Server) handleRepoMap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if s.RepoMapData == nil {
		http.Error(w, "No repomap data available", 404)
		return
	}

	json.NewEncoder(w).Encode(s.RepoMapData)
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

type ChatContext struct {
	SelectedNode string `json:"selectedNode"`
	ViewMode     string `json:"viewMode"`
}

type ChatRequest struct {
	Message   string       `json:"message"`
	SessionID string       `json:"sessionId"`
	Context   *ChatContext `json:"context,omitempty"`
}

func (s *Server) handleChat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "GET" {
		sessionID := r.URL.Query().Get("sessionId")
		if sessionID == "" {
			http.Error(w, "Session ID required", 400)
			return
		}

		sess, err := s.Sessions.GetSession(sessionID)
		if err != nil {
			http.Error(w, "Session not found", 404)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sess.Messages)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	// Session Handling
	if req.SessionID != "" {
		// Ensure session exists
		_, err := s.Sessions.GetSession(req.SessionID)
		if err != nil {
			// Create it if not found
			if err := s.Sessions.CreateSessionWithID(req.SessionID); err != nil {
				fmt.Printf("Failed to create session: %v\n", err)
			}
		}

		// Save User Message
		s.Sessions.AddMessage(req.SessionID, "user", req.Message)
	}

	// SSE Setup
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	tokens := make(chan string)

	// Context Construction
	var fullContext string
	if req.SessionID != "" {
		history, _ := s.Sessions.FormatHistory(req.SessionID) // Get full history including recent user msg
		fullContext = history
		// Note: The history already includes the current message we just added!
		// But duplicate? No, FormatHistory gets all messages.
		// So we should NOT append req.Message again to prompt if we pass full history?
		// Providers usually take "Prompt".
		// If we pass Full History as Prompt, the model might get confused about what to answer.
		// Standard pattern: System Prompt + History (Context) + Last Message.
		// Our `FormatHistory` returns everything.
		// Let's just pass `fullContext` as the prompt.
	} else {
		fullContext = req.Message
	}

	// Tool-aware System Prompt
	// Tool-aware System Prompt
	// Tool-aware System Prompt
	var dynamicContext string
	if req.Context != nil {
		dynamicContext = fmt.Sprintf("\nCurrent Context:\n- View Mode: %s\n- Selected Node: %s\n", 
			req.Context.ViewMode, req.Context.SelectedNode)
	}

	systemPrompt := `You are an expert architect for 'repomap'. Help the user refine their project plan.
Tools: read_file, write_file, list_dir.
Use tools to modify plan files when requested. Always explain your actions.
Be concise and thoughtful.` + dynamicContext

	s.Provider.SetSystemPrompt(systemPrompt)

	go func() {
		defer close(tokens)

		// Stream generation
		err := s.Provider.GenerateStream(fullContext, nil, tokens)
		if err != nil {
			tokens <- fmt.Sprintf("\n[Error: %v]", err)
		}

		// Capture full response for history
		// We need to intercept the tokens to build the full string string.
		// But `GenerateStream` writes to channel directly.
		// We can wrap the channel or just accumulate logic here?
		// Wait, GenerateStream blocks until done.
		// We can't see the tokens here easily unless we change the architecture
		// to return collected response OR we wrap the channel.
		// Let's assume we can't save the assistant message EASILY without changing architecture.
		// BUT wait, we need to save it for history!
		// Refactoring `GenerateStream` is risky right now.
		// Alternative: Wrap the token channel.
	}()

	// We need to intercept tokens to save to history
	// We need to intercept tokens to save to history
	var assistantResponse strings.Builder

	type StreamResponse struct {
		Token string `json:"token"`
	}

	for token := range tokens {
		assistantResponse.WriteString(token)

		resp := StreamResponse{Token: token}
		jsonBytes, err := json.Marshal(resp)
		if err == nil {
			fmt.Fprintf(w, "data: %s\n\n", jsonBytes)
			w.(http.Flusher).Flush()
		}
	}

	// Save Assistant Message
	if req.SessionID != "" {
		s.Sessions.AddMessage(req.SessionID, "assistant", assistantResponse.String())
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
