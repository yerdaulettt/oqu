package main

import (
	"oqu/internal/app"
	"oqu/internal/configs"
	"oqu/internal/repository/postgresql"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	cfg := configs.NewPostgresqlConfig()
	db := postgresql.NewPostgresqlConn(cfg)
	defer db.Close()

	app.Bastau(db)
}
