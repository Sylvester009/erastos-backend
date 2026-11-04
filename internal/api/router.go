package api

import (
	"erastos-backend/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func BookRouter(pool *pgxpool.Pool) *gin.Engine {
	router := gin.Default()
	router.GET("/books", handlers.GetBooksHandler(pool))
	router.GET("/books/:id", handlers.GetBookByIDHandler(pool))
	router.POST("/books", handlers.CreateBookHandler(pool))
	router.PUT("/books/:id", handlers.UpdateBookHandler(pool))
	router.DELETE("/books/:id", handlers.DeleteBookHandler(pool))
	return router
}
