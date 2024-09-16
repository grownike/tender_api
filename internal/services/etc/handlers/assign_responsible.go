package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


func (h *handler) AssignResponsible() gin.HandlerFunc {
    return func(c *gin.Context) {
        organizationIDStr := c.Param("orgId")
        organizationID, err := uuid.Parse(organizationIDStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization ID"})
            return
        }


		username := c.Param("user")
        
        err = h.storage.AssignResponsible(organizationID, username)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Responsible assigned successfully"})
    }
}