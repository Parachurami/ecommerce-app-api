package config

import "os"

var (
	PORT           string = GetStringEnv("PORT", ":3000")
	DB_CONN_STRING string = GetStringEnv("DB_CONN_STRING", "postgres://admin:adminpassword@localhost:5432/ecommerce_db?sslmode=disable")
	JWT_SECRET     string = GetStringEnv("JWT_SECRET", "my-secret")
	REDIS_URL      string = GetStringEnv("REDIS_URL", "rediss://default:gQAAAAAAAXScAAIgcDIxZTY5YTkxZjQ2YjI0NWQ4YjQ4ZWJmMzU1NDA5NWUzMQ@pumped-yak-95388.upstash.io:6379")
)

func GetStringEnv(key, fallback string) string {
	env := os.Getenv(key)
	if env == "" {
		return fallback
	}
	return env
}
