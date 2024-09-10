package handlers

import (
	"avito_tenders/internal/models"
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
)

func (h *handler) CreateTender() gin.HandlerFunc {
	return func(c *gin.Context) {
		
		tender := models.Tender{}

		jsonDataBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(500, gin.H{"error1": err.Error()})
			return
		}
		err = json.Unmarshal(jsonDataBytes, &tender)
		if err != nil {
			c.JSON(500, gin.H{"error2": err.Error()})
			return
		}
		err = h.storage.CreateTender(c, &tender)
		if err != nil {
			c.JSON(500, gin.H{"error3": err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "ok"})
	}
}
