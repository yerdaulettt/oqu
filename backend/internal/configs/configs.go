package configs

import (
	"os"
	"strconv"
	"time"
)

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

type jwtConfig struct {
	Secret     []byte
	AccessTtl  time.Duration
	RefreshTtl time.Duration
}

func NewJwtConfig() *jwtConfig {
	secret := []byte(getEnv("JWT_SECRET", "secret123"))
	aTtl, err := strconv.Atoi(getEnv("JWT_ACCESS_TTL", "15"))
	if err != nil {
		aTtl = 15
	}

	rTtl, err := strconv.Atoi(getEnv("JWT_REFRESH_TTL", "24"))
	if err != nil {
		rTtl = 24
	}

	return &jwtConfig{
		Secret:     secret,
		AccessTtl:  time.Duration(aTtl) * time.Minute,
		RefreshTtl: time.Duration(rTtl) * time.Hour,
	}
}
