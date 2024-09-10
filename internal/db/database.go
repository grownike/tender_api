package db

import (
	"avito_tenders/internal/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *Database {
	return &Database{DB: db}
}

func InitDB() (*Database, error) {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dsn := os.Getenv("POSTGRES_CONN")
	db, errDB := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if errDB != nil {
		return nil, errDB
	}

	err = db.AutoMigrate(&models.Employee{}, &models.Organization{}, &models.Tender{})
	if err != nil {
		return nil, err
	}

	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		log.Fatalf("Error creating extension uuid-ossp: %v", err)
	}

	log.Println("Database connected and migrations applied successfully")
	return &Database{DB: db}, nil

}
