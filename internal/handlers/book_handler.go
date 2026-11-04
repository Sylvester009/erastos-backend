package handlers

import (
	"log"
	"net/http"

	"erastos-backend/internal/models"
	"erastos-backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateBookRequest struct {
	Title         string   `json:"title" binding:"required"`
	Authors       []string `json:"authors" binding:"required,min=1"`
	CoverURL      *string  `json:"cover_url"`
	Description   *string  `json:"description"`
	PageCount     *int     `json:"page_count"`
	PublishedDate *string  `json:"published_date"`
	Publisher     *string  `json:"publisher"`
	ISBN          *string  `json:"isbn"`
	AddMode       string   `json:"add_mode" binding:"required,oneof=MANUAL SCAN SEARCH"`
	VolumeID      *string  `json:"volume_id"`
	BookshelfID   int      `json:"bookshelf_id" binding:"required,gt=0"`
}

type UpdateBookRequest struct {
	Title         string   `json:"title" binding:"required"`
	Authors       []string `json:"authors" binding:"required,min=1"`
	CoverURL      *string  `json:"cover_url"`
	Description   *string  `json:"description"`
	PageCount     *int     `json:"page_count"`
	PublishedDate *string  `json:"published_date"`
	Publisher     *string  `json:"publisher"`
	ISBN          *string  `json:"isbn"`
	AddMode       string   `json:"add_mode" binding:"required,oneof=MANUAL SCAN SEARCH"`
	VolumeID      *string  `json:"volume_id"`
	BookshelfID   int      `json:"bookshelf_id" binding:"required,gt=0"`
}

func CreateBookHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request CreateBookRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid book data"})
			return
		}

		createdBook, err := repository.CreateBook(pool, models.Book{
			Title:         request.Title,
			Authors:       request.Authors,
			CoverURL:      request.CoverURL,
			Description:   request.Description,
			PageCount:     request.PageCount,
			PublishedDate: request.PublishedDate,
			Publisher:     request.Publisher,
			ISBN:          request.ISBN,
			AddMode:       request.AddMode,
			VolumeID:      request.VolumeID,
			BookshelfID:   request.BookshelfID,
		})
		if err != nil {
			log.Printf("failed to create book: %v", err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to create book"})
			return
		}

		c.IndentedJSON(http.StatusCreated, createdBook)
	}
}

func GetBooksHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		books, err := repository.GetBooks(pool)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to retrieve books"})
			return
		}
		c.IndentedJSON(http.StatusOK, books)
	}
}

func GetBookByIDHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := parseBookID(c)
		if !ok {
			return
		}

		book, err := repository.GetBookByID(pool, id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to retrieve book"})
			return
		}
		if book == nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
			return
		}
		c.IndentedJSON(http.StatusOK, book)
	}
}

func UpdateBookHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := parseBookID(c)
		if !ok {
			return
		}

		var request UpdateBookRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid book data"})
			return
		}
		updatedBook := models.Book{
			ID:            id,
			Title:         request.Title,
			Authors:       request.Authors,
			CoverURL:      request.CoverURL,
			Description:   request.Description,
			PageCount:     request.PageCount,
			PublishedDate: request.PublishedDate,
			Publisher:     request.Publisher,
			ISBN:          request.ISBN,
			AddMode:       request.AddMode,
			VolumeID:      request.VolumeID,
			BookshelfID:   request.BookshelfID,
		}

		if updatedBook.Title == "" && updatedBook.Authors == nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "At least one of the fields must be provided for update"})
			return
		}

		book, err := repository.UpdateBook(pool, updatedBook)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to update book"})
			return
		}

		if book == nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
			return
		}

		c.IndentedJSON(http.StatusOK, book)
	}
}

func DeleteBookHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := parseBookID(c)
		if !ok {
			return
		}

		deleted, err := repository.DeleteBook(pool, id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to delete book"})
			return
		}
		if !deleted {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
			return
		}

		c.IndentedJSON(http.StatusNoContent, gin.H{"message": "book deleted successfully"})
	}
}

func parseBookID(c *gin.Context) (string, bool) {
	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "book ID is required"})
		return "", false
	}
	return id, true
}
