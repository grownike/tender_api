package handlers

import (
	"avito_tenders/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *handler) GetMyTenders() gin.HandlerFunc {
	return func(c *gin.Context) {

		limitParam := c.DefaultQuery("limit", "5")
		offsetParam := c.DefaultQuery("offset", "0")
		username := c.Query("username")
		
		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
			return
		}

		limit, err := strconv.Atoi(limitParam)
		if err != nil || limit <= 0 {
			limit = 5
		}

		offset, err := strconv.Atoi(offsetParam)
		if err != nil || offset < 0 {
			offset = 0
		}

		var tenders []models.Tender

		tenders, err = h.storage.GetMyTenders(username, limit, offset)
		if err != nil {
			if err.Error() == "user not found" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, tenders)
	}
}
