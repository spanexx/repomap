package output

import (
	"testing"
)

func TestCountTokens(t *testing.T) {
	tests := []struct {
		text string
		want int
	}{
		{"", 0},
		{"a", 0}, // 1/4 = 0
		{"abcd", 1},
		{"12345678", 2},
		{"Hello World!", 3}, // 12 / 4 = 3
	}

	for _, tt := range tests {
		got := CountTokens(tt.text)
		if got != tt.want {
			t.Errorf("CountTokens(%q) = %d, want %d", tt.text, got, tt.want)
		}
	}
}
