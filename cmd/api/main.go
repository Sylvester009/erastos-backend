package main

import (
	"log"

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

	pool, err := database.Connect(cfg.DatabaseURL)

	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	defer pool.Close()

	router := api.BookRouter(pool)
	router.SetTrustedProxies(nil)
	
	if err := router.Run("localhost:" + cfg.Port); err != nil {
		log.Fatal("Error running server: ", err)
	}
}
