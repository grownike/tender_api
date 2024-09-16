package handlers

import (
	"avito_tenders/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) CreateCompany() gin.HandlerFunc {
	return func(c *gin.Context) {
		var company models.Organization
		if err := c.ShouldBindJSON(&company); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := h.storage.CreateCompany(c, &company); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, company)
	}
}