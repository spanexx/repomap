package output

import (
	"encoding/json"
)

// RenderJSON converts a list of FileNodes into a JSON string.
// maxTokens specifies the maximum number of tokens allowed in the output.
// If maxTokens is 0, no limit is applied.
func RenderJSON(nodes []*FileNode, maxTokens int) (string, error) {
	var includedNodes []*FileNode
	currentTokens := 0

	// Add wrapper tokens approx count
	// {"files":[]} is ~12 chars -> 3 tokens
	currentTokens += 3

	for _, node := range nodes {
		nodeCost := node.TokenCount
		if nodeCost == 0 {
			// Estimate if missing
			for _, def := range node.Definitions {
				nodeCost += CountTokens(def)
			}
			nodeCost += CountTokens(node.Path) + 5
		}

		if maxTokens > 0 && currentTokens+nodeCost > maxTokens {
			// Budget exceeded
			truncatedNode := &FileNode{
				Path:        "...",
				Importance:  "truncated",
				Definitions: []string{"Output truncated due to token limit"},
			}
			includedNodes = append(includedNodes, truncatedNode)
			break
		}

		includedNodes = append(includedNodes, node)
		currentTokens += nodeCost
	}

	repoMap := RepoMap{
		Files: includedNodes,
	}

	output, err := json.MarshalIndent(repoMap, "", "  ")
	if err != nil {
		return "", err
	}

	return string(output), nil
}
