package main

import (
	"github.com/joho/godotenv"
	"go/adv-demo/internal/link"
	"go/adv-demo/internal/stat"
	"go/adv-demo/internal/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Error connecting to database while auto migration", err)
	}
	err = db.AutoMigrate(&link.Link{}, &user.User{}, &stat.Stat{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database migrated")
}
