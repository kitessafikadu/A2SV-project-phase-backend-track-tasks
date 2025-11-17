package services

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"library_management/models"
)

// LibraryManager defines the public API.
type LibraryManager interface {
	AddBook(b *models.Book)
	GetBook(id int) (*models.Book, error)
	ReserveBook(bookID int, memberID int) error
	ProcessBorrow(bookID int, memberID int) error
	StartReservationWorker()
	Shutdown()
}

// reservation request that flows through channel
type ReservationRequest struct {
	BookID   int
	MemberID int
	Resp     chan error
}

// InMemoryLibrary is a concrete implementation
type InMemoryLibrary struct {
	books    map[int]*models.Book
	members  map[int]*models.Member
	mu       sync.Mutex // protects books and members maps & their fields
	reqCh    chan ReservationRequest
	quit     chan struct{}
	wg       sync.WaitGroup
	workerOn bool
}

func NewInMemoryLibrary(buffer int) *InMemoryLibrary {
	l := &InMemoryLibrary{
		books:   make(map[int]*models.Book),
		members: make(map[int]*models.Member),
		reqCh:   make(chan ReservationRequest, buffer),
		quit:    make(chan struct{}),
	}
	return l
}

func (l *InMemoryLibrary) AddBook(b *models.Book) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.books[b.ID] = b
}

func (l *InMemoryLibrary) GetBook(id int) (*models.Book, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	b, ok := l.books[id]
	if !ok {
		return nil, errors.New("book not found")
	}
	// return pointer (caller should not modify without library mutex)
	return b, nil
}

// ReserveBook sends a reservation request into the worker queue and waits for worker's result.
// It returns an error if reservation couldn't be made (book missing or already reserved).
func (l *InMemoryLibrary) ReserveBook(bookID int, memberID int) error {
	if !l.workerOn {
		return errors.New("reservation worker not started")
	}
	resp := make(chan error, 1)
	req := ReservationRequest{
		BookID:   bookID,
		MemberID: memberID,
		Resp:     resp,
	}
	select {
	case l.reqCh <- req:
		// wait for worker to respond
		err := <-resp
		return err
	case <-time.After(2 * time.Second):
		return errors.New("failed to queue reservation")
	}
}

// ProcessBorrow is the action that finalizes the borrow.
// It should be called under l.mu to avoid races, but here we lock inside.
func (l *InMemoryLibrary) ProcessBorrow(bookID int, memberID int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	b, ok := l.books[bookID]
	if !ok {
		return errors.New("book not found")
	}
	// If it's reserved by memberID and not already borrowed -> complete borrow
	if b.ReservedBy != memberID {
		return errors.New("book not reserved by this member")
	}
	if b.BorrowedBy != 0 {
		return errors.New("book already borrowed")
	}
	// finalize borrow
	b.BorrowedBy = memberID
	// clear reservation
	b.ReservedBy = 0
	b.Available = false
	return nil
}

// StartReservationWorker launches worker goroutines to process requests concurrently.
func (l *InMemoryLibrary) StartReservationWorker() {
	l.mu.Lock()
	if l.workerOn {
		l.mu.Unlock()
		return
	}
	l.workerOn = true
	l.mu.Unlock()

	// single goroutine that dispatches each request to its own goroutine for concurrency,
	// allowing many reservations to be processed in parallel while queueing remains ordered.
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		for {
			select {
			case req := <-l.reqCh:
				// dispatch to handler goroutine
				l.wg.Add(1)
				go func(r ReservationRequest) {
					defer l.wg.Done()
					r.Resp <- l.handleReservation(r.BookID, r.MemberID)
				}(req)
			case <-l.quit:
				return
			}
		}
	}()
}

// handleReservation performs the reservation logic, starts auto-cancel timer and
// starts an asynchronous borrow attempt (simulated) that's racing with auto-cancel.
func (l *InMemoryLibrary) handleReservation(bookID, memberID int) error {
	// Step 1: check and set reservation atomically
	l.mu.Lock()
	b, ok := l.books[bookID]
	if !ok {
		l.mu.Unlock()
		return errors.New("book not found")
	}
	if b.ReservedBy != 0 {
		l.mu.Unlock()
		return fmt.Errorf("book %d already reserved by member %d", bookID, b.ReservedBy)
	}
	if b.BorrowedBy != 0 {
		l.mu.Unlock()
		return fmt.Errorf("book %d already borrowed by member %d", bookID, b.BorrowedBy)
	}
	// reserve
	b.ReservedBy = memberID
	// mark as not available for others
	b.Available = false
	l.mu.Unlock()

	// Start auto-cancel timer goroutine (5 seconds)
	cancelCh := make(chan struct{})
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		timer := time.NewTimer(5 * time.Second)
		select {
		case <-timer.C:
			// timer expired -> auto-cancel if still reserved and not borrowed
			l.mu.Lock()
			defer l.mu.Unlock()
			cur := l.books[bookID]
			if cur == nil {
				return
			}
			if cur.ReservedBy == memberID && cur.BorrowedBy == 0 {
				// auto-unreserve
				cur.ReservedBy = 0
				cur.Available = true
				fmt.Printf("[AutoCancel] reservation for book %d by member %d cancelled (timeout)\n", bookID, memberID)
			}
		case <-cancelCh:
			// reservation completed (borrow happened), stop timer
			timer.Stop()
		}
	}()

	// Process borrowing asynchronously.
	// The spec says: "If the book is available, reserve it and process borrowing asynchronously."
	// We'll simulate an asynchronous borrow attempt that may happen before the 5s timer.
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		// simulate variable processing time (between 500ms - 4s)
		rand.Seed(time.Now().UnixNano())
		delay := time.Duration(500+rand.Intn(3500)) * time.Millisecond
		time.Sleep(delay)

		// Attempt to borrow
		err := l.ProcessBorrow(bookID, memberID)
		if err != nil {
			// borrow failed - could be because auto-cancel happened first
			fmt.Printf("[AsyncBorrow] member %d failed to borrow book %d: %v\n", memberID, bookID, err)
			return
		}
		// success
		fmt.Printf("[AsyncBorrow] member %d successfully borrowed book %d (after %v)\n", memberID, bookID, delay)
		// inform timer goroutine to stop
		close(cancelCh)
	}()

	return nil
}

// Shutdown gracefully shuts down worker(s).
func (l *InMemoryLibrary) Shutdown() {
	close(l.quit)
	// drain channel to avoid goroutine leaks
	l.mu.Lock()
	l.workerOn = false
	l.mu.Unlock()
	l.wg.Wait()
}
