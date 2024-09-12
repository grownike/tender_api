package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *handler) GetBidsByTender() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenderIdParam := c.Param("tenderId")
		tenderId, err := uuid.Parse(tenderIdParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"reason": "Invalid tenderId"})
			return
		}

		username := c.Query("username")
		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"reason": "Username is required"})
			return
		}

		limitParam := c.DefaultQuery("limit", "5")
		offsetParam := c.DefaultQuery("offset", "0")

		limit, err := strconv.Atoi(limitParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"reason": "Invalid limit value"})
			return
		}

		offset, err := strconv.Atoi(offsetParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"reason": "Invalid offset value"})
			return
		}

		bids, err := h.storage.GetBidsByTender(c, tenderId, username, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"reason": "Failed to fetch bids"})
			return
		}
		if len(bids) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"reason": "No bids found"})
			return
		}

		c.JSON(http.StatusOK, bids)
	}
}
