package wordfreq

import (
	"regexp"
	"strings"
)

// WordFrequencyCount returns a map of word -> frequency from the input string.
func WordFrequencyCount(text string) map[string]int {
	text = strings.ToLower(text)
	re := regexp.MustCompile(`[^\w\s]`)
	cleanText := re.ReplaceAllString(text, "")
	words := strings.Fields(cleanText)

	frequency := make(map[string]int)
	for _, word := range words {
		frequency[word]++
	}
	return frequency
}
