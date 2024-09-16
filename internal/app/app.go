package app

import (
	"avito_tenders/internal/db"
	bids_handlers "avito_tenders/internal/services/bids/handlers"
	bids_storage "avito_tenders/internal/services/bids/storage"
	tender_handlers "avito_tenders/internal/services/tender/handlers"
	tender_storage "avito_tenders/internal/services/tender/storage"
	etc_handlers "avito_tenders/internal/services/etc/handlers"
	etc_storage "avito_tenders/internal/services/etc/storage"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func SetupRoutes(database *db.Database) *gin.Engine {
	r := gin.Default()
	r.Use(CORSMiddleware())

	r.GET("/api/ping", pingHandler)

	tenderStorage := tender_storage.New(database)
	tenderHandler := tender_handlers.New(tenderStorage)

	etcStorage := etc_storage.New(database)
	etcHandler := etc_handlers.New(etcStorage)

	bidsStorage := bids_storage.New(database)
	bidsHandler := bids_handlers.New(bidsStorage)


	//Прочее для тестов POST

	
	r.POST("/api/employees/new", etcHandler.CreateEmployee())
	r.POST("/api/organizations/new", etcHandler.CreateCompany())
	r.POST("/api/newAssign/:orgId/:user", etcHandler.AssignResponsible())

	//Тендеры. Tenders GET - POST - PUT - PATCH

	r.GET("/api/tenders", tenderHandler.GetTenders())
	r.GET("/api/tenders/my", tenderHandler.GetMyTenders())
	r.GET("/api/tenders/:tenderId/status", tenderHandler.GetTenderStatus())

	r.POST("/api/tenders/new", tenderHandler.CreateTender())

	r.PUT("/api/tenders/:tenderId/status", tenderHandler.EditTenderStatus())
	r.PUT("/api/tenders/:tenderId/rollback/:version", tenderHandler.RollbackTender())

	r.PATCH("/api/tenders/:tenderId/edit", tenderHandler.EditTender())

	//Предложение. Bids GET - POST - PUT - PATCH

	r.GET("/api/bids/my", bidsHandler.GetMyBids())

	r.GET("/api/bids/tender/:tenderId/list", bidsHandler.GetBidsByTender())
	r.GET("/api/bids/:bidId/status", bidsHandler.GetBidStatus())

	r.GET("/api/bids/reviews/:tenderId", bidsHandler.GetReviewsByTender())

	r.POST("/api/bids/new", bidsHandler.CreateBids())

	r.PUT("/api/bids/:bidId/rollback/:version", bidsHandler.RollbackBid())
	//error!
	r.PUT("/api/bids/:bidId/status", bidsHandler.EditBidStatus())

	r.PATCH("/api/bids/:bidId/edit", bidsHandler.EditBid())

	

	return r
}

func Start_Server() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	r := SetupRoutes(database)
	serverAddress := os.Getenv("SERVER_ADDRESS")
	if err := r.Run(serverAddress); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	}
}

func pingHandler(c *gin.Context) {
	c.String(200, "ok")
}
