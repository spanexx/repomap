package config

import (
	"os"
	"reflect"
	"testing"
)

func TestGetEnv(t *testing.T) {
	key := "TEST_ENV_VAR"
	expected := "value"

	os.Setenv(key, expected)
	defer os.Unsetenv(key)

	if val := GetEnv(key, "default"); val != expected {
		t.Errorf("expected %s, got %s", expected, val)
	}

	if val := GetEnv("NON_EXISTENT", "default"); val != "default" {
		t.Errorf("expected default, got %s", val)
	}
}

func TestLoadEnv(t *testing.T) {
	prefix := "TEST_APP_"
	os.Setenv("TEST_APP_DEBUG", "true")
	os.Setenv("TEST_APP_OUTPUT", "json")
	os.Setenv("OTHER_VAR", "ignore")

	defer func() {
		os.Unsetenv("TEST_APP_DEBUG")
		os.Unsetenv("TEST_APP_OUTPUT")
		os.Unsetenv("OTHER_VAR")
	}()

	env := LoadEnv(prefix)

	expected := map[string]string{
		"debug": "true",
		"output": "json",
	}

	if !reflect.DeepEqual(env, expected) {
		t.Errorf("expected %v, got %v", expected, env)
	}
}
