package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const AUTOR_USER = "User"
const AUTOR_ORGANIZATION = "Organization"

type Employee struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (Employee) TableName() string {
	return "employee"
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

func (Organization) TableName() string {
	return "organization"
}

type Tender struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name            string    `gorm:"not null" json:"name"`
	Description     string    `json:"description"`
	ServiceType     string    `gorm:"not null;column:service_type" json:"serviceType"`
	Status          string    `gorm:"default:Created" json:"status"`
	OrganizationID  uuid.UUID `gorm:"not null" json:"organizationId"`
	CreatorUsername string    `gorm:"not null" json:"creatorUsername"`
	Version         int       `gorm:"default:1" json:"version"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (Tender) TableName() string {
	return "tender"
}

func (t *Tender) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("field 'name' is required")
	}
	if t.ServiceType == "" {
		return fmt.Errorf("field 'serviceType' is required")
	}
	if t.CreatorUsername == "" {
		return fmt.Errorf("field 'creatorUsername' is required")
	}

	if t.OrganizationID == uuid.Nil {
		return fmt.Errorf("field 'organizationId' must be a valid UUID")
	}

	return nil
}

type Bid struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name            string    `gorm:"not null" json:"name"`
	Description     string    `json:"description"`
	Status          string    `gorm:"default:Created" json:"status"`
	TenderID        uuid.UUID `gorm:"not null;type:uuid" json:"tenderId"`
	OrganizationID  uuid.UUID `gorm:"not null" json:"organizationId"`
	CreatorUsername string    `gorm:"not null" json:"creatorUsername"`
	AuthorType      string    `gorm:"not null" json:"authorType"`
	Version         int       `gorm:"default:1" json:"version"`
	AuthorId        uuid.UUID `gorm:"type:uuid" json:"authorId"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (b *Bid) Validate() error {
	if b.Name == "" {
		return fmt.Errorf("field 'name' is required")
	}
	if b.AuthorType != AUTOR_USER && b.AuthorType != AUTOR_ORGANIZATION {
		return fmt.Errorf("field 'authorType' is required")
	}
	if b.AuthorType == AUTOR_USER && b.CreatorUsername == "" {
		return fmt.Errorf("field 'creatorUsername' is required")
	}
	return nil
}

func (Bid) TableName() string {
	return "bid"
}

type Review struct {
	ID               uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Content          string    `gorm:"not null" json:"content"`
	BidID            uuid.UUID `gorm:"type:uuid;not null" json:"bidId"`
	ReviewerUsername string    `gorm:"not null" json:"reviewerUsername"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

func (Review) TableName() string {
	return "review"
}

type Responsible struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null" json:"organizationId"`
	UserID         uuid.UUID `gorm:"type:uuid;not null" json:"userId"`
}

func (Responsible) TableName() string {
	return "organization_responsible"
}

type TenderVersion struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	TenderID        uuid.UUID `gorm:"not null" json:"tenderId"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	ServiceType     string    `json:"serviceType"`
	OrganizationID  uuid.UUID `gorm:"not null" json:"organizationId"`
	Status          string    `gorm:"default:Created" json:"status"`
	Version         int       `json:"version"`
	CreatorUsername string    `gorm:"not null" json:"creatorUsername"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

func (TenderVersion) TableName() string {
	return "tender_version"
}

type BidVersion struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	BidID           uuid.UUID `gorm:"not null;type:uuid" json:"bidId"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Version         int       `json:"version"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"createdAt"`
	CreatorUsername string    `gorm:"not null" json:"creatorUsername"`
	AuthorId        uuid.UUID `gorm:"type:uuid" json:"authorId"`
	AuthorType      string    `gorm:"not null" json:"authorType"`
	Status          string    `gorm:"default:Created" json:"status"`
	TenderID        uuid.UUID `gorm:"not null;type:uuid" json:"tenderId"`
	OrganizationID  uuid.UUID `gorm:"not null" json:"organizationId"`
}

func (BidVersion) TableName() string {
	return "bid_versions"
}
