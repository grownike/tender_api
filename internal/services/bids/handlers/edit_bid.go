package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *handler) EditBid() gin.HandlerFunc {
	return func(c *gin.Context) {
		BidIdParam := c.Param("bidId")
		username := c.Query("username")

		print(username)

		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
			return
		}

		BidId, err := uuid.Parse(BidIdParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bid ID"})
			return
		}


		jsonDataBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		var updates map[string]interface{}

		err = json.Unmarshal(jsonDataBytes, &updates)
		if err != nil {
			c.JSON(500, gin.H{"error Unmarshal": err.Error()})
			return
		}

		updatedBid, err := h.storage.UpdateBid(c, BidId, updates, username)
		if err != nil {
			if err.Error() == "unauthorized: you are not the creator of this bid" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}
			if err.Error() == "bid not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return

			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update bid"})
			return
		}

		c.JSON(http.StatusOK, updatedBid)
	}
}

