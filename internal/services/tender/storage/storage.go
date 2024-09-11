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

	tenderVersion := models.TenderVersion{
		TenderID:    tender.ID,
		Name:        tender.Name,
		Description: tender.Description,
		ServiceType: tender.ServiceType,
		Version:     tender.Version,
	}
	s.Database.DB.Create(&tenderVersion)

	if err := s.Database.DB.Model(&tender).Updates(updates).Error; err != nil {
		return nil, err
	}
	s.Database.DB.Model(&tender).Update("version", tender.Version+1)

	return &tender, nil
}

func (s *storage) GetTenderById(tenderId uuid.UUID) (*models.Tender, error) {
	var tender models.Tender
	if err := s.Database.DB.First(&tender, tenderId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tender not found")
		}
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

func (s *storage) GetTenderVersion(tenderId uuid.UUID, version int, tenderVersion *models.TenderVersion, username string) error {
	var tender models.Tender

	if err := s.Database.DB.First(&tender, tenderId).Error; err != nil {
		return errors.New("tender not found")
	}

	if tender.CreatorUsername != username {
		return errors.New("unauthorized: you are not the creator of this tender")
	}

	return s.Database.DB.Where("tender_id = ? AND version = ?", tenderId, version).First(tenderVersion).Error
}
func (s *storage) RollbackTender(tenderId uuid.UUID, version models.TenderVersion) error {

	var tender models.Tender

	if err := s.Database.DB.First(&tender, tenderId).Error; err != nil {
		return errors.New("tender not found")
	}

	updates := map[string]interface{}{
		"name":        version.Name,
		"description": version.Description,
		"service_type": version.ServiceType,
		"version":     tender.Version + 1,
	}
	return s.Database.DB.Model(&models.Tender{}).Where("id = ?", tenderId).Updates(updates).Error
}


