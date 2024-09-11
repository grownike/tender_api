package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *handler) EditTender() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем ID тендера из параметров URL
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


		jsonDataBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		var updates map[string]interface{}

		err = json.Unmarshal(jsonDataBytes, &updates)
		if err != nil {
			c.JSON(500, gin.H{"error Unmarshal": err.Error()})
			return
		}

		updates = convertKeys(updates)

		updatedTender, err := h.storage.UpdateTender(c, tenderId, updates, username)
		if err != nil {
			if err.Error() == "unauthorized: you are not the creator of this tender" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}
			if err.Error() == "tender not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return

			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tender"})
			return
		}

		c.JSON(http.StatusOK, updatedTender)
	}
}


//Я столкнулся с тем, что при передачи в json "serviceType" у меня почему-то он так и остаётся.
//В create_tendeer всё работает как надо. Пока что я решил проблему через такой костыль...
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
