package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Message represents a single turn in the conversation.
type Message struct {
	Role      string    `json:"role"` // "user" or "assistant" (or "model")
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// Session represents a conversation history.
type Session struct {
	ID        string    `json:"id"`
	Messages  []Message `json:"messages"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Manager handles session storage and retrieval.
type Manager struct {
	storageDir string
	mu         sync.RWMutex
}

// NewManager creates a new session manager using the given storage directory.
func NewManager(storageDir string) (*Manager, error) {
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create session storage dir: %w", err)
	}
	return &Manager{
		storageDir: storageDir,
	}, nil
}

// CreateSession initializes a new session.
func (m *Manager) CreateSession() (*Session, error) {
	id := uuid.New().String()
	session := &Session{
		ID:        id,
		Messages:  []Message{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := m.SaveSession(session); err != nil {
		return nil, err
	}
	return session, nil
}

// CreateSessionWithID initializes a new session with a specific ID.
func (m *Manager) CreateSessionWithID(id string) error {
	session := &Session{
		ID:        id,
		Messages:  []Message{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := m.SaveSession(session); err != nil {
		return err
	}
	return nil
}

// GetSession loads a session by ID.
func (m *Manager) GetSession(id string) (*Session, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	path := filepath.Join(m.storageDir, id+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("failed to parse session: %w", err)
	}
	return &session, nil
}

// SaveSession persists the session to disk.
func (m *Manager) SaveSession(session *Session) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	session.UpdatedAt = time.Now()
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	path := filepath.Join(m.storageDir, session.ID+".json")
	return os.WriteFile(path, data, 0644)
}

// AddMessage appends a message to the session and saves it.
func (m *Manager) AddMessage(sessionID string, role, content string) error {
	s, err := m.GetSession(sessionID)
	if err != nil {
		return err
	}

	s.Messages = append(s.Messages, Message{
		Role:      role,
		Content:   content,
		Timestamp: time.Now(),
	})

	return m.SaveSession(s)
}

// ClearSession clears the message history for a session.
func (m *Manager) ClearSession(sessionID string) error {
	s, err := m.GetSession(sessionID)
	if err != nil {
		return err
	}
	s.Messages = []Message{}
	return m.SaveSession(s)
}

// FormatHistory returns the conversation history as a formatted string context.
func (m *Manager) FormatHistory(sessionID string) (string, error) {
	s, err := m.GetSession(sessionID)
	if err != nil {
		return "", err
	}

	var sb string
	for _, msg := range s.Messages {
		role := "User"
		if msg.Role == "assistant" || msg.Role == "model" {
			role = "Assistant"
		}
		sb += fmt.Sprintf("--- %s ---\n%s\n\n", role, msg.Content)
	}
	return sb, nil
}
