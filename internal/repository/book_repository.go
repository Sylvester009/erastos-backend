package repository

import (
	"context"
	"errors"
	"time"

	"erastos-backend/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const bookColumns = `
	id, title, authors, cover_url, description, page_count, published_date,
	publisher, isbn, add_mode, volume_id, bookshelf_id, created_at,
	updated_at
`

type rowScanner interface {
	Scan(dest ...any) error
}

func scanBook(row rowScanner) (models.Book, error) {
	var book models.Book
	err := row.Scan(
		&book.ID,
		&book.Title,
		&book.Authors,
		&book.CoverURL,
		&book.Description,
		&book.PageCount,
		&book.PublishedDate,
		&book.Publisher,
		&book.ISBN,
		&book.AddMode,
		&book.VolumeID,
		&book.BookshelfID,
		&book.CreatedAt,
		&book.UpdatedAt,
	)
	return book, err
}

func CreateBook(pool *pgxpool.Pool, book models.Book) (*models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const query = `
		INSERT INTO books (
			title, authors, cover_url, description, page_count, published_date,
			publisher, isbn, add_mode, volume_id, bookshelf_id, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING ` + bookColumns

	createdBook, err := scanBook(pool.QueryRow(ctx, query,
		book.Title,
		book.Authors,
		book.CoverURL,
		book.Description,
		book.PageCount,
		book.PublishedDate,
		book.Publisher,
		book.ISBN,
		book.AddMode,
		book.VolumeID,
		book.BookshelfID,
		time.Now(),
		time.Now(),
	))
	if err != nil {
		return nil, err
	}

	return &createdBook, nil
}

func GetBooks(pool *pgxpool.Pool) ([]models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := pool.Query(ctx, `SELECT `+bookColumns+` FROM books`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := make([]models.Book, 0)
	for rows.Next() {
		book, err := scanBook(rows)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func GetBookByID(pool *pgxpool.Pool, id string) (*models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT ` + bookColumns + ` FROM books WHERE id = $1`
	book, err := scanBook(pool.QueryRow(ctx, query, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &book, nil
}

func UpdateBook(pool *pgxpool.Pool, book models.Book) (*models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const query = `
		UPDATE books
		SET title = $1, authors = $2, cover_url = $3, description = $4,
			page_count = $5, published_date = $6, publisher = $7, isbn = $8,
			add_mode = $9, volume_id = $10, bookshelf_id = $11,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $12
		RETURNING ` + bookColumns

	updatedBook, err := scanBook(pool.QueryRow(ctx, query,
		book.Title,
		book.Authors,
		book.CoverURL,
		book.Description,
		book.PageCount,
		book.PublishedDate,
		book.Publisher,
		book.ISBN,
		book.AddMode,
		book.VolumeID,
		book.BookshelfID,
		book.ID,
	))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &updatedBook, nil
}

func DeleteBook(pool *pgxpool.Pool, id string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	commandTag, err := pool.Exec(ctx, `DELETE FROM books WHERE id = $1`, id)
	if err != nil {
		return false, err
	}
	return commandTag.RowsAffected() > 0, nil
}
