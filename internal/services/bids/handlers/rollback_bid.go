package handlers

import (
	"avito_tenders/internal/models"
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

		var bidVersion models.BidVersion
		if err := h.storage.GetBidVersion(bidId, version, &bidVersion, username); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Version not found"})
			return
		}

		if err := h.storage.RollbackBid(bidId, bidVersion); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to rollback bid"})
			return
		}
		bid, err := h.storage.GetBidById(bidId)
		if err != nil {
			if err.Error() == "bid not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return

			}
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version"})
			return
		}

		c.JSON(http.StatusOK, bid)
	}
}
