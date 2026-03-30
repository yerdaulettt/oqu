package main

import (
	"oqu/internal/app"
	"oqu/internal/dbconn"
	"oqu/pkg/config"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	cfg := config.NewPostgresqlConfig()
	db := dbconn.GetDBConn(cfg)
	defer db.Close()

	app.Bastau(db)
}
