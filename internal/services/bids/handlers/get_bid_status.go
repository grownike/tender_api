package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *handler) GetBidStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		bidIdParam := c.Param("bidId")
		bidId, err := uuid.Parse(bidIdParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"reason": "Invalid bidId"})
			return
		}

		username := c.Query("username")
		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"reason": "Username is required"})
			return
		}

		status, err := h.storage.GetBidStatus(c, bidId, username)
		if err != nil {
			if err.Error() == "bid not found" {
				c.JSON(http.StatusNotFound, gin.H{"reason": "Bid not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"reason": "Failed to fetch bid status"})
			}
			return
		}

		c.JSON(http.StatusOK, status)
	}
}
