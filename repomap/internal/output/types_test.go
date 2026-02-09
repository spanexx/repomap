package output

import (
	"testing"
)

func TestFileNode_Validate(t *testing.T) {
	tests := []struct {
		name    string
		node    *FileNode
		wantErr bool
	}{
		{
			name: "Valid Node",
			node: &FileNode{
				Path:       "main.go",
				Language:   "go",
				Importance: "high",
				Rank:       0.9,
			},
			wantErr: false,
		},
		{
			name: "Empty Path",
			node: &FileNode{
				Path: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.node.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("FileNode.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
