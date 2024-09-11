package handlers

import (
	"avito_tenders/internal/models"
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
)

func (h *handler) CreateBids() gin.HandlerFunc {
	return func(c *gin.Context) {

		bid := models.Bid{}

		jsonDataBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		err = json.Unmarshal(jsonDataBytes, &bid)
		if err != nil {
			c.JSON(500, gin.H{"error Unmarshal": err.Error()})
			return
		}

		if err := bid.Validate(); err != nil {
			c.JSON(400, gin.H{"Bad validation": err.Error()})
			return
		}
		
		err = h.storage.CreateBid(c, &bid)
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to create bid"})
			return
		}
		c.JSON(200, gin.H{"status": "ok"})
	}
}
