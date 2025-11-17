package main

import (
	"fmt"
	"concurrent-book-reservation/controllers"
	"concurrent-book-reservation/models"
	"concurrent-book-reservation/services"
	"time"
)

func main() {
	// Initialize library
	lib := services.NewInMemoryLibrary(20)
	lib.AddBook(&models.Book{ID: 1, Title: "Concurrency in Go", Author: "K. Cox", Available: true})
	lib.AddBook(&models.Book{ID: 2, Title: "Clean Code", Author: "R. Martin", Available: true})

	// Start worker
	lib.StartReservationWorker()
	defer lib.Shutdown()

	// Controller
	ctrl := controllers.NewLibraryController(lib)

	fmt.Println("=== Concurrent Reservation Simulation ===")
	// Simulate multiple members attempting to reserve the same book (ID=1) almost at once
	memberCount := 5
	for i := 1; i <= memberCount; i++ {
		memberID := i
		// simulate slight staggering
		go ctrl.RequestReserve(1, memberID)
	}

	// Also simulate some members requesting for another book
	ctrl.RequestReserve(2, 10)
	ctrl.RequestReserve(2, 11)

	// Optionally show explicit borrow after some delay for one of the members
	// (This demonstrates explicit borrow competing with auto-cancel)
	ctrl.RequestExplicitBorrowAfter(1, 2, 3*time.Second) // member 2 tries to explicitly borrow after 3s

	// Give enough time for async operations, borrows, and auto-cancellations to print.
	time.Sleep(8 * time.Second)

	// Print final state
	b1, _ := lib.GetBook(1)
	b2, _ := lib.GetBook(2)
	fmt.Println("----- Final state -----")
	fmt.Printf("Book 1: ID=%d Title=%s Available=%v ReservedBy=%d BorrowedBy=%d\n", b1.ID, b1.Title, b1.Available, b1.ReservedBy, b1.BorrowedBy)
	fmt.Printf("Book 2: ID=%d Title=%s Available=%v ReservedBy=%d BorrowedBy=%d\n", b2.ID, b2.Title, b2.Available, b2.ReservedBy, b2.BorrowedBy)
}
