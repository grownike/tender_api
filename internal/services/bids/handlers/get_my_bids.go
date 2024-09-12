package handlers

import (
	"avito_tenders/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) GetMyBids() gin.HandlerFunc {
	return func(c *gin.Context) {

		username := c.Query("username")

		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
			return
		}

		var bid []models.Bid

		query := h.storage.GetMyBids(username)

		if err := query.Find(&bid).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve bids"})
			return
		}

		c.JSON(http.StatusOK, bid)
	}
}
