package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *handler) EditBid() gin.HandlerFunc {
	return func(c *gin.Context) {
		BidIdParam := c.Param("bidId")
		username := c.Query("username")

		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
			return
		}

		BidId, err := uuid.Parse(BidIdParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bid ID"})
			return
		}

		var updates map[string]interface{}

		if err := c.ShouldBindJSON(&updates); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		updatedBid, err := h.storage.UpdateBid(c, BidId, updates, username)
		if err != nil {
			if err.Error() == "unauthorized: user is not responsible for bid's organization" {
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
			if err.Error() == "bid is not in Created status" {
				c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update bid"})
			return
		}

		c.JSON(http.StatusOK, updatedBid)
	}
}
