package services

import (
	"github.com/abhidarbey/user-service/pkg/logic/helpers"
	"github.com/abhidarbey/user-service/pkg/models"
	"github.com/abhidarbey/user-service/pkg/repository/dao"
	"github.com/abhidarbey/user-service/pkg/utils/datetime"
	"github.com/abhidarbey/user-service/pkg/utils/resterr"
)

// OrganizationServiceInterface interface
type OrganizationServiceInterface interface {
	FindAll(*models.AuthUser) (models.Organizations, *resterr.RestErr)
	GetByID(string, *models.AuthUser) (*models.Organization, *resterr.RestErr)
	Update(models.Organization, *models.AuthUser) (*models.Organization, *resterr.RestErr)
	Delete(string, *models.AuthUser) *resterr.RestErr
}

type organizationService struct{}

// OrganizationService variable
var (
	OrganizationService OrganizationServiceInterface = &organizationService{}
)

// FindAll organization
func (s *organizationService) FindAll(au *models.AuthUser) (models.Organizations, *resterr.RestErr) {
	organizations, err := dao.OrganizationDao.FindAll()
	if err != nil {
		return nil, err
	}

	return organizations, nil
}

// GetByID organization
func (s *organizationService) GetByID(id string, au *models.AuthUser) (*models.Organization, *resterr.RestErr) {
	// Get organization
	organization, err := dao.OrganizationDao.GetByID(id)
	if err != nil {
		return nil, err
	}

	return organization, nil
}

// Update organization
func (s *organizationService) Update(organization models.Organization, au *models.AuthUser) (*models.Organization, *resterr.RestErr) {
	// Verify permission --> IsGranted
	if !au.IsSuperuser && !au.IsOrgAdmin {
		if !helpers.IsGranted("CanUpdateOrganization", *au) {
			return nil, resterr.NewUnauthorizedError("Permission not granted")
		}
	}

	current, err := s.GetByID(organization.ID, au)
	if err != nil {
		return nil, err
	}

	// Verify auth user's organization
	if au.Organization != current.ID {
		return nil, resterr.NewUnauthorizedError("Unauthorized request")
	}

	if organization.Name != "" {
		current.Name = organization.Name
	}
	if organization.Website != "" {
		current.Website = organization.Website
	}

	current.UpdatedAt = datetime.GetDateTimeString()

	if updateErr := dao.OrganizationDao.Update(*current); updateErr != nil {
		return nil, updateErr
	}

	return current, nil
}

// Delete organization
func (s *organizationService) Delete(id string, au *models.AuthUser) *resterr.RestErr {
	// Verify permission --> IsGranted
	if !au.IsOrgAdmin {
		if !helpers.IsGranted("CanDeleteOrganization", *au) {
			return resterr.NewUnauthorizedError("Permission not granted")
		}
	}

	return dao.OrganizationDao.Delete(id)
}
