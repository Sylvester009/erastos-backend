# Erastos Backend

A Go REST API for managing library books. It uses Gin for HTTP routing and
PostgreSQL for persistent book storage through pgx.

## Current progress

- Book list, detail, create, update, and delete endpoints are implemented.
- API handlers use the PostgreSQL repository; they do not read or write the
  in-memory sample slice in `data/books.go`.
- The book table uses UUID IDs and stores book metadata,
  bookshelf ID, and creation/update timestamps.
- Database configuration is loaded from `.env`, and the application verifies
  the database connection at startup.
- Migrations are run manually; the API does not apply them automatically.

The `bookshelf_id` column is required but is not yet a foreign key because a
bookshelves table has not been added. User ownership/authentication is also not
implemented in the current model or API.

## Project structure

```text
.
├── cmd/api/main.go                         # Application entry point
├── data/books.go                           # Legacy in-memory sample data
├── internal/
│   ├── api/router.go                       # Gin router and route registration
│   ├── config/config.go                    # Environment configuration
│   ├── database/postgres.go                # PostgreSQL connection pool
│   ├── handlers/book_handler.go            # Book HTTP handlers
│   ├── models/book.go                      # Book model
│   └── repository/book_repository.go       # PostgreSQL queries
├── migrations/
│   ├── 000001_create_books_table.*.sql     # Creates/drops books table
│   ├── 000002_create_user_table.*.sql      # Creates/drops users table
│   └── 000003_convert_book_ids_to_uuid.*.sql # Converts existing IDs to UUIDs
├── .air.toml                               # Air live-reload configuration
├── go.mod
└── go.sum
```

## Requirements

- Go version specified in `go.mod`.
- A reachable PostgreSQL database.
- The `migrate` CLI if you want to apply migrations with the commands below.

## Configure the database

Create a `.env` file in the project root:

```dotenv
DATABASE_URL="postgresql://USER:PASSWORD@HOST:5432/DATABASE?sslmode=require"
PORT=8080
```

Replace the placeholders with your database credentials. Keep `.env` private
and never commit real credentials. `DATABASE_URL` is used by the API to connect
to PostgreSQL; `PORT` sets the HTTP server port.

## Apply the migration

The application does not run migrations at startup. From the project root,
load the variables from `.env` into your shell, then run:

```bash
set -a
source .env
set +a
migrate -path migrations -database "$DATABASE_URL" up
```

To roll back the latest migration:

```bash
migrate -path migrations -database "$DATABASE_URL" down 1
```

If migration `000001` was already applied before its SQL file changed, the
migration tool will not reapply the edited file. Use a new forward migration
or update the database schema explicitly.

Migration `000003` preserves existing book rows while assigning new UUIDs.
Existing numeric book IDs will change when it is applied.

## Run the API

From the project root:

```bash
go run ./cmd/api
```

With `PORT=8080`, the API listens at `http://localhost:8080`. To build:

```bash
go build ./cmd/api
```

For live reload, install [Air](https://github.com/air-verse/air) and run
`air` from the project root.

## API endpoints

| Method | Path | Description |
|---|---|---|
| `GET` | `/books` | List books |
| `GET` | `/books/:id` | Retrieve a book by UUID |
| `POST` | `/books` | Create a book |
| `PUT` | `/books/:id` | Replace a book |
| `DELETE` | `/books/:id` | Delete a book |

Example: list books

```bash
curl --include http://localhost:8080/books
```

Example: create a book

```bash
curl http://localhost:8080/books \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"title":"The Name of the Wind","authors":["Patrick Rothfuss"],"cover_url":null,"description":null,"page_count":null,"published_date":null,"publisher":null,"isbn":"9780756404741","add_mode":"MANUAL","volume_id":null,"bookshelf_id":1}'
```

The database generates `id`, `created_at`, and `updated_at`. The request must
include a positive `bookshelf_id`; `add_mode` must be `MANUAL`, `SCAN`, or
`SEARCH`. `PUT` replaces the full book record, so include all required fields.
