package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *handler) EditTender() gin.HandlerFunc {
	return func(c *gin.Context) {

		tenderIdParam := c.Param("tenderId")
		username := c.Query("username")

		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
			return
		}

		tenderId, err := uuid.Parse(tenderIdParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tender ID"})
			return
		}

		var updates map[string]interface{}

		if err := c.ShouldBindJSON(&updates); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		updates = convertKeys(updates)

		updatedTender, err := h.storage.UpdateTender(c, tenderId, updates, username)
		if err != nil {
			if err.Error() == "user not found" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}
			if err.Error() == "unauthorized: user is not responsible for tender's organization" {
				c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
				return
			}
			if err.Error() == "tender not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			if err.Error() == "tender is not in CREATED status" {
				c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tender"})
			return
		}

		c.JSON(http.StatusOK, updatedTender)
	}
}

//Я столкнулся с тем, что при передачи в json "serviceType" у меня почему-то он так и остаётся.
//В create_tender всё работает как надо. Пока что я решил проблему через такой костыль...

func convertKeys(updates map[string]interface{}) map[string]interface{} {
	converted := make(map[string]interface{})
	for key, value := range updates {
		switch key {
		case "serviceType":
			converted["service_type"] = value
		default:
			converted[key] = value
		}
	}
	return converted
}
