// Word frequency count in Go
package main

import (
	"fmt"
	"regexp"
	"strings"
)


func WordFrequencyCount(text string) map[string]int {
	text = strings.ToLower(text)

	// Remove punctuation using regex
	re := regexp.MustCompile(`[^\w\s]`)
	cleanText := re.ReplaceAllString(text, "")

	// Split into words
	words := strings.Fields(cleanText)

	// Count word frequencies
	frequency := make(map[string]int)
	for _, word := range words {
		frequency[word]++
	}

	return frequency
}

func main() {
	text := "Hello, hello! How are you? You are doing great, right?"
	fmt.Println(WordFrequencyCount(text))
}
