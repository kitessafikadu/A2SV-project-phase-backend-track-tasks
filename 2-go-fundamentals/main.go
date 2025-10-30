package main

import (
	"bufio"
	"fmt"
	"go-fundamentals/palindrome"
	"go-fundamentals/wordfreq"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n=== Go Fundamentals CLI ===")
		fmt.Println("1. Word Frequency Count")
		fmt.Println("2. Palindrome Check")
		fmt.Println("3. Exit")
		fmt.Print("Choose an option (1-3): ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			fmt.Print("\nEnter text for word frequency count:\n> ")
			text, _ := reader.ReadString('\n')
			result := wordfreq.WordFrequencyCount(text)
			fmt.Println("\n🔍 Word Frequency Result:")
			for word, count := range result {
				fmt.Printf("%s: %d\n", word, count)
			}

		case "2":
			fmt.Print("\nEnter text to check palindrome:\n> ")
			text, _ := reader.ReadString('\n')
			if palindrome.IsPalindrome(text) {
				fmt.Println("✅ It's a palindrome!")
			} else {
				fmt.Println("❌ Not a palindrome.")
			}

		case "3":
			fmt.Println("👋 Exiting program...")
			return

		default:
			fmt.Println("⚠️ Invalid choice, please try again.")
		}
	}
}
