package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *handler) GetTenderStatus() gin.HandlerFunc {
	return func(c *gin.Context) {

		tenderIdParam := c.Param("tenderId")
		username := c.Query("username")

		print(username)

		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
			return
		}

		tenderId, err := uuid.Parse(tenderIdParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tender ID"})
			return
		}


		status, err := h.storage.GetTenderStatus(tenderId, username)
		if err != nil {
			if err.Error() == "unauthorized: you are not the creator of this tender" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}
			if err.Error() == "tender not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return

			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tender status"})
			return
		}

		
		c.JSON(http.StatusOK, status)
	}
}
