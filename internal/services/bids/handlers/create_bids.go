package handlers

import (
	"github.com/gin-gonic/gin"
)

func (h *handler) CreateBids() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.storage.CreateBids(c)
		c.JSON(200, gin.H{"status": "ok"})
	}
}
