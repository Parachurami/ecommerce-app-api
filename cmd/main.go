package main

import (
	"context"
	"log"

	"github.com/Parachurami/ecommerce-app-api/cmd/api"
	"github.com/Parachurami/ecommerce-app-api/config"
	"github.com/Parachurami/ecommerce-app-api/internal/db"
)

func main() {
	config := api.NewConfig(
		config.DB_CONN_STRING,
		config.PORT,
	)
	ctx := context.Background()
	conn, redisClient, dbErr := db.InitDB(ctx)
	if dbErr != nil {
		log.Print(dbErr)
		return
	}
	defer conn.Close()
	app := api.NewApp(conn, redisClient, config)

	db.RunMigrations(conn)

	if err := app.Run(app.Mount()); err != nil {
		log.Fatal(err)
		return
	}
}
