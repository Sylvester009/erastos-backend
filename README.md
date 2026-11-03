# Erastos Backend

A Go REST API for managing books in a library. It uses Gin for HTTP routing
and pgx to connect to PostgreSQL.

## Current progress

Implemented:

- Book endpoints for listing, retrieving, adding, replacing, and deleting books.
- A `Book` model with authors, optional bibliographic fields, add mode
  (`MANUAL`, `SCAN`, or `SEARCH`), and bookshelf ID.
- Six sample books in `data/books.go`.
- PostgreSQL connection configuration and a startup connection/ping.
- SQL up/down migration files for creating and dropping the `books` table.
- Air configuration for optional live reload.

**Important:** The API handlers still read and modify the in-memory
`data.Books` slice. They do not yet read from or write to PostgreSQL. The
database connection is currently established at startup, but the books
migration is not automatically applied by the app. API changes are lost when
the process stops.

The `bookshelf_id` column is not yet a foreign key because a bookshelves table
has not been defined.

## Project structure

```text
.
├── cmd/api/main.go                         # Application entry point
├── data/books.go                           # In-memory sample book data
├── internal/
│   ├── api/
│   │   ├── book_handlers.go                # HTTP handlers
│   │   └── router.go                       # Gin router and route registration
│   ├── config/config.go                    # Reads environment configuration
│   └── database/postgres.go                # Creates and pings pgx connection pool
├── migrations/
│   ├── 000001_create_books_table.up.sql    # Creates books table
│   └── 000001_create_books_table.down.sql  # Drops books table
├── type/book.go                            # Book model
├── .air.toml                               # Air live-reload configuration
├── go.mod
└── go.sum
```

## Requirements

- Go 1.27.3 or compatible, as specified in `go.mod`.
- A reachable PostgreSQL database.
- Network access on the first run so Go can download dependencies.

## Configuration

The application loads a `.env` file from the current working directory. Create
one in the project root with:

```dotenv
DATABASE_URL=postgresql://USER:PASSWORD@HOST:5432/DATABASE?sslmode=require
PORT=8080
```

Replace the placeholders with your database connection details. Keep `.env`
and all real credentials private; do not commit them to source control.
`DATABASE_URL` is used to establish and ping the PostgreSQL connection.
`PORT` sets the HTTP server port.

## Run the API

From the project root:

```bash
go run ./cmd/api
```

The server listens at `http://localhost:8080` when `PORT=8080`.

To build the API:

```bash
go build ./cmd/api
```

### Live reload

The repository includes an `.air.toml` configuration for
[Air](https://github.com/air-verse/air). Install Air, then run from the project
root:

```bash
air
```

## Database migration

The migration files define the schema, but the application does not run them
automatically. Apply `migrations/000001_create_books_table.up.sql` to your
PostgreSQL database using your preferred migration tool or SQL client. To
reverse it, apply `migrations/000001_create_books_table.down.sql`.

The table stores authors in a PostgreSQL `TEXT[]`. `created_at` and
`updated_at` default to the current timestamp when a row is inserted; updating
`updated_at` on later edits still needs application logic or a database trigger.

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
