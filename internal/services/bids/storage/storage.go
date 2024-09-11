package storage

import (
	"avito_tenders/internal/db"
	"avito_tenders/internal/models"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

func (s *storage) CreateBid(c *gin.Context, bid *models.Bid) error {

	if bid.AuthorType == models.AUTOR_ORGANIZATION {
		var org models.Organization
		if err := s.Database.DB.First(&org, bid.AuthorId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("org not found")
			}
		}
		bid.CreatorUsername = org.Name
		bid.OrganizationID = org.ID
	}

	if bid.AuthorType == models.AUTOR_USER {
		var resp models.Responsible
		if err := s.Database.DB.First(&resp, "user_id = ?", bid.AuthorId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("responsibility record not found for user")
			}
			return err
		}

		var org models.Organization
		if err := s.Database.DB.First(&org, "id = ?", resp.OrganizationID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("organization not found")
			}
			return err
		}

		bid.OrganizationID = org.ID
	}

	if err := s.Database.DB.Create(bid).Error; err != nil {
		return err
	}
	return nil
}

func (s *storage) GetMyBids(username string) *gorm.DB {
	query := s.Database.DB.Where("creator_username = ?", username).Order("name ASC")
	return query
}
