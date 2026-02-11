package cli

import (
	"testing"
)

func TestNewApp(t *testing.T) {
	appName := "test-app"
	appVersion := "1.0.0"

	app := NewApp(appName, appVersion)

	if app == nil {
		t.Fatal("NewApp returned nil")
	}

	if app.Name != appName {
		t.Errorf("expected Name to be %s, got %s", appName, app.Name)
	}

	if app.Version != appVersion {
		t.Errorf("expected Version to be %s, got %s", appVersion, app.Version)
	}
}

func TestAppMethods(t *testing.T) {
	app := NewApp("test-app", "1.0.0")

	app.SetDescription("A test application")
	if app.Description != "A test application" {
		t.Errorf("App description not set correctly")
	}

	app.SetUsage("test-app [options]")
	if app.UsageText != "test-app [options]" {
		t.Errorf("App usage text not set correctly")
	}

	app.AddExample("test-app --verbose")
	if len(app.Examples) != 1 || app.Examples[0] != "test-app --verbose" {
		t.Errorf("App example not added correctly")
	}
}
