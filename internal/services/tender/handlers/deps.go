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
	GetMyTenders(username string) *gorm.DB
	UpdateTender(c *gin.Context, tenderId uuid.UUID, updates map[string]interface{}, username string) (*models.Tender, error)
	GetTenderStatus(tenderId uuid.UUID, username string) (string, error)
	EditTenderStatus(tenderId uuid.UUID, username string, newStatus string) (*models.Tender, error)
	GetTenderVersion(tenderId uuid.UUID, version int, tenderVersion *models.TenderVersion, username string) error
	RollbackTender(tenderId uuid.UUID, version models.TenderVersion) error
	GetTenderById(tenderId uuid.UUID) (*models.Tender, error)
}
