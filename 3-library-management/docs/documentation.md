# Library Management System (Go)

## Overview

This console-based application demonstrates a simple library management system built in Go using structs, interfaces, maps, and slices.

## Features

- Add, remove, borrow, and return books
- Track available and borrowed books
- Support for multiple members

## Architecture

- **models/**: Defines `Book` and `Member` structs.
- **services/**: Implements `LibraryManager` interface and business logic.
- **controllers/**: Handles user input/output.
- **main.go**: Entry point.

## Run Instructions

1. Initialize module:

```
go mod init library-management
go mod tidy
```

2. Run the program

```
go run main.go
```

## Example Usage

1. Add Book
2. Remove Book
3. Borrow Book
4. Return Book
5. List Available Books
6. List Borrowed Books
7. Exit
