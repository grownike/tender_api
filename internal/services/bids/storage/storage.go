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

func (s *storage) GetBidsByTender(c *gin.Context, tenderId uuid.UUID, username string, limit int, offset int) ([]models.Bid, error) {
	var bids []models.Bid

	err := s.Database.DB.
		Where("tender_id = ? AND creator_username = ?", tenderId, username).
		Order("name ASC").
		Limit(limit).
		Offset(offset).
		Find(&bids).Error

	if err != nil {
		return nil, err
	}

	return bids, nil
}

func (s *storage) GetBidStatus(c *gin.Context, bidId uuid.UUID, username string) (string, error) {
	var bid models.Bid

	err := s.Database.DB.
		Where("id = ? AND creator_username = ?", bidId, username).
		First(&bid).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", errors.New("bid not found")
	} else if err != nil {
		return "", err
	}

	return bid.Status, nil
}

func (s *storage) UpdateBid(c *gin.Context, bidId uuid.UUID, updates map[string]interface{}, username string) (*models.Bid, error) {
	var bid models.Bid
	if err := s.Database.DB.First(&bid, bidId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("bid not found")
		}
		return nil, err
	}
	if bid.CreatorUsername != username {
		return nil, errors.New("unauthorized: you are not the creator of this bid")
	}

	bidVersion := models.BidVersion{
		ID:          uuid.New(),
		BidID:       bid.ID,
		Name:        bid.Name,
		Description: bid.Description,
		Version:     bid.Version,
	}
	s.Database.DB.Create(&bidVersion)

	if err := s.Database.DB.Model(&bid).Updates(updates).Error; err != nil {
		return nil, err
	}
	s.Database.DB.Model(&bid).Update("version", bid.Version+1)

	return &bid, nil
}

func (s *storage) RollbackBid(bidId uuid.UUID, version models.BidVersion) error {

	var bid models.Bid

	if err := s.Database.DB.First(&bid, bidId).Error; err != nil {
		return errors.New("bid not found")
	}

	bidVersion := models.BidVersion{
		ID:          uuid.New(),
		BidID:       bid.ID,
		Name:        bid.Name,
		Description: bid.Description,
		Version:     bid.Version,
	}
	s.Database.DB.Create(&bidVersion)

	updates := map[string]interface{}{
		"name":         version.Name,
		"description":  version.Description,
		"version":      bid.Version + 1,
	}
	return s.Database.DB.Model(&models.Bid{}).Where("id = ?", bidId).Updates(updates).Error
}


func (s *storage) GetBidById(bidId uuid.UUID) (*models.Bid, error) {
	var bid models.Bid
	if err := s.Database.DB.First(&bid, bidId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("bid not found")
		}
		return nil, err
	}

	return &bid, nil
}

func (s *storage) GetBidVersion(bidId uuid.UUID, version int, bidVersion *models.BidVersion, username string) error {
	var bid models.Bid

	if err := s.Database.DB.First(&bid, bidId).Error; err != nil {
		return errors.New("bid not found")
	}

	if bid.CreatorUsername != username {
		return errors.New("unauthorized: you are not the creator of this bid")
	}

	return s.Database.DB.Where("bid_id = ? AND version = ?", bidId, version).First(bidVersion).Error
}


func (s *storage) EditBidStatus(bidId uuid.UUID, username string, newStatus string) (*models.Bid, error) {
	var bid models.Bid

	if err := s.Database.DB.Where("id = ?", bidId).First(&bid).Error; err != nil {
		return nil, errors.New("bid not found")
	}

	if bid.CreatorUsername != username {
		return nil, errors.New("unauthorized: you are not the creator of this bid")
	}

	bid.Status = newStatus

	if err := s.Database.DB.Save(&bid).Error; err != nil {
		return nil, err
	}

	return &bid, nil
}
