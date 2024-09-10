package storage

import (
	"avito_tenders/internal/db"
	"avito_tenders/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type storage struct {
	Database *db.Database
}

// New создаёт новый сервис для работы с
func New(database *db.Database) *storage {
	return &storage{
		Database: database,
	}
}

// sd
func (s *storage) CreateBids(c *gin.Context) {
	var tender models.Tender
	if err := c.BindJSON(&tender); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Сохранение в базе данных
	if err := s.Database.DB.Create(&tender).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tender"})
		return
	}

	c.JSON(http.StatusOK, tender)
}
