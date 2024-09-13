package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *handler) RollbackTender() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenderIdParam := c.Param("tenderId")
		versionParam := c.Param("version")
		username := c.Query("username")

		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
			return
		}

		tenderId, err := uuid.Parse(tenderIdParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tender ID"})
			return
		}

		version, err := strconv.Atoi(versionParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version"})
			return
		}

	
		tender,err := h.storage.RollbackTender(tenderId, version, username)
		if err != nil {
			if err.Error() == "tender not found"{
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			if err.Error() == "version not found"{
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			if err.Error() == "user not found" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}
			if err.Error() == "unauthorized: user is not responsible for this organization" {
				c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to rollback tender"})
			return
		}

		c.JSON(http.StatusOK, tender)
	}
}
