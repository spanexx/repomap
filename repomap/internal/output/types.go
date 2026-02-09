package output

import (
	"errors"
)

// FileNode represents a single file in the repository map.
type FileNode struct {
	Path        string   `json:"path" xml:"path,attr"`
	Language    string   `json:"language" xml:"language,attr"`
	Importance  string   `json:"importance" xml:"importance,attr"`
	Rank        float64  `json:"rank" xml:"rank,attr"`
	Definitions []string `json:"definitions" xml:"definition"`
	TokenCount  int      `json:"token_count" xml:"token_count,attr"`
}

// RepoMap represents the complete repository map output.
type RepoMap struct {
	Files []*FileNode `json:"files" xml:"file"`
	XMLName struct{} `json:"-" xml:"repomap"`
}

// Validate checks if the FileNode is valid.
func (n *FileNode) Validate() error {
	if n.Path == "" {
		return errors.New("file path cannot be empty")
	}
	return nil
}
