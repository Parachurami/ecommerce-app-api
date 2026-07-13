package db

import (
	"context"
	"log"
	"time"

	"github.com/Parachurami/ecommerce-app-api/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitDB(ctx context.Context) (*pgxpool.Pool, *redis.Client, error) {
	options, redisErr := redis.ParseURL(config.REDIS_URL)
	if redisErr != nil {
		log.Print(redisErr)
		return nil, nil, redisErr
	}
	client := redis.NewClient(options)
	log.Print("Redis Connected Successfully: ", options.ClientName)
	config, err := pgxpool.ParseConfig(config.DB_CONN_STRING)
	if err != nil {
		log.Printf("Error occurred while parsing config for DB: %v", err)
		return nil, nil, err
	}
	config.MaxConns = 25                      // Maximum number of active connections
	config.MinConns = 5                       // Minimum number of idle connections
	config.MaxConnLifetime = time.Hour        // How long a connection can live before being destroyed
	config.MaxConnIdleTime = 30 * time.Minute // How long a connection can sit idle

	conn, connErr := pgxpool.NewWithConfig(ctx, config)
	if connErr != nil {
		log.Printf("Error occurred while connecting to DB: %v", err)
		return nil, nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		log.Printf("Error occurred while pinging DB: %v", err)
		return nil, nil, err
	}
	log.Print("Successfully connected to DB!!")
	return conn, client, nil
}
