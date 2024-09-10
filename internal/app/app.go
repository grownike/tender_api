package app

import (
	"avito_tenders/internal/db"
	bids_handlers "avito_tenders/internal/services/bids/handlers"
	bids_storage "avito_tenders/internal/services/bids/storage"
	tender_handlers "avito_tenders/internal/services/tender/handlers"
	tender_storage "avito_tenders/internal/services/tender/storage"
	"log"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func SetupRoutes(database *db.Database) *gin.Engine {
	r := gin.Default()

	tenderStorage := tender_storage.New(database)
	tenderHandler := tender_handlers.New(tenderStorage)

	bidsStorage := bids_storage.New(database)
	bidsHandler := bids_handlers.New(bidsStorage)

	// Маршруты для предложений (tender)

	r.POST("/api/tenders/new", tenderHandler.CreateTender())

	// Маршруты для предложений (bids)

	r.POST("/api/bids/new", bidsHandler.CreateBids())
	// r.GET("/api/bids/my", bidService.GetMyBids)

	return r
}

func Start_Server() {
	// Загружаем переменные окружения
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Инициализация маршрутов
	r := SetupRoutes(database)

	// Запуск HTTP сервера
	serverAddress := os.Getenv("SERVER_ADDRESS")
	if err := r.Run(serverAddress); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}
