package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

func ConnectToPostgreSQL() (*pgxpool.Pool, error) {
	// Connect to PostgreSQL using the environment variable
	connectionString := os.Getenv("POSTGRES_CONNECTION_STRING")
	if connectionString == "" {
		log.Fatal("POSTGRES_CONNECTION_STRING environment variable not set")
	}

	// Set client options
	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, err
	}

	// Connect to PostgreSQL
	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	// Ping the PostgreSQL server to check the connection
	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to PostgreSQL!")
	return pool, nil
}
