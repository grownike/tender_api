package handlers

import (
	"avito_tenders/internal/models"

	"github.com/gin-gonic/gin"
)

type storage interface {
	CreateBid(c *gin.Context, bid *models.Bid) error
}
