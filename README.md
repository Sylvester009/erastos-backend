# Erastos

Erastos is a personal library project with a Go REST API, PostgreSQL database,
and a Next.js frontend. This milestone includes authenticated book management
in the API and a responsive frontend prototype.

## Project status

### Backend API

- Book CRUD endpoints use PostgreSQL and UUID book IDs.
- Each book belongs to a user. All `/books` routes require a valid JWT and
  operate only on books owned by the authenticated user.
- User registration stores passwords as bcrypt hashes. Login returns an
  HS256-signed JWT that expires after 72 hours.
- `GET /auth/user` returns the authenticated user's profile.
- CORS allows the configured frontend origins.
- Database schema changes are managed with manual `migrate` commands.

### Frontend

The Next.js app provides a responsive library UI, sample book covers and
details, and designed add/edit, sign-up, and login pages. These pages are
currently a visual prototype: the sample library is local data, and the forms
do not submit to the API or perform authentication.

### Current scope

- Books require a positive `bookshelf_id`, but a bookshelves table and
  bookshelf foreign key have not been implemented.
- The frontend and backend are not yet integrated into an end-to-end
  authenticated application.

## Repository layout

```text
erastos-backend/
├── cmd/api/main.go
├── internal/
│   ├── config/config.go
│   ├── database/postgres.go
│   ├── handlers/
│   │   ├── book_handler.go
│   │   └── user_handler.go
│   ├── middleware/auth_middlware.go
│   ├── models/
│   │   ├── book.go
│   │   └── user.go
│   └── repository/
│       ├── book_repository.go
│       └── user_repository.go
├── migrations/
│   ├── 000001_create_books_table.*.sql
│   ├── 000002_create_user_table.*.sql
│   ├── 000003_convert_book_ids_to_uuid.*.sql
│   └── 000004_add_userid_to_book_table.*.sql
└── README.md

erastos-frontend/
├── app/
│   ├── [id]/                   # Sample book details and edit design
│   ├── add-book/               # Add-book form design
│   ├── login/                  # Login design
│   ├── signup/                 # Sign-up design
│   └── page.tsx                # Sample library page
├── next.config.ts
└── package.json
```

## Requirements

- Go version specified in `go.mod`
- PostgreSQL
- The `migrate` CLI
- Node.js and npm for the frontend

## Backend setup

From `erastos-backend`, create a private `.env` file:

```dotenv
DATABASE_URL="postgres://USER:PASSWORD@HOST:5432/DATABASE?sslmode=require"
PORT=8080
JWT_SECRET="replace-with-a-long-random-secret"
CORS_ORIGINS="http://localhost:3000,http://127.0.0.1:3000"
```

Use your database provider's connection URL and a strong, private JWT secret.
Do not commit `.env` or real credentials. The application currently expects
`.env` to exist when it starts.

`CORS_ORIGINS` is a comma-separated list of exact origins (scheme and port,
without a path). It defaults to `http://localhost:3000` and
`http://127.0.0.1:3000`. Add the deployed frontend origin to the list when
deploying. Restart the API after changing `.env`.

### Apply migrations

The API does not run migrations at startup. From the backend directory:

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

Migration `000004` adds a non-null `user_id` foreign key to `books`. It can
only apply directly when the `books` table has no existing rows. If books
already exist, assign each one to a valid user as part of a migration before
making `user_id` non-null. Do not edit migrations that have already been
applied; create a new forward migration for later schema changes.

Migration `000003` changes existing book IDs to UUIDs, so numeric IDs from an
older schema do not remain the same.

### Run the API

```bash
go run ./cmd/api
```

With `PORT=8080`, the API listens at `http://localhost:8080`. Build and test:

```bash
go build ./cmd/api
go test ./...
```

The current server binds to `localhost`; change the bind address before
deploying it where external network access is required.

## Frontend setup

From `erastos-frontend`:

```bash
npm install
npm run dev
```

The development app runs at `http://localhost:3000`. To create a production
build:

```bash
npm run build
npm start
```

The Next.js config includes a `/api/books` rewrite to the local backend. Set
`BOOKS_API_URL` to the backend base URL when using that proxy outside local
development. The current book forms are design-only and do not submit requests.

## API reference

| Method | Path | Authentication | Description |
|---|---|---|---|
| `POST` | `/auth/register` | No | Create an account |
| `POST` | `/auth/login` | No | Verify credentials and return a JWT |
| `GET` | `/auth/user` | Bearer token | Get the authenticated user's profile |
| `GET` | `/books` | Bearer token | List the authenticated user's books |
| `GET` | `/books/:id` | Bearer token | Get one of the user's books |
| `POST` | `/books` | Bearer token | Create a book for the authenticated user |
| `PUT` | `/books/:id` | Bearer token | Replace one of the user's books |
| `DELETE` | `/books/:id` | Bearer token | Delete one of the user's books |

The API derives book ownership from the verified JWT. Do not send `user_id` in
book create or update requests.

### Register

Registration requires a username, valid email, and password of at least eight
characters:

```bash
curl http://localhost:8080/auth/register \
    --include \
    --header "Content-Type: application/json" \
    --request POST \
    --data '{"username":"reader","email":"reader@example.com","password":"replace-with-a-password"}'
```

### Log in

Login returns a JSON object containing `token`. Copy that value and set it as
`TOKEN` in your shell:

```bash
curl http://localhost:8080/auth/login \
    --include \
    --header "Content-Type: application/json" \
    --request POST \
    --data '{"email":"reader@example.com","password":"replace-with-a-password"}'

TOKEN="paste-the-token-value-here"
```

Send it on protected requests as `Authorization: Bearer $TOKEN`.

### List books

```bash
curl http://localhost:8080/books \
    --include \
    --header "Authorization: Bearer $TOKEN"
```

### Create a book

The database generates the book UUID and timestamps. `add_mode` must be
`MANUAL`, `SCAN`, or `SEARCH`; `bookshelf_id` must be a positive integer.

```bash
curl http://localhost:8080/books \
    --include \
    --header "Content-Type: application/json" \
    --header "Authorization: Bearer $TOKEN" \
    --request POST \
    --data '{"title":"The Name of the Wind","authors":["Patrick Rothfuss"],"cover_url":null,"description":null,"page_count":null,"published_date":null,"publisher":null,"isbn":"9780756404741","add_mode":"MANUAL","volume_id":null,"bookshelf_id":1}'
```
