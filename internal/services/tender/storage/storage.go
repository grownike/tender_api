package storage

import (
	"avito_tenders/internal/db"
	"avito_tenders/internal/models"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type storage struct {
	Database *db.Database
}

func New(database *db.Database) *storage {
	return &storage{
		Database: database,
	}
}

func (s *storage) CreateTender(c *gin.Context, tender *models.Tender) error {
	if err := s.Database.DB.Create(tender).Error; err != nil {
		return err
	}
	return nil
}

func (s *storage) GetTenders(limit, offset int, serviceType []string) *gorm.DB {
	query := s.Database.DB.Limit(limit).Offset(offset).Order("name ASC")
	if len(serviceType) > 0 {
		query = query.Where("service_type IN ?", serviceType)
	}
	return query
}

func (s *storage) GetMyTenders(username string) *gorm.DB {
	query := s.Database.DB.Where("creator_username = ?", username).Order("name ASC")
	return query
}

func (s *storage) UpdateTender(c *gin.Context, tenderId uuid.UUID, updates map[string]interface{}, username string) (*models.Tender, error) {
	var tender models.Tender
	if err := s.Database.DB.First(&tender, tenderId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tender not found")
		}
		return nil, err
	}
	if tender.CreatorUsername != username {
		return nil, errors.New("unauthorized: you are not the creator of this tender")
	}

	if err := s.Database.DB.Model(&tender).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &tender, nil
}

func (s *storage) GetTenderStatus(tenderId uuid.UUID, username string) (string, error) {
	var tender models.Tender
	if err := s.Database.DB.First(&tender, tenderId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("tender not found")
		}
		return "", err
	}

	if tender.CreatorUsername != username {
		return "", errors.New("unauthorized: you are not the creator of this tender")
	}

	return tender.Status, nil
}

func (s *storage) EditTenderStatus(tenderId uuid.UUID, username string, newStatus string) (*models.Tender, error) {
	var tender models.Tender

	if err := s.Database.DB.Where("id = ?", tenderId).First(&tender).Error; err != nil {
		return nil, errors.New("tender not found")
	}

	if tender.CreatorUsername != username {
		return nil, errors.New("unauthorized: you are not the creator of this tender")
	}

	tender.Status = newStatus

	if err := s.Database.DB.Save(&tender).Error; err != nil {
		return nil, err
	}

	return &tender, nil
}
