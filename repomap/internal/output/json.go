package output

import (
	"encoding/json"
)

// RenderJSON converts a list of FileNodes into a JSON string.
func RenderJSON(nodes []*FileNode) (string, error) {
	repoMap := RepoMap{
		Files: nodes,
	}

	output, err := json.MarshalIndent(repoMap, "", "  ")
	if err != nil {
		return "", err
	}

	return string(output), nil
}
