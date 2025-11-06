package controllers

import (
	"fmt"
	"library-management/models"
	"library-management/services"
)

func StartConsoleApp() {
	library := services.NewLibrary()

	// Add some sample members for testing
	library.Members[1] = models.Member{ID: 1, Name: "Alice"}
	library.Members[2] = models.Member{ID: 2, Name: "Bob"}

	for {
		fmt.Println("\n===== Library Management System =====")
		fmt.Println("1. Add Book")
		fmt.Println("2. Remove Book")
		fmt.Println("3. Borrow Book")
		fmt.Println("4. Return Book")
		fmt.Println("5. List Available Books")
		fmt.Println("6. List Borrowed Books")
		fmt.Println("7. Exit")
		fmt.Print("Enter your choice: ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			var id int
			var title, author string
			fmt.Print("Enter book ID: ")
			fmt.Scanln(&id)
			fmt.Print("Enter book title: ")
			fmt.Scanln(&title)
			fmt.Print("Enter book author: ")
			fmt.Scanln(&author)
			library.AddBook(models.Book{ID: id, Title: title, Author: author})
			fmt.Println("✅ Book added successfully.")

		case 2:
			var id int
			fmt.Print("Enter book ID to remove: ")
			fmt.Scanln(&id)
			library.RemoveBook(id)
			fmt.Println("🗑️ Book removed successfully.")

		case 3:
			var bookID, memberID int
			fmt.Print("Enter book ID: ")
			fmt.Scanln(&bookID)
			fmt.Print("Enter member ID: ")
			fmt.Scanln(&memberID)
			if err := library.BorrowBook(bookID, memberID); err != nil {
				fmt.Println("❌", err)
			} else {
				fmt.Println("📚 Book borrowed successfully.")
			}

		case 4:
			var bookID, memberID int
			fmt.Print("Enter book ID: ")
			fmt.Scanln(&bookID)
			fmt.Print("Enter member ID: ")
			fmt.Scanln(&memberID)
			if err := library.ReturnBook(bookID, memberID); err != nil {
				fmt.Println("❌", err)
			} else {
				fmt.Println("✅ Book returned successfully.")
			}

		case 5:
			fmt.Println("\n📖 Available Books:")
			for _, book := range library.ListAvailableBooks() {
				fmt.Printf("[%d] %s by %s\n", book.ID, book.Title, book.Author)
			}

		case 6:
			var memberID int
			fmt.Print("Enter member ID: ")
			fmt.Scanln(&memberID)
			fmt.Printf("\n👤 Borrowed Books for Member %d:\n", memberID)
			for _, book := range library.ListBorrowedBooks(memberID) {
				fmt.Printf("[%d] %s by %s\n", book.ID, book.Title, book.Author)
			}

		case 7:
			fmt.Println("👋 Exiting... Goodbye!")
			return

		default:
			fmt.Println("❌ Invalid choice. Try again.")
		}
	}
}
