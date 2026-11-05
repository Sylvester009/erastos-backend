# Erastos

Erastos is a personal library project with a Go/Gin REST API, PostgreSQL
storage, and a Next.js frontend.

## Current progress

### Backend

- Book endpoints support listing, retrieving by UUID, creating, replacing,
  and deleting books using PostgreSQL storage.
- User registration stores a bcrypt-hashed password. Login checks the password
  and returns a signed JWT that expires after 72 hours.
- The API connects to PostgreSQL through a pgx connection pool and checks the
  connection at startup.
- Database migrations are managed manually with `migrate`; the server does not
  apply migrations automatically.

Authentication is not yet enforced on the book endpoints: they are currently
public, and book records are not scoped to a user. `bookshelf_id` is required
for a book but is not yet a foreign key because a bookshelves table has not
been implemented.

### Frontend

The sibling `../erastos-frontend` Next.js app has a library page and book
details, add/edit book screens, and signup/login screens. The signup, login,
and book form screens are currently visual designs only; they are not connected
to the API or authentication flow.

## Project structure

```text
.
├── cmd/api/main.go                              # Server setup and route registration
├── internal/
│   ├── config/config.go                         # .env configuration
│   ├── database/postgres.go                     # PostgreSQL connection pool
│   ├── handlers/book_handler.go                 # Book HTTP handlers
│   ├── handlers/user_handler.go                 # Registration and login handlers
│   ├── models/book.go                           # Book model
│   ├── models/user.go                           # User model
│   ├── repository/book_repository.go             # Book database queries
│   └── repository/user_repository.go             # User database queries
├── migrations/
│   ├── 000001_create_books_table.*.sql          # Books table
│   ├── 000002_create_user_table.*.sql           # Users table
│   └── 000003_convert_book_ids_to_uuid.*.sql    # Converts old book IDs to UUIDs
├── .air.toml                                    # Air live-reload configuration
├── go.mod
└── go.sum
```

## Requirements

- The Go version specified in `go.mod`.
- A reachable PostgreSQL database.
- The `migrate` CLI to apply or roll back database migrations.

## Configure the application

Create a `.env` file in the project root:

```dotenv
DATABASE_URL="postgres://USER:PASSWORD@HOST:5432/DATABASE?sslmode=require"
PORT=8080
JWT_SECRET="replace-with-a-long-random-secret"
```

Use your database provider's connection URL and a strong, private JWT secret.
Keep `.env` out of version control; never commit real credentials. The current
configuration loader expects a `.env` file.

## Apply database migrations

Migrations are not run by the API. From the project root, load the variables
from `.env` and apply them:

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

Migration files that have already been applied should not be edited to change
the database schema. Add a new forward migration instead. Migration `000003`
assigns UUIDs to existing book rows, so existing numeric book IDs change when
that migration is applied.

## Run the API

From the project root:

```bash
go run ./cmd/api
```

With `PORT=8080`, the server listens at `http://localhost:8080`. To build:

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
| `POST` | `/auth/register` | Register a user |
| `POST` | `/auth/login` | Log in and receive a JWT |

### List books

```bash
curl --include http://localhost:8080/books
```

### Create a book

```bash
curl http://localhost:8080/books \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"title":"The Name of the Wind","authors":["Patrick Rothfuss"],"cover_url":null,"description":null,"page_count":null,"published_date":null,"publisher":null,"isbn":"9780756404741","add_mode":"MANUAL","volume_id":null,"bookshelf_id":1}'
```

The database generates `id`, `created_at`, and `updated_at`. A create request
requires a title, at least one author, a positive `bookshelf_id`, and an
`add_mode` of `MANUAL`, `SCAN`, or `SEARCH`. `PUT` replaces the full book
record, so include all required fields.

### Register a user

Passwords must contain at least eight characters. The stored password is
bcrypt-hashed and is not included in the response.

```bash
curl http://localhost:8080/auth/register \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"username":"reader","email":"reader@example.com","password":"replace-with-a-password"}'
```

### Log in

Successful login returns a signed JWT. The current book endpoints do not yet
require this token.

```bash
curl http://localhost:8080/auth/login \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"email":"reader@example.com","password":"replace-with-a-password"}'
```
