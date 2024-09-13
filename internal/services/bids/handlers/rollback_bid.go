package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *handler) RollbackBid() gin.HandlerFunc {
	return func(c *gin.Context) {
		bidIdParam := c.Param("bidId")
		versionParam := c.Param("version")
		username := c.Query("username")

		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
			return
		}

		bidId, err := uuid.Parse(bidIdParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bid ID"})
			return
		}

		version, err := strconv.Atoi(versionParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version"})
			return
		}

		bid,err := h.storage.RollbackBid(bidId, version, username)
		if err != nil {
			if err.Error() == "bid not found"{
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			if err.Error() == "version not found"{
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			if err.Error() == "user not found" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}
			if err.Error() == "unauthorized: user is not responsible for this bid's organization" {
				c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to rollback bid"})
			return
		}

		c.JSON(http.StatusOK, bid)
	}
}
