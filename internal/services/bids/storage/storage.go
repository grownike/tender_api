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
func (s *storage) CreateBid(c *gin.Context, bid *models.Bid) error {
	var tender models.Tender
	if err := s.Database.DB.First(&tender, bid.TenderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("tender not found")
		}
	}

	if bid.AuthorType == models.AUTOR_ORGANIZATION {
		var organization models.Organization
		if err := s.Database.DB.First(&organization, bid.AuthorId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("organization not found")
			}
		}
		bid.CreatorUsername = organization.Name
		bid.OrganizationID = organization.ID
	}

	if bid.AuthorType == models.AUTOR_USER {
		var employee models.Employee

		if err := s.Database.DB.Where("username = ?", bid.CreatorUsername).First(&employee).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("user not found")
			}
			return err
		}

		var responsible models.Responsible
		if err := s.Database.DB.First(&responsible, "user_id = ?", employee.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("user is not responsible for any organization")
			}
			return err
		}
		var organization models.Organization
		if err := s.Database.DB.First(&organization, "id = ?", responsible.OrganizationID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("user is responsible for organization that not found")
			}
			return err
		}

		bid.OrganizationID = organization.ID
	}

	if err := s.Database.DB.Create(bid).Error; err != nil {
		return err
	}
	return nil
}

// ready
func (s *storage) GetMyBids(username string, limit int, offset int) ([]models.Bid, error) {
	employee := models.Employee{}

	if err := s.Database.DB.Where("username = ?", username).First(&employee).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	var bids []models.Bid
	if err := s.Database.DB.Where("creator_username = ?", username).Order("name ASC").
		Limit(limit).
		Offset(offset).Find(&bids).Error; err != nil {
		return nil, err
	}

	return bids, nil
}

// ready
func (s *storage) GetBidsByTender(c *gin.Context, tenderId uuid.UUID, username string, limit int, offset int) ([]models.Bid, error) {
	var bids []models.Bid
	var tender models.Tender
	var employee models.Employee

	if err := s.Database.DB.First(&tender, tenderId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tender not found")
		}
		return nil, err
	}

	if err := s.Database.DB.Where("username = ?", username).First(&employee).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if tender.Status != "published" {
		var responsible models.Responsible
		if err := s.Database.DB.
			Where("organization_id = ? AND user_id = ?", tender.OrganizationID, employee.ID).
			First(&responsible).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("unauthorized: user in not responsible of this tender's bids")
			}
			return nil, err
		}
	}

	err := s.Database.DB.
		Where("tender_id = ?", tenderId).
		Order("name ASC").
		Limit(limit).
		Offset(offset).
		Find(&bids).Error

	if err != nil {
		return nil, err
	}

	return bids, nil
}

// ready
func (s *storage) GetBidStatus(c *gin.Context, bidId uuid.UUID, username string) (string, error) {
	var bid models.Bid
	var employee models.Employee

	if err := s.Database.DB.First(&bid, bidId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("bid not found")
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
		Where("organization_id = ? AND user_id = ?", bid.OrganizationID, employee.ID).
		First(&responsible).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("unauthorized: you are not responsible for this bid")
		}
		return "", err
	}

	return bid.Status, nil
}

// ready
func (s *storage) UpdateBid(c *gin.Context, bidId uuid.UUID, updates map[string]interface{}, username string) (*models.Bid, error) {
	var employee models.Employee

	if err := s.Database.DB.Where("username = ?", username).First(&employee).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	var bid models.Bid
	if err := s.Database.DB.First(&bid, bidId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("bid not found")
		}
		return nil, err
	}

	if bid.Status != "Created" {
		return nil, errors.New("bid is not in Created status")
	}

	var responsible models.Responsible

	if err := s.Database.DB.
		Where("organization_id = ? AND user_id = ?", bid.OrganizationID, employee.ID).
		First(&responsible).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("unauthorized: user is not responsible for bid's organization")
		}
		return nil, err
	}

	bidVersion := models.BidVersion{
		BidID:           bid.ID,
		Name:            bid.Name,
		Description:     bid.Description,
		Status:          bid.Status,
		TenderID:        bid.TenderID,
		CreatorUsername: bid.CreatorUsername,
		Version:         bid.Version,
		AuthorId:        bid.AuthorId,
		AuthorType:      bid.AuthorType,
		CreatedAt:       bid.CreatedAt,
		OrganizationID:  bid.OrganizationID,
	}

	s.Database.DB.Create(&bidVersion)

	if err := s.Database.DB.Model(&bid).Updates(updates).Error; err != nil {
		return nil, err
	}
	s.Database.DB.Model(&bid).Update("version", bid.Version+1)

	return &bid, nil
}

// ready
func (s *storage) RollbackBid(bidId uuid.UUID, version int, username string) (*models.Bid, error) {

	var bid models.Bid
	var employee models.Employee
	var responsible models.Responsible

	if err := s.Database.DB.Where("username = ?", username).First(&employee).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if err := s.Database.DB.First(&bid, bidId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("bid not found")
		}
		return nil, err
	}

	if err := s.Database.DB.
		Where("organization_id = ? AND user_id = ?", bid.OrganizationID, employee.ID).
		First(&responsible).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("unauthorized: user is not responsible for this bid's organization")
		}
		return nil, err
	}

	if bid.Status != "Created" {
		return nil, errors.New("bid is not in Created status")
	}

	var bidRollback models.BidVersion
	if err := s.Database.DB.Where("bid_id = ? AND version = ?", bidId, version).First(&bidRollback).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("version not found")
		}
		return nil, err

	}

	bidVersion := models.BidVersion{
		BidID:           bid.ID,
		Name:            bid.Name,
		Description:     bid.Description,
		Status:          bid.Status,
		TenderID:        bid.TenderID,
		CreatorUsername: bid.CreatorUsername,
		Version:         bid.Version,
		AuthorId:        bid.AuthorId,
		AuthorType:      bid.AuthorType,
		CreatedAt:       bid.CreatedAt,
		OrganizationID:  bid.OrganizationID,
	}

	s.Database.DB.Create(&bidVersion)

	updates := map[string]interface{}{
		"name":             bidRollback.Name,
		"description":      bidRollback.Description,
		"status":           bidRollback.Status,
		"tender_id":        bidRollback.TenderID,
		"creator_username": bidRollback.CreatorUsername,
		"version":          bid.Version + 1,
		"author_id":        bidRollback.AuthorId,
		"author_type":      bidRollback.AuthorType,
		"created_at":       bidRollback.CreatedAt,
		"organization_id":  bidRollback.OrganizationID,
	}

	if err := s.Database.DB.Model(&models.Bid{}).Where("id = ?", bidId).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &bid, nil
}

func (s *storage) EditBidStatus(bidId uuid.UUID, username string, newStatus string) (*models.Bid, error) {
	var bid models.Bid
	var employee models.Employee
	var responsible models.Responsible

	if err := s.Database.DB.Where("username = ?", username).First(&employee).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if err := s.Database.DB.First(&bid, bidId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("bid not found")
		}
		return nil, err
	}

	if err := s.Database.DB.
		Where("organization_id = ? AND user_id = ?", bid.OrganizationID, employee.ID).
		First(&responsible).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("unauthorized: user is not responsible for this bid's organization")
		}
		return nil, err
	}

	bid.Status = newStatus

	if err := s.Database.DB.Save(&bid).Error; err != nil {
		return nil, err
	}

	return &bid, nil
}

func (s *storage) GetReviewsByTender(tenderId uuid.UUID, limit int, offset int, authorUsername string, requesterUsername string) ([]models.Review, error) {
	var reviews []models.Review
	var author models.Employee
	var requester models.Employee
	var tender models.Tender
	var responsible models.Responsible

	if err := s.Database.DB.Where("username = ?", authorUsername).First(&author).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("author not found")
		}
		return nil, err
	}

	if err := s.Database.DB.Where("username = ?", requesterUsername).First(&requester).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("requester not found")
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
		Where("organization_id = ? AND user_id = ?", tender.OrganizationID, requester.ID).
		First(&responsible).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("unauthorized: requester is not responsible for tender's organization")
		}
		return nil, err
	}

	err := s.Database.DB.
		Where("bid_id IN (SELECT id FROM bid WHERE tender_id = ? AND creator_username = ?)", tenderId, authorUsername).
		Limit(limit).
		Offset(offset).
		Find(&reviews).Error

	if err != nil {
		return nil, err
	}

	return reviews, nil
}
