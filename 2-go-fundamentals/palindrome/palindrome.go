package palindrome

import (
	"regexp"
	"strings"
)

// IsPalindrome checks if a string is a palindrome ignoring punctuation and case.
func IsPalindrome(s string) bool {
	s = strings.ToLower(s)
	re := regexp.MustCompile(`[^a-z0-9]`)
	clean := re.ReplaceAllString(s, "")

	n := len(clean)
	for i := 0; i < n/2; i++ {
		if clean[i] != clean[n-1-i] {
			return false
		}
	}
	return true
}
