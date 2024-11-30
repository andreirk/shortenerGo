package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Port string
	Db   DbConfig
	Auth AuthConfig
}

type DbConfig struct {
	Dsn string
}
type AuthConfig struct {
	Secret string
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file, using default config")
	}
	return &Config{
		Db: DbConfig{
			Dsn: os.Getenv("DSN"),
		},
		Auth: AuthConfig{
			Secret: os.Getenv("SECRET"),
		},
		Port: os.Getenv("PORT"),
	}
}