package output

import (
	"encoding/xml"
)

// RenderXML converts a list of FileNodes into an XML string.
func RenderXML(nodes []*FileNode) (string, error) {
	repoMap := RepoMap{
		Files: nodes,
	}

	output, err := xml.MarshalIndent(repoMap, "", "  ")
	if err != nil {
		return "", err
	}

	// Add XML header
	return xml.Header + string(output), nil
}
