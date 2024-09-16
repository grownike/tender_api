package handlers

import (
	"avito_tenders/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) CreateEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		var employee models.Employee
		if err := c.ShouldBindJSON(&employee); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := h.storage.CreateEmployee(c, &employee); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, employee)
	}
}