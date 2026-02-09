package output

import (
	"encoding/xml"
)

// RenderXML converts a list of FileNodes into an XML string.
// maxTokens specifies the maximum number of tokens allowed in the output.
// If maxTokens is 0, no limit is applied.
func RenderXML(nodes []*FileNode, maxTokens int) (string, error) {
	var includedNodes []*FileNode
	currentTokens := 0

	// Add header tokens approx count (not precise but safe buffer)
	// <?xml version="1.0" encoding="UTF-8"?>\n<repomap>\n</repomap> is ~50 chars -> 12 tokens
	currentTokens += 12

	for _, node := range nodes {
		nodeCost := node.TokenCount
		if nodeCost == 0 {
			// Estimate if missing
			for _, def := range node.Definitions {
				nodeCost += CountTokens(def)
			}
			// Add path/metadata cost
			nodeCost += CountTokens(node.Path) + 5 // +5 for attributes overhead
		}

		if maxTokens > 0 && currentTokens+nodeCost > maxTokens {
			// Budget exceeded
			// Add a truncation notice node
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

	output, err := xml.MarshalIndent(repoMap, "", "  ")
	if err != nil {
		return "", err
	}

	// Add XML header
	return xml.Header + string(output), nil
}
