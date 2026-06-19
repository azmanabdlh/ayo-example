package config

import (
	"fmt"
	"os"
)

type Config struct {
	Server   Server
	Database Database
}

type Server struct {
	Address string
}

type Database struct {
	DSN string
}

func New() (*Config, error) {
	return newConfig(), nil
}

func newConfig() *Config {

	dbDSN := os.Getenv("DATABASE_URL")

	if dbDSN == "" {
		dbDSN = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"),
		)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		Server: Server{
			Address: port,
		},
		Database: Database{
			DSN: dbDSN,
		},
	}
}
