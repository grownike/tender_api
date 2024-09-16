package handlers

import (
	"avito_tenders/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type storage interface {
	CreateEmployee(c *gin.Context, employee *models.Employee) error
	CreateCompany(c *gin.Context, employee *models.Organization) error
	AssignResponsible(organizationID uuid.UUID, username string) error
}
