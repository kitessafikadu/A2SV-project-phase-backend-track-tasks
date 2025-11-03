package services

import (
	"errors"
	"library-management/models"
)

type libraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int)
	BorrowBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) []models.Book
}

type Library struct{
	Books map[int]models.Book
	Members map[int]models.Member
}

func NewLibrary() *Library {
	return &Library{
		Books: make(map[int]models.Book),
		Members: make(map[int]models.Member),
	}
}

func (l *Library) AddBook(book models.Book){
	book.Status="Available"
	l.Books[book.ID]=book
}

func (l *Library) RemoveBook(bookID int){
	delete(l.Books, bookID)
}
