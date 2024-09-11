package handlers

import (
	"avito_tenders/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *handler) GetTenders() gin.HandlerFunc {
	return func(c *gin.Context) {

		limitParam := c.DefaultQuery("limit", "5")
		offsetParam := c.DefaultQuery("offset", "0")
		serviceType := c.QueryArray("service_type")

		limit, err := strconv.Atoi(limitParam)
		if err != nil || limit <= 0 {
			limit = 5
		}

		offset, err := strconv.Atoi(offsetParam)
		if err != nil || offset < 0 {
			offset = 0
		}

		var tenders []models.Tender

		query := h.storage.GetTenders(limit, offset, serviceType)

		if err := query.Find(&tenders).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tenders"})
			return
		}

		c.JSON(http.StatusOK, tenders)
	}
}
