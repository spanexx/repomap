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
	Imports     []string `json:"imports,omitempty" xml:"import,omitempty"`
	TokenCount  int       `json:"token_count" xml:"token_count,attr"`
	// Planning Features
	Status   string    `json:"status,omitempty" xml:"status,attr,omitempty"`
	Intent   string    `json:"intent,omitempty" xml:"intent,omitempty"`
	Issues   []Issue   `json:"issues,omitempty" xml:"issue,omitempty"`
	Comments []Comment `json:"comments,omitempty" xml:"comment,omitempty"`
}

type Issue struct {
	Type        string `json:"type" xml:"type,attr"`
	Description string `json:"description" xml:"description"`
	Severity    string `json:"severity" xml:"severity,attr"`
}

type Comment struct {
	User string `json:"user" xml:"user,attr"`
	Text string `json:"text" xml:"text"`
}

// RepoMap represents the complete repository map output.
type RepoMap struct {
	Files   []*FileNode `json:"files" xml:"file"`
	XMLName struct{}    `json:"-" xml:"repomap"`
}

// Validate checks if the FileNode is valid.
func (n *FileNode) Validate() error {
	if n.Path == "" {
		return errors.New("file path cannot be empty")
	}
	return nil
}
