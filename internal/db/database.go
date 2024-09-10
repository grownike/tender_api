package db

import (
	"avito_tenders/internal/models"
	"log"
	"os"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
// 	"avito_tenders/services/tender"
// 	"avito_tenders/services/bids"
)

type Database struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *Database {
	return &Database{DB: db}
}

func InitDB() (*Database, error) {
	// Загрузка переменных окружения
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Получение строки соединения из переменных окружения
	dsn := os.Getenv("POSTGRES_CONN")
	db, errDB := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if errDB != nil {
		return nil, errDB
	}

	// Выполнение миграций для моделей
	err = db.AutoMigrate(&models.Employee{}, &models.Organization{}, &models.Tender{})
	if err != nil {
		return nil, err
	}

	log.Println("Database connected and migrations applied successfully")
	return &Database{DB: db}, nil
}

