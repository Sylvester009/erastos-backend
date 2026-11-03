# Erastos Backend

A Go backend for a book library, built with the Gin web framework.

## Current progress

The project currently provides an in-memory REST API for books:

- Lists all books and retrieves a book by ID.
- Adds, replaces, and deletes books.
- Starts with six sample books.
- Uses the `Book` model in `type/book.go`.
- Keeps the sample book data in `data/books.go`.
- Defines the HTTP handlers in `book_handlers.go` and registers routes in `router.go`.

The book model includes title, authors, optional bibliographic details, the way
the book was added (`MANUAL`, `SCAN`, or `SEARCH`), and a bookshelf ID.

Data is stored only in memory. Books added, updated, or deleted through the API
are lost when the server stops. A database, persistent storage, validation,
authentication, and bookshelf relationship enforcement are not implemented yet.

## Requirements

- Go 1.27.3 or compatible version, as specified in `go.mod`.
- Network access on the first run so Go can download dependencies.

## Run the server

From the project directory:

```bash
go run .
```

The server listens at `http://localhost:8080`.

## API endpoints

| Method | Path | Description |
|---|---|---|
| `GET` | `/books` | Return all books |
| `GET` | `/books/:id` | Return one book by its string ID |
| `POST` | `/books` | Add a book from a JSON request body |
| `PUT` | `/books/:id` | Replace a book by ID using a JSON request body |
| `DELETE` | `/books/:id` | Delete a book by ID |

Example request to list books:

```bash
curl --include http://localhost:8080/books
```

Example request to add a book:

```bash
curl http://localhost:8080/books \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"id":"7","title":"The Name of the Wind","authors":["Patrick Rothfuss"],"coverUrl":null,"description":null,"pageCount":null,"publishedDate":null,"publisher":null,"isbn":"9780756404741","addMode":"MANUAL","volumeId":null,"bookshelfId":1}'
```

Use string IDs in the URL and JSON body, for example `/books/7` and `"id":"7"`.
The `PUT` endpoint replaces the full book record rather than updating only the
fields included in the request.
