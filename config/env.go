package config

import "os"

var (
	PORT           string = GetStringEnv("PORT", ":3000")
	DB_CONN_STRING string = GetStringEnv("DB_CONN_STRING", "postgres://admin:adminpassword@localhost:5432/ecommerce_db?ssl_mode=disable")
	JWT_SECRET     string = GetStringEnv("JWT_SECRET", "my-secret")
)

func GetStringEnv(key, fallback string) string {
	env := os.Getenv(key)
	if env == "" {
		return fallback
	}
	return env
}
