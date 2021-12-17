package str

import "strings"

// Case in-sensitive search: Returns true if s constains substr. False otherwise.
func CaseInsensitiveSearch(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
