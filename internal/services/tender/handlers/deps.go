//go:generate mockgen -source=deps.go -destination=mock_test.go -package=handlers

package handlers

import (
	"avito_tenders/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type storage interface {
	CreateTender(c *gin.Context, tender *models.Tender) error
	GetTenders(limit, offset int, serviceType []string) *gorm.DB
	GetMyTenders(username string, limit int, offset int) ([]models.Tender, error)
	UpdateTender(c *gin.Context, tenderId uuid.UUID, updates map[string]interface{}, username string) (*models.Tender, error)
	GetTenderStatus(tenderId uuid.UUID, username string) (string, error)
	EditTenderStatus(tenderId uuid.UUID, username string, newStatus string) (*models.Tender, error)
	RollbackTender(tenderId uuid.UUID, version int, username string) (*models.Tender, error)
}
