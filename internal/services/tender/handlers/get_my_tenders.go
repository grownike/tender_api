package handlers

import (
	"avito_tenders/internal/models"
	"net/http"
	
	"github.com/gin-gonic/gin"
)

func (h *handler) GetMyTenders() gin.HandlerFunc {
	return func(c *gin.Context) {
		
		username := c.Query("username")

		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
			return
		}

		var tenders []models.Tender

		query := h.storage.GetMyTenders(username)

		if err := query.Find(&tenders).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tenders"})
			return
		}

		c.JSON(http.StatusOK, tenders)
	}
}
