package ranking

import (
	"testing"

	"github.com/spanexx/agents-cli/repomap/internal/graph"
)

func TestRank(t *testing.T) {
	// Create a mock graph
	// Nodes: A, B, C, D
	// Edges:
	// A -> B
	// A -> C
	// B -> C
	// D -> (none)

	// In-degrees:
	// A: 0
	// B: 1 (from A)
	// C: 2 (from A, B)
	// D: 0

	// Max In-Degree: 2

	// Expected Scores:
	// A: 0.0
	// B: 0.5
	// C: 1.0
	// D: 0.0

	g := &graph.Graph{
		Nodes: map[string]*graph.Node{
			"A": {Path: "A", InDegree: 0},
			"B": {Path: "B", InDegree: 1},
			"C": {Path: "C", InDegree: 2},
			"D": {Path: "D", InDegree: 0},
		},
		Edges: map[string][]string{
			"A": {"B", "C"},
			"B": {"C"},
		},
	}

	scores := Rank(g)

	expected := map[string]float64{
		"A": 0.0,
		"B": 0.5,
		"C": 1.0,
		"D": 0.0,
	}

	for path, want := range expected {
		got, ok := scores[path]
		if !ok {
			t.Errorf("Score for %s missing", path)
			continue
		}
		if got != want {
			t.Errorf("Score(%s) = %f, want %f", path, got, want)
		}
	}
}

func TestRank_Empty(t *testing.T) {
	g := &graph.Graph{
		Nodes: map[string]*graph.Node{},
	}
	scores := Rank(g)
	if len(scores) != 0 {
		t.Errorf("Expected empty scores, got %d", len(scores))
	}
}

func TestRank_ZeroEdges(t *testing.T) {
	g := &graph.Graph{
		Nodes: map[string]*graph.Node{
			"A": {Path: "A", InDegree: 0},
			"B": {Path: "B", InDegree: 0},
		},
	}
	scores := Rank(g)

	// Max degree is 0, so all scores should be 0.0
	if scores["A"] != 0.0 || scores["B"] != 0.0 {
		t.Errorf("Expected 0.0 scores for disconnected graph")
	}
}
