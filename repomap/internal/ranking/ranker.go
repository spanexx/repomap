package ranking

import (
	"github.com/spanexx/agents-cli/repomap/internal/graph"
)

// Rank calculates the importance score for each file in the graph based on in-degree centrality.
// Returns a map of file paths to their normalized score (0.0 - 1.0).
func Rank(g *graph.Graph) map[string]float64 {
	scores := make(map[string]float64)
	if len(g.Nodes) == 0 {
		return scores
	}

	maxInDegree := 0
	for _, node := range g.Nodes {
		if node.InDegree > maxInDegree {
			maxInDegree = node.InDegree
		}
	}

	for path, node := range g.Nodes {
		if maxInDegree == 0 {
			scores[path] = 0.0
		} else {
			scores[path] = float64(node.InDegree) / float64(maxInDegree)
		}
	}

	return scores
}

// AssignImportance maps normalized scores to importance levels.
// > 0.7: high
// 0.3 - 0.7: medium
// < 0.3: low
func AssignImportance(scores map[string]float64) map[string]string {
	importance := make(map[string]string)
	for path, score := range scores {
		if score > 0.7 {
			importance[path] = "high"
		} else if score >= 0.3 {
			importance[path] = "medium"
		} else {
			importance[path] = "low"
		}
	}
	return importance
}
