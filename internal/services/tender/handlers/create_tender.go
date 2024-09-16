package handlers

import (
	"avito_tenders/internal/models"
	"net/http"
	"github.com/gin-gonic/gin"
)

func (h *handler) CreateTender() gin.HandlerFunc {
	return func(c *gin.Context) {

		tender := models.Tender{}

		if err := c.ShouldBindJSON(&tender); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

		if err := tender.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		err := h.storage.CreateTender(c, &tender)
		if err != nil {
			if err.Error() == "user not found"{
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}
			if err.Error() == "organization not found" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}
			if err.Error() == "user is not responsible for this organization" {
				c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, tender)
	}
}
