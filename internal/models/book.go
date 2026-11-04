package models

import "time"

type Book struct {
	ID            string    `json:"id" db:"id"`
	Title         string    `json:"title" db:"title"`
	Authors       []string  `json:"authors" db:"authors"`
	CoverURL      *string   `json:"cover_url" db:"cover_url"`
	Description   *string   `json:"description" db:"description"`
	PageCount     *int      `json:"page_count" db:"page_count"`
	PublishedDate *string   `json:"published_date" db:"published_date"`
	Publisher     *string   `json:"publisher" db:"publisher"`
	ISBN          *string   `json:"isbn" db:"isbn"`
	AddMode       string    `json:"add_mode" db:"add_mode"`
	VolumeID      *string   `json:"volume_id" db:"volume_id"`
	BookshelfID   int       `json:"bookshelf_id" db:"bookshelf_id"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}
