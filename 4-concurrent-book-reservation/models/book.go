package models

import "sync"

type Book struct {
	ID         int
	Title      string
	Author     string

	// Protected by LibraryManager's mutex
	Available  bool
	ReservedBy int 
	BorrowedBy int 
	
	// optional per-book mutex if you want more-granular locking
	mu sync.Mutex
}
