# Erastos Backend

A Go REST API for managing books in a library. The HTTP server uses Gin.

## Project structure

```text
.
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── data/
│   └── books.go                 # In-memory sample book data
├── internal/
│   └── api/
│       ├── book_handlers.go     # Book request handlers
│       └── router.go            # Gin router and route registration
├── type/
│   └── book.go                  # Book data model
├── .air.toml                   # Air live-reload configuration
├── go.mod
└── go.sum
```

## Current progress

The API currently supports:

- Listing all books and retrieving one book by its string ID.
- Adding a book, replacing a book, and deleting a book.
- Six sample books loaded from `data/books.go`.
- A `Book` model with authors, optional bibliographic fields, an add mode
  (`MANUAL`, `SCAN`, or `SEARCH`), and a bookshelf ID.

The books are stored in memory only. Changes made through the API are lost when
the server stops. Persistent database storage and bookshelf relationship
enforcement are not implemented yet.

## Requirements

- Go 1.27.3 or compatible, as specified in `go.mod`.
- Network access on the first run so Go can download dependencies.

## Run the API

From the project root:

```bash
go run ./cmd/api
```

The server listens at `http://localhost:8080`.

To build the API:

```bash
go build ./cmd/api
```

### Live reload

The repository includes an `.air.toml` configuration for [Air](https://github.com/air-verse/air).
After installing Air, run this from the project root:

```bash
air
```

## API endpoints

| Method | Path | Description |
|---|---|---|
| `GET` | `/books` | Return all books |
| `GET` | `/books/:id` | Return a book by its string ID |
| `POST` | `/books` | Add a book from a JSON request body |
| `PUT` | `/books/:id` | Replace a book using a JSON request body |
| `DELETE` | `/books/:id` | Delete a book |

Example: list books

```bash
curl --include http://localhost:8080/books
```

Example: add a book

```bash
curl http://localhost:8080/books \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"id":"7","title":"The Name of the Wind","authors":["Patrick Rothfuss"],"coverUrl":null,"description":null,"pageCount":null,"publishedDate":null,"publisher":null,"isbn":"9780756404741","addMode":"MANUAL","volumeId":null,"bookshelfId":1}'
```

Use the same string ID in the URL and the book JSON, such as `/books/7` and
`"id":"7"`. `PUT` replaces the whole book record; include all fields in its
request body.
