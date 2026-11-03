package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(databaseURL string) (*pgxpool.Pool, error) {
	var ctx context.Context = context.Background()

	var config *pgxpool.Config
	var err error

	config, err = pgxpool.ParseConfig(databaseURL)

	if err != nil {
		log.Printf("Error parsing database URL: %v", err)
		return nil, err
	}
	var pool *pgxpool.Pool

	pool, err = pgxpool.NewWithConfig(ctx, config)

	if err != nil {
		log.Printf("Error creating database pool: %v", err)
		return nil, err
	}

	err = pool.Ping(ctx)
	
	if err != nil {
		log.Printf("Error pinging database: %v", err)
		return nil, err
	}

	log.Println("Successfully connected to the database")

	return pool, nil
}