package util

import (
	"path/filepath"
	"testing"
)

func TestNormalizePath(t *testing.T) {
	input := filepath.Join("foo", "bar", "..", "baz")
	expected := "foo/baz"

	if result := NormalizePath(input); result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestMakeRelative(t *testing.T) {
	base := "/home/user/project"
	target := "/home/user/project/src/main.go"
	expected := "src/main.go"

	// Note: filepath.Rel requires absolute paths on some platforms to work predictably?
	// Or just cleaner paths.
	// For testing, let's use relative paths if running in unknown env,
	// or mock the OS separator behavior?
	// We'll trust filepath behavior but ensure ToSlash works.

	// Let's use simple relative paths that work everywhere
	base = "project"
	target = filepath.Join("project", "src", "file.go")
	expected = "src/file.go"

	if result := MakeRelative(base, target); result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestHasExtension(t *testing.T) {
	if !HasExtension("file.go", []string{".go", ".js"}) {
		t.Error("expected true for .go extension")
	}

	if HasExtension("file.txt", []string{".go"}) {
		t.Error("expected false for .txt extension")
	}
}
