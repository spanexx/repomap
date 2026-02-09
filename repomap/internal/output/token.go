package output

// CountTokens estimates the number of tokens in a string.
// For the MVP, it uses a simple heuristic: length of string divided by 4.
func CountTokens(text string) int {
	return len(text) / 4
}
