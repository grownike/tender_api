package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Employee struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Organization struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Type        string    `gorm:"type:organization_type" json:"type"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Tender struct {
	gorm.Model
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name            string    `gorm:"not null" json:"name"`
	Description     string    `json:"description"`
	ServiceType     string    `json:"service_type"`
	Status          string    `gorm:"default:CREATED" json:"status"`
	OrganizationID  uint      `gorm:"not null" json:"organization_id"`
	CreatorUsername string    `gorm:"not null" json:"creator_username"`
	Version         int       `gorm:"default:1" json:"version"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Bid struct {
	gorm.Model
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name            string    `gorm:"not null" json:"name"`
	Description     string    `json:"description"`
	Status          string    `gorm:"default:CREATED" json:"status"`
	TenderID        uint      `gorm:"not null" json:"tender_id"`
	OrganizationID  uint      `gorm:"not null" json:"organization_id"`
	CreatorUsername string    `gorm:"not null" json:"creator_username"`
	Version         int       `gorm:"default:1" json:"version"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Review struct {
	gorm.Model
	ID               uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Content          string    `gorm:"not null" json:"content"`
	BidID            uint      `gorm:"not null" json:"bid_id"`
	ReviewerUsername string    `gorm:"not null" json:"reviewer_username"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type Responsible struct {
	gorm.Model
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	OrganizationID uint      `gorm:"not null" json:"organization_id"`
	UserID         uint      `gorm:"not null" json:"user_id"`
}
