package handlers

import (
	"avito_tenders/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	
)

type storage interface {
	CreateBid(c *gin.Context, bid *models.Bid) error
	GetMyBids(username string, limit int, offset int) ([]models.Bid, error)
	GetBidsByTender(c *gin.Context, tenderId uuid.UUID, username string, limit int, offset int) ([]models.Bid, error)
	GetBidStatus(c *gin.Context, bidId uuid.UUID, username string) (string, error)
	UpdateBid(c *gin.Context, bidId uuid.UUID, updates map[string]interface{}, username string) (*models.Bid, error)
	RollbackBid(bidId uuid.UUID, version int, username string) (*models.Bid, error)
	EditBidStatus(bidId uuid.UUID, username string, newStatus string) (*models.Bid, error)
}
