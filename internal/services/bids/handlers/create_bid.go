package handlers

import (
	"avito_tenders/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) CreateBids() gin.HandlerFunc {
	return func(c *gin.Context) {

		bid := models.Bid{}

		if err := c.ShouldBindJSON(&bid); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

		if err := bid.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		if err := h.storage.CreateBid(c, &bid); err!= nil{
			if err.Error() == "user not found"{
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}
			if err.Error() == "organization not found"{
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}
			if err.Error() == "user is not responsible for any organization"{
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}
			if err.Error() == "user is responsible for organization that not found"{
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}
			if err.Error() == "tender not found"{
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, bid)
	}
}
