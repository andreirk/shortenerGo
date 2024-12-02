package main

import (
	"github.com/joho/godotenv"
	"go/adv-demo/internal/link"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database while auto migration", err)
	}
	db.AutoMigrate(&link.Link{})
}
