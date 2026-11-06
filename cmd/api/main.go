package main

import (
	"log"
	"time"

	"erastos-backend/internal/config"
	"erastos-backend/internal/database"
	"erastos-backend/internal/handlers"
	"erastos-backend/internal/middleware"

	"github.com/gin-contrib/cors"
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
	router.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSOrigins,
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		MaxAge:       12 * time.Hour,
	}))

	books := router.Group("/books", middleware.AuthMiddleware(cfg))
	books.GET("", handlers.GetBooksHandler(pool))
	books.GET("/:id", handlers.GetBookByIDHandler(pool))
	books.POST("", handlers.CreateBookHandler(pool))
	books.PUT("/:id", handlers.UpdateBookHandler(pool))
	books.DELETE("/:id", handlers.DeleteBookHandler(pool))

	router.POST("/auth/register", handlers.CreateUserHandler(pool))
	router.POST("/auth/login", handlers.LoginHandler(pool, cfg))
	router.GET("/auth/user", middleware.AuthMiddleware(cfg), handlers.GetUserHandler(pool))

	if err := router.Run("localhost:" + cfg.Port); err != nil {
		log.Fatal("Error running server: ", err)
	}

}
