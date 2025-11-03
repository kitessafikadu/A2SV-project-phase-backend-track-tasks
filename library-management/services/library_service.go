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