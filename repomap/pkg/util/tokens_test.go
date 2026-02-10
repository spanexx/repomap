package util

import "testing"

func TestCountTokens(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"", 0},
		{"a", 0}, // 1/4 = 0
		{"abcd", 1},
		{"12345678", 2},
		{"Hello World!", 3}, // 12 chars / 4 = 3
	}

	for _, tt := range tests {
		if got := CountTokens(tt.input); got != tt.expected {
			t.Errorf("CountTokens(%q) = %d; want %d", tt.input, got, tt.expected)
		}
	}
}
