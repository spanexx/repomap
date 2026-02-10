package cli

import (
	"reflect"
	"testing"
)

func TestAddFlag(t *testing.T) {
	app := NewApp("test-app", "1.0.0")

	app.AddFlag("string-flag", "a string flag", "default")
	app.AddFlag("int-flag", "an int flag", 42)
	app.AddFlag("bool-flag", "a bool flag", true)
	app.AddFlag("slice-flag", "a slice flag", []string{"default"})

	// Parse with defaults
	flags, err := app.Parse([]string{})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if val := flags.GetString("string-flag"); val != "default" {
		t.Errorf("expected string-flag to be 'default', got '%s'", val)
	}
	if val := flags.GetInt("int-flag"); val != 42 {
		t.Errorf("expected int-flag to be 42, got %d", val)
	}
	if val := flags.GetBool("bool-flag"); !val {
		t.Errorf("expected bool-flag to be true, got %v", val)
	}
	if val := flags.GetStringSlice("slice-flag"); !reflect.DeepEqual(val, []string{"default"}) {
		t.Errorf("expected slice-flag to be ['default'], got %v", val)
	}
}

func TestParse(t *testing.T) {
	app := NewApp("test-app", "1.0.0")

	app.AddFlag("name", "name flag", "unknown")
	app.AddFlag("age", "age flag", 0)
	app.AddFlag("admin", "admin flag", false)
	app.AddFlag("tags", "tags flag", []string{})

	args := []string{
		"-name", "Alice",
		"-age", "30",
		"-admin",
		"-tags", "dev,ops",
		"-tags", "lead",
	}

	flags, err := app.Parse(args)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if val := flags.GetString("name"); val != "Alice" {
		t.Errorf("expected name to be 'Alice', got '%s'", val)
	}
	if val := flags.GetInt("age"); val != 30 {
		t.Errorf("expected age to be 30, got %d", val)
	}
	if val := flags.GetBool("admin"); !val {
		t.Errorf("expected admin to be true, got %v", val)
	}

	expectedTags := []string{"dev", "ops", "lead"}
	if val := flags.GetStringSlice("tags"); !reflect.DeepEqual(val, expectedTags) {
		t.Errorf("expected tags to be %v, got %v", expectedTags, val)
	}
}

func TestUnknownFlag(t *testing.T) {
	app := NewApp("test-app", "1.0.0")
	// Since we use flag.ContinueOnError, Parse should return error
	_, err := app.Parse([]string{"-unknown"})
	if err == nil {
		t.Fatal("expected error for unknown flag, got nil")
	}
}

func TestUnsupportedFlagType(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	app := NewApp("test-app", "1.0.0")
	app.AddFlag("float-flag", "unsupported", 3.14)
}
