package handlers

import (
	"avito_tenders/internal/models"
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

		var tenderVersion models.TenderVersion
		if err := h.storage.GetTenderVersion(tenderId, version, &tenderVersion, username); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Version not found"})
			return
		}

		if err := h.storage.RollbackTender(tenderId, tenderVersion); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to rollback tender"})
			return
		}
		tender, err := h.storage.GetTenderById(tenderId)
		if err != nil {
			if err.Error() == "tender not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return

			}
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version"})
			return
		}

		c.JSON(http.StatusOK, tender)
	}
}
