package configs

import "os"

type PostgresqlConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	Sslmode  string
}

func NewPostgresqlConfig() *PostgresqlConfig {
	return &PostgresqlConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		Sslmode:  os.Getenv("DB_SSLMODE"),
	}
}
