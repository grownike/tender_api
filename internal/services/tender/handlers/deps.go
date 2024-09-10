package handlers

import (
	"avito_tenders/internal/models"

	"github.com/gin-gonic/gin"
)

type storage interface{
	CreateTender(c *gin.Context, tender *models.Tender) error
}