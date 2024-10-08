package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const STATUS_CREATED = "Created"
const STATUS_PUBLISHED = "Published"
const STATUS_CANCELED = "Canceled"

func (h *handler) EditBidStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		bidIdParam := c.Param("bidId")
		username := c.Query("username")
		newStatus := c.Query("status")

		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
			return
		}

		if newStatus != STATUS_CREATED && newStatus != STATUS_PUBLISHED && newStatus != STATUS_CANCELED {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
			return
		}

		bidId, err := uuid.Parse(bidIdParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bid ID"})
			return
		}

		updatedBid, err := h.storage.EditBidStatus(bidId, username, newStatus)
		if err != nil {
			if err.Error() == "unauthorized: user is not responsible for this bid's organization" {
				c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
				return
			}
			if err.Error() == "bid not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			if err.Error() == "user not found" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update bid status"})
			return
		}

		c.JSON(http.StatusOK, updatedBid)
	}
}
