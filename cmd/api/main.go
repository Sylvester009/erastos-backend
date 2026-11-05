package main

import (
	"log"

	"erastos-backend/internal/config"
	"erastos-backend/internal/database"
	"erastos-backend/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	var cfg *config.Config
	var err error

	cfg, err = config.LoadConfig()

	if err != nil {
		log.Fatal("Error loading configuration: ", err)
	}

	var pool *pgxpool.Pool

	pool, err = database.Connect(cfg.DatabaseURL)

	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	defer pool.Close()

	var router *gin.Engine = gin.Default()

	router.SetTrustedProxies(nil)

	router.GET("/books", handlers.GetBooksHandler(pool))
	router.GET("/books/:id", handlers.GetBookByIDHandler(pool))
	router.POST("/books", handlers.CreateBookHandler(pool))
	router.PUT("/books/:id", handlers.UpdateBookHandler(pool))
	router.DELETE("/books/:id", handlers.DeleteBookHandler(pool))

	router.POST("/auth/register", handlers.CreateUserHandler(pool))
	router.POST("/auth/login", handlers.LoginHandler(pool, cfg))

	if err := router.Run("localhost:" + cfg.Port); err != nil {
		log.Fatal("Error running server: ", err)
	}

}
