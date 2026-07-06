package main

import (
	"os"

	"oqu/internal/app"
	"oqu/internal/configs"
	"oqu/internal/repository/postgresql"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	cfg := configs.NewPostgresqlConfig()
	db := postgresql.NewPostgresqlConn(cfg)
	defer db.Close()

	cache := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "",
		DB:       0,
		Protocol: 2,
	})
	defer cache.Close()

	app.Bastau(db, cache)
}
