package util

import (
	"reflect"
	"testing"
)

func TestFilterByExtension(t *testing.T) {
	paths := []string{
		"main.go",
		"README.md",
		"app.js",
		"test.go",
	}

	exts := []string{".go"}
	expected := []string{"main.go", "test.go"}

	result := FilterByExtension(paths, exts)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}

	// Case insensitivity
	result = FilterByExtension([]string{"test.GO"}, []string{".go"})
	if len(result) != 1 {
		t.Error("expected case insensitive match")
	}
}

func TestExcludeByPattern(t *testing.T) {
	paths := []string{
		"main.go",
		"main_test.go",
		"vendor/lib.go",
		"README.md",
	}

	patterns := []string{"*_test.go", "vendor/*"}
	// vendor/* matches "vendor/lib.go" because it contains a slash?
	// filepath.Match: "The pattern syntax is... match does not support **"
	// "vendor/*" matches "vendor/lib.go"

	expected := []string{"main.go", "README.md"}

	result := ExcludeByPattern(paths, patterns)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}
