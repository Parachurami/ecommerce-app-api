package db

import (
	"context"
	"log"
	"time"

	"github.com/Parachurami/ecommerce-app-api/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB(ctx context.Context) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(config.DB_CONN_STRING)
	if err != nil {
		log.Printf("Error occurred while parsing config for DB: %v", err)
		return nil, err
	}
	config.MaxConns = 25                      // Maximum number of active connections
	config.MinConns = 5                       // Minimum number of idle connections
	config.MaxConnLifetime = time.Hour        // How long a connection can live before being destroyed
	config.MaxConnIdleTime = 30 * time.Minute // How long a connection can sit idle

	conn, connErr := pgxpool.NewWithConfig(ctx, config)
	if connErr != nil {
		log.Printf("Error occurred while connecting to DB: %v", err)
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		log.Printf("Error occurred while pinging DB: %v", err)
		return nil, err
	}
	log.Print("Successfully connected to DB!!")
	return conn, nil
}
