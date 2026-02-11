package server

import (
	"testing"
)

func TestNew(t *testing.T) {
	port := "9090"
	planPath := "test_plan.json"
	srv := New(port, planPath, nil)

	if srv == nil {
		t.Fatal("New returned nil")
	}

	if srv.Port != port {
		t.Errorf("expected port %s, got %s", port, srv.Port)
	}
}
