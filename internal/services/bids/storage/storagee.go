package storage

import (
	"avito_tenders/internal/db"
	"github.com/gin-gonic/gin"
)

type storage struct {
	Database *db.Database
}

// New создаёт новый сервис для работы с
func New(database *db.Database) *storage {
	return &storage{
		Database: database,
	}
}

// sd
func (s *storage) CreateBids(c *gin.Context) {
	
}
