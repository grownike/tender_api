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
func (s *storage) CreateEmployee(c *gin.Context, employee *models.Employee) error {
	//тут бы сделать проверку 
    return s.Database.DB.Create(employee).Error
}

// ready
func (s *storage) CreateCompany(c *gin.Context, company *models.Organization) error {
	
	return s.Database.DB.Create(company).Error
}

// ready
func (s *storage) AssignResponsible(organizationID uuid.UUID, username string) error {
	
	var organization models.Organization

	if err := s.Database.DB.Where("id = ?", organizationID).First(&organization).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("organization not found")
		}
		return err

	}

    var employee models.Employee
    if err := s.Database.DB.Where("username = ?", username).First(&employee).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

    
    responsible := models.Responsible{
        OrganizationID: organizationID,
        UserID:         employee.ID,
    }

    if err := s.Database.DB.Create(&responsible).Error; err != nil {
        return err
    }

    return nil
}
