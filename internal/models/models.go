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
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

type Organization struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Type        string    `gorm:"type:organization_type" json:"type"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

type Tender struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name            string    `gorm:"not null" json:"name"`
	Description     string    `json:"description"`
	ServiceType     string    `gorm:"not null" json:"serviceType"`
	Status          string    `gorm:"default:CREATED" json:"status"`
	OrganizationID  uuid.UUID `gorm:"not null" json:"organizationId"` 
	CreatorUsername string    `gorm:"not null" json:"creatorUsername"` 
	Version         int       `gorm:"default:1" json:"version"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

type Bid struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name            string    `gorm:"not null" json:"name"`
	Description     string    `json:"description"`
	Status          string    `gorm:"default:CREATED" json:"status"`
	TenderID        uint      `gorm:"not null" json:"tenderId"`
	OrganizationID  uint      `gorm:"not null" json:"organizationId"`
	CreatorUsername string    `gorm:"not null" json:"creatorUsername"`
	Version         int       `gorm:"default:1" json:"version"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

type Review struct {
	ID               uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Content          string    `gorm:"not null" json:"content"`
	BidID            uint      `gorm:"not null" json:"bidId"`
	ReviewerUsername string    `gorm:"not null" json:"reviewerUsername"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

type Responsible struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	OrganizationID uint      `gorm:"not null" json:"organizationId"`
	UserID         uint      `gorm:"not null" json:"userId"`
}
