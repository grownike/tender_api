package handlers

import (
	"github.com/gin-gonic/gin"
)

type storage interface {
	CreateBids(c *gin.Context)
}
