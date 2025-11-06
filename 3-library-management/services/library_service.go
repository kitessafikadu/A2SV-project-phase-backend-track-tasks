package services

import (
	"errors"
	"library_management/models"
)

type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int)
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) []models.Book
}

type Library struct {
	Books   map[int]models.Book
	Members map[int]models.Member
}

func NewLibrary() *Library {
	return &Library{
		Books:   make(map[int]models.Book),
		Members: make(map[int]models.Member),
	}
}

func (l *Library) AddBook(book models.Book) {
	book.Status = "Available"
	l.Books[book.ID] = book
}

func (l *Library) RemoveBook(bookID int) {
	delete(l.Books, bookID)
}

func (l *Library) BorrowBook(bookID int, memberID int) error {
	book, ok := l.Books[bookID]
	if !ok {
		return errors.New("book not found")
	}
	if book.Status == "Borrowed" {
		return errors.New("book already borrowed")
	}

	member, ok := l.Members[memberID]
	if !ok {
		return errors.New("member not found")
	}

	book.Status = "Borrowed"
	member.BorrowedBooks = append(member.BorrowedBooks, book)
	l.Members[memberID] = member
	l.Books[bookID] = book
	return nil
}

func (l *Library) ReturnBook(bookID int, memberID int) error {
	member, ok := l.Members[memberID]
	if !ok {
		return errors.New("member not found")
	}

	book, ok := l.Books[bookID]
	if !ok {
		return errors.New("book not found")
	}

	found := false
	newBorrowedBooks := []models.Book{}
	for _, b := range member.BorrowedBooks {
		if b.ID == bookID {
			found = true
		} else {
			newBorrowedBooks = append(newBorrowedBooks, b)
		}
	}

	if !found {
		return errors.New("book not borrowed by this member")
	}

	book.Status = "Available"
	member.BorrowedBooks = newBorrowedBooks
	l.Books[bookID] = book
	l.Members[memberID] = member
	return nil
}

func (l *Library) ListAvailableBooks() []models.Book {
	var available []models.Book
	for _, book := range l.Books {
		if book.Status == "Available" {
			available = append(available, book)
		}
	}
	return available
}

func (l *Library) ListBorrowedBooks(memberID int) []models.Book {
	member, ok := l.Members[memberID]
	if !ok {
		return []models.Book{}
	}
	return member.BorrowedBooks
}
