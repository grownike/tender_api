package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *handler) GetReviewsByTender() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenderIdParam := c.Param("tenderId")

		tenderId, err := uuid.Parse(tenderIdParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tender ID"})
			return
		}

		// Optional pagination parameters
		limitParam := c.DefaultQuery("limit", "5")
		offsetParam := c.DefaultQuery("offset", "0")
		authorUsername := c.Query("authorUsername")
		requesterUsername := c.Query("requesterUsername")

		if authorUsername == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "authorUsername is required"})
			return
		}

		if authorUsername == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "requesterUsername is required"})
			return
		}

		limit, err := strconv.Atoi(limitParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit value"})
			return
		}

		offset, err := strconv.Atoi(offsetParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset value"})
			return
		}

		// Call storage function to retrieve reviews
		reviews, err := h.storage.GetReviewsByTender(tenderId, limit, offset, authorUsername, requesterUsername)
		if err != nil {
			if err.Error() == "unauthorized: requester is not responsible for tender's organization" {
				c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
				return
			}
			if err.Error() == "tender not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			if err.Error() == "author not found" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}
			if err.Error() == "requester not found" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get review"})
			return
		}

		if len(reviews) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "No reviews found"})
			return
		}

		c.JSON(http.StatusOK, reviews)
	}
}
