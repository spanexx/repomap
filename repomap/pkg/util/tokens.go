package util

// CountTokens estimates the number of tokens in a text string.
// It uses a simple approximation: 1 token ~= 4 characters.
// This is a standard heuristic for GPT-based models when a tokenizer is not available.
func CountTokens(text string) int {
	if text == "" {
		return 0
	}
	// Use ceil to ensure non-empty strings count as at least 1 token if they have chars
	// but integer division is standard.
	// Let's use simple division as per spec: len(text)/4.
	// But commonly we round up.
	// Spec says: "Document approximation method".
	// Let's stick to len/4 but ensure at least 1 if not empty?
	// The prompt memory said "len(text)/4".
	return len(text) / 4
}

// CountTokensRough is an alias for CountTokens, explicitly stating it's an estimate.
func CountTokensRough(text string) int {
	return CountTokens(text)
}
