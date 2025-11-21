# Task Manager API Documentation

Base URL: `http://localhost:8080`

## Endpoints

### GET /tasks

- Description: Return list of all tasks.
- Response: 200 OK

```json
{
  "tasks": [
    {
      "id": 1,
      "title": "Example",
      "description": "...",
      "due_date": "2025-11-20T12:00:00Z",
      "status": "todo"
    }
  ]
}
```

### GET /tasks/:id

- Description: Return a single task by id.
- Response: 200 OK (task) or 404 Not Found

### POST /tasks

- Description: Create a new task.
- Request body (JSON):

```json
{
  "title": "Finish report",
  "description": "Complete the Q4 report",
  "due_date": "2025-11-30T15:00:00Z",
  "status": "in-progress"
}
```

- Response: 201 Created (returns created task object)
- Validation: `title` and `status` are required. `due_date` must be RFC3339 (ISO8601) if provided.

### PUT /tasks/:id

- Description: Update an existing task. Provide full payload similar to POST.
- Request body (JSON): same schema as POST
- Response: 200 OK (updated task) or 404 Not Found

### DELETE /tasks/:id

- Description: Delete a task.
- Response: 204 No Content or 404 Not Found

## Testing with Postman / curl

- Create:

```
curl -X POST http://localhost:8080/tasks \
 -H "Content-Type: application/json" \
 -d '{"title":"Do homework","description":"Math","due_date":"2025-11-30T10:00:00Z","status":"todo"}'
```

- List:

```
curl http://localhost:8080/tasks
```

- Get:

```
curl http://localhost:8080/tasks/1
```

- Update:

```
curl -X PUT http://localhost:8080/tasks/1 -H "Content-Type: application/json" -d '{"title":"Updated","description":"desc","due_date":"2025-12-01T10:00:00Z","status":"done"}'
```

- Delete:

```
curl -X DELETE http://localhost:8080/tasks/1
```

````

---

## go.mod

```text
module task_manager

go 1.20

require github.com/gin-gonic/gin v1.9.0
````

---

# Notes

- This project uses an in-memory map protected by a mutex for concurrent access.
- Date/time fields expect RFC3339 (time.RFC3339) format. Empty or omitted due_date will be zero time.
- Responses use appropriate HTTP status codes: 201 for Created, 200 for OK, 204 for No Content, 400 for Bad Request, 404 for Not Found.

# How to run

1. Ensure Go (>=1.18) is installed.
2. From the project root run:

```bash
go mod tidy
go run main.go
```

Server will start at [http://localhost:8080](http://localhost:8080)

Open `docs/api_documentation.md` for Postman examples and payloads.
