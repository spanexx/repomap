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
