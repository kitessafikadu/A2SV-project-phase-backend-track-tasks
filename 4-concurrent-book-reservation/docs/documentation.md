# Concurrency Documentation — Library Management (Reservation Worker)

## Overview

This implementation adds concurrent book reservation support using:

- Goroutines to handle concurrent tasks and to simulate asynchronous work.
- A buffered channel to queue incoming reservation requests.
- A global `sync.Mutex` to protect shared state (books map and book fields) and prevent races.
- A timer-based goroutine per reservation that auto-cancels the reservation if not borrowed within 5 seconds.

## Key components

- `InMemoryLibrary` (services):
  - `reqCh` (chan ReservationRequest): queue for reservation requests.
  - `StartReservationWorker()`: launches a dispatcher goroutine that reads from `reqCh` and dispatches each request to a concurrent handler goroutine.
  - `handleReservation(bookID, memberID)`: reserves a book, starts a 5-second auto-cancel timer, and triggers an asynchronous borrow attempt.
  - `ProcessBorrow(bookID, memberID)`: finalizes borrowing. Called by async borrow routine or explicitly by controller.

## Concurrency patterns used

1. **Channel-based queue**: external callers call `ReserveBook`, which pushes a `ReservationRequest` to `reqCh`. The worker reads requests and spawns handler goroutines. This decouples producers (controllers) from consumers (workers).
2. **Dispatcher + worker goroutines**: a single dispatcher reads the queue and spawns multiple worker goroutines so multiple reservations can be processed concurrently.
3. **Mutex (sync.Mutex)**: a single mutex `mu` on `InMemoryLibrary` protects the `books` map and modifications to `Book` fields (ReservedBy, BorrowedBy, Available).
4. **Timer goroutines for auto-cancel**: each successful reservation launches a timer goroutine that waits 5 seconds; if the book isn't borrowed by then, it clears the reservation.
5. **Asynchronous borrow**: after reservation, the system simulates a borrow attempt in a separate goroutine (random delay). If the borrow completes before the 5-second timer, the timer is canceled.

## Safety notes

- All reads and writes to book state are protected by the central `mu` mutex to prevent race conditions.
- The design is intentionally simple: a single mutex simplifies correctness. For higher throughput you can:
  - Use per-book mutexes or
  - Use sharded locks or
  - Use atomic operations for some fields
- Worker dispatching uses a buffered channel to handle bursts of incoming requests.

## How the auto-cancellation works

1. Reservation succeeds and `ReservedBy` is set.
2. A timer goroutine is created and waits 5 seconds.
3. If the book is borrowed before timeout, the borrow routine closes a cancel channel to stop the timer.
4. If timer fires first, it re-acquires the mutex and clears `ReservedBy`, setting `Available = true`.

## How to run the demo

1. `go run main.go`
2. The `main` function simulates multiple members attempting to reserve the same book concurrently.
3. Observe printed log lines showing which reservations were accepted, which asynchronous borrows succeeded, and which reservations auto-cancelled.

## Further improvements

- Persist state to DB.
- Replace global mutex with per-book locks for better parallelism.
- Add retry/backoff for failed reservations.
- Convert controller to HTTP endpoints.
