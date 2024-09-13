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

// ready
func (s *storage) CreateTender(c *gin.Context, tender *models.Tender) error {
	var organization models.Organization
	var responsible models.Responsible
	var employee models.Employee

	if err := s.Database.DB.Where("username = ?", tender.CreatorUsername).First(&employee).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	if err := s.Database.DB.Where("id = ?", tender.OrganizationID).First(&organization).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("organization not found")
		}
		return err

	}

	if err := s.Database.DB.
		Joins("JOIN employee ON employee.id = organization_responsible.user_id").
		Where("employee.username = ? AND organization_responsible.organization_id = ?", tender.CreatorUsername, tender.OrganizationID).
		First(&responsible).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user is not responsible for this organization")
		}
		return err
	}

	if err := s.Database.DB.Create(tender).Error; err != nil {
		return err
	}
	return nil
}

// ready
func (s *storage) GetTenders(limit, offset int, serviceType []string) *gorm.DB {
	query := s.Database.DB.Limit(limit).Offset(offset).Order("name ASC")
	query = query.Where("status = ?", "PUBLISHED")

	if len(serviceType) > 0 {
		query = query.Where("service_type IN ?", serviceType)
	}
	return query
}

// ready
func (s *storage) GetMyTenders(username string, limit int, offset int) ([]models.Tender, error) {
	employee := models.Employee{}

	if err := s.Database.DB.Where("username = ?", username).First(&employee).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	var tenders []models.Tender
	if err := s.Database.DB.Where("creator_username = ?", username).Order("name ASC").
		Limit(limit).
		Offset(offset).Find(&tenders).Error; err != nil {
		return nil, err
	}

	return tenders, nil
}

// ready
func (s *storage) UpdateTender(c *gin.Context, tenderId uuid.UUID, updates map[string]interface{}, username string) (*models.Tender, error) {
	var employee models.Employee

	if err := s.Database.DB.Where("username = ?", username).First(&employee).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	var tender models.Tender
	if err := s.Database.DB.First(&tender, tenderId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tender not found")
		}
		return nil, err
	}

	if tender.Status != "Created" {
		return nil, errors.New("tender is not in Created status")
	}

	var responsible models.Responsible

	if err := s.Database.DB.
		Where("organization_id = ? AND user_id = ?", tender.OrganizationID, employee.ID).
		First(&responsible).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("unauthorized: user is not responsible for tender's organization")
		}
		return nil, err
	}

	tenderVersion := models.TenderVersion{
		TenderID:        tenderId,
		Name:            tender.Name,
		Description:     tender.Description,
		ServiceType:     tender.Description,
		Status:          tender.Status,
		OrganizationID:  tender.OrganizationID,
		CreatorUsername: tender.CreatorUsername,
		Version:         tender.Version,
		CreatedAt:       tender.CreatedAt,
	}
	s.Database.DB.Create(&tenderVersion)

	if err := s.Database.DB.Model(&tender).Updates(updates).Error; err != nil {
		return nil, err
	}
	s.Database.DB.Model(&tender).Update("version", tender.Version+1)

	return &tender, nil
}

// ready work
func (s *storage) GetTenderStatus(tenderId uuid.UUID, username string) (string, error) {
	var tender models.Tender
	var employee models.Employee

	if err := s.Database.DB.First(&tender, tenderId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("tender not found")
		}
		return "", err
	}

	if err := s.Database.DB.Where("username = ?", username).First(&employee).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("user not found")
		}
		return "", err
	}

	var responsible models.Responsible

	if err := s.Database.DB.
		Where("organization_id = ? AND user_id = ?", tender.OrganizationID, employee.ID).
		First(&responsible).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("unauthorized: you are not responsible for this organization's tenders")
		}
		return "", err
	}

	return tender.Status, nil
}

// ready
func (s *storage) EditTenderStatus(tenderId uuid.UUID, username string, newStatus string) (*models.Tender, error) {
	var tender models.Tender
	var employee models.Employee
	var responsible models.Responsible

	if err := s.Database.DB.Where("username = ?", username).First(&employee).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if err := s.Database.DB.First(&tender, tenderId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tender not found")
		}
		return nil, err
	}

	if err := s.Database.DB.
		Where("organization_id = ? AND user_id = ?", tender.OrganizationID, employee.ID).
		First(&responsible).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("unauthorized: user is not responsible for this organization")
		}
		return nil, err
	}

	tender.Status = newStatus

	if err := s.Database.DB.Save(&tender).Error; err != nil {
		return nil, err
	}

	return &tender, nil
}

// ready
func (s *storage) RollbackTender(tenderId uuid.UUID, version int, username string) (*models.Tender, error) {

	var tender models.Tender
	var employee models.Employee
	var responsible models.Responsible

	if err := s.Database.DB.Where("username = ?", username).First(&employee).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if err := s.Database.DB.First(&tender, tenderId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tender not found")
		}
		return nil, err
	}

	if err := s.Database.DB.
		Where("organization_id = ? AND user_id = ?", tender.OrganizationID, employee.ID).
		First(&responsible).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("unauthorized: user is not responsible for this organization")
		}
		return nil, err
	}

	if tender.Status != "Created" {
		return nil, errors.New("tender is not in Created status")
	}

	var tenderRollback models.TenderVersion
	if err := s.Database.DB.Where("tender_id = ? AND version = ?", tenderId, version).First(&tenderRollback).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("version not found")
		}
		return nil, err

	}

	tenderVersion := models.TenderVersion{
		TenderID:        tenderId,
		Name:            tender.Name,
		Description:     tender.Description,
		ServiceType:     tender.Description,
		Status:          tender.Status,
		OrganizationID:  tender.OrganizationID,
		CreatorUsername: tender.CreatorUsername,
		Version:         tender.Version,
		CreatedAt:       tender.CreatedAt,
	}

	s.Database.DB.Create(&tenderVersion)

	updates := map[string]interface{}{
		"name":             tenderRollback.Name,
		"description":      tenderRollback.Description,
		"service_type":     tenderRollback.Description,
		"status":           tenderRollback.Status,
		"organization_id":  tenderRollback.OrganizationID,
		"creator_username": tenderRollback.CreatorUsername,
		"version":          tender.Version + 1,
		"created_at":       tenderRollback.CreatedAt,
	}
	if err := s.Database.DB.Model(&models.Tender{}).Where("id = ?", tenderId).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &tender, nil
}
