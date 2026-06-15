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
		log.Printf("Unable to parse DATABASE_URL: %v\n", err)
		return nil, err
	}

	var pool *pgxpool.Pool
	pool, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Printf("Unable to create connection pool: %v\n", err)
		pool.Close()
		return nil, err
	}

	log.Printf("Successfully connected to database\n")

	return pool, nil
}
