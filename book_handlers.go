package main

import (
	"net/http"

	"erastos-backend/data"
	types "erastos-backend/type"
	"github.com/gin-gonic/gin"
)

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, data.Books)
}

func getBookByID(c *gin.Context) {
	id := c.Param("id")

	for _, b := range data.Books {
		if string(b.ID) == id {
			c.IndentedJSON(http.StatusOK, b)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
}

func postBooks(c *gin.Context) {
	var newBook types.Book

	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	data.Books = append(data.Books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func updateBook(c *gin.Context) {
	id := c.Param("id")
	var updatedBook types.Book

	if err := c.BindJSON(&updatedBook); err != nil {
		return
	}

	for i, b := range data.Books {
		if string(b.ID) == id {
			data.Books[i] = updatedBook
			c.IndentedJSON(http.StatusOK, updatedBook)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
}

func deleteBook(c *gin.Context) {
	id := c.Param("id")

	for i, b := range data.Books {
		if string(b.ID) == id {
			data.Books = append(data.Books[:i], data.Books[i+1:]...)
			c.IndentedJSON(http.StatusOK, gin.H{"message": "book deleted"})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
}
