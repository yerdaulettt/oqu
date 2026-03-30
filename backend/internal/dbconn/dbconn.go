package dbconn

import (
	"database/sql"
	"fmt"

	"oqu/pkg/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func autoMigrate(cfg *config.PostgresqlConfig) {
	sourceURL := "file://migrations"
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.Sslmode)

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		panic(err)
	}

	err = m.Up()

	if err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
}

func GetDBConn(cfg *config.PostgresqlConfig) *sql.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.Sslmode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	autoMigrate(cfg)

	return db
}
