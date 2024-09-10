package storage

import (
	"avito_tenders/internal/db"
	"avito_tenders/internal/models"
	"github.com/gin-gonic/gin"
)

type storage struct {
	Database *db.Database
}

func New(database *db.Database) *storage {
	return &storage{
		Database: database,
	}
}

// CreateTender метод для создания тендера
func (s *storage) CreateTender(c *gin.Context, tender *models.Tender) error {
	if err := s.Database.DB.Create(tender).Error; err != nil {
		return  err
	}
	return nil
}
