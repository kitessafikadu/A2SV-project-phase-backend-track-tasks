package controllers

import (
	"fmt"
	"concurrent-book-reservation/services"
	"time"
)

// Very small "controller" layer to simulate incoming client calls.
// In a real app you'd wire HTTP endpoints here and call service.ReserveBook from handlers.

type LibraryController struct {
	service services.LibraryManager
}

func NewLibraryController(s services.LibraryManager) *LibraryController {
	return &LibraryController{service: s}
}

// Simulate a reservation request from a member
func (c *LibraryController) RequestReserve(bookID, memberID int) {
	go func() {
		// calling ReserveBook will send a request to the reservation queue and wait for result
		err := c.service.ReserveBook(bookID, memberID)
		if err != nil {
			fmt.Printf("[Controller] member %d: reserve book %d -> ERROR: %v\n", memberID, bookID, err)
			return
		}
		fmt.Printf("[Controller] member %d: reservation request accepted for book %d\n", memberID, bookID)
	}()
}

// Optionally simulate a member explicitly calling borrow (separate from asynchronous borrow)
// This demonstrates concurrent actions: send an explicit borrow after some delay.
func (c *LibraryController) RequestExplicitBorrowAfter(bookID, memberID int, after time.Duration) {
	go func() {
		time.Sleep(after)
		err := c.service.ProcessBorrow(bookID, memberID)
		if err != nil {
			fmt.Printf("[Controller] member %d explicit borrow book %d -> ERROR: %v\n", memberID, bookID, err)
			return
		}
		fmt.Printf("[Controller] member %d explicit borrow book %d -> SUCCESS\n", memberID, bookID)
	}()
}
