package api

import "github.com/gin-gonic/gin"

func BookRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", getBookByID)
	router.POST("/books", postBooks)
	router.PUT("/books/:id", updateBook)
	router.DELETE("/books/:id", deleteBook)
	return router
}
