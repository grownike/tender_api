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

	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		log.Fatalf("Error creating extension uuid-ossp: %v", err)
	}

	err = createEnumType(db)
	if err != nil {
		log.Fatalf("Failed to create enum type: %v", err)
	}

	err = db.AutoMigrate(&models.Employee{}, &models.Organization{}, &models.Tender{}, &models.TenderVersion{})
	if err != nil {
		return nil, err
	}

	log.Println("Database connected and migrations applied successfully")
	return &Database{DB: db}, nil

}

func createEnumType(db *gorm.DB) error {
	query := `
    DO $$ 
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'organization_type') THEN
            CREATE TYPE organization_type AS ENUM (
                'IE',
                'LLC',
                'JSC'
            );
        END IF;
    END $$;
    `
	return db.Exec(query).Error
}
