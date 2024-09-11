package handlers

import (
	"avito_tenders/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type storage interface {
	CreateBid(c *gin.Context, bid *models.Bid) error
	GetMyBids(username string) *gorm.DB
}
