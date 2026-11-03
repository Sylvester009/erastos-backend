package main

import (
	"log"
	"github.com/jackc/pgx/v5/pgxpool"

	"erastos-backend/internal/api"
	"erastos-backend/internal/config"
	"erastos-backend/internal/database"
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

	router := api.BookRouter()
	router.SetTrustedProxies(nil)
	router.Run("localhost:" + cfg.Port)
}
