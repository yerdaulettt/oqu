package main

import (
	"os"

	"oqu/internal/app"
	"oqu/internal/auth"
	"oqu/internal/configs"
	"oqu/internal/repository/postgresql"
	"oqu/internal/repository/rediscache"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

// @title Oqu REST API
func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	postgresqlCfg := configs.NewPostgresqlConfig()
	db := postgresql.NewPostgresqlConn(postgresqlCfg)
	defer db.Close()

	cache := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "",
		DB:       0,
		Protocol: 2,
	})
	defer cache.Close()

	cacheRepo := rediscache.NewCacheRepository(cache)

	jwtCfg := configs.NewJwtConfig()
	jwtService := auth.NewJwtAuth(jwtCfg.Secret, jwtCfg.AccessTtl, jwtCfg.RefreshTtl)

	app.Bastau(db, cacheRepo, jwtService)
}
