package services

import (
	"github.com/abhidarbey/user-service/pkg/logic/helpers"
	"github.com/abhidarbey/user-service/pkg/models"
	"github.com/abhidarbey/user-service/pkg/repository/dao"
	"github.com/abhidarbey/user-service/pkg/utils/datetime"
	"github.com/abhidarbey/user-service/pkg/utils/encrypt"
	"github.com/abhidarbey/user-service/pkg/utils/resterr"
)

// DepartmentServiceInterface interface
type DepartmentServiceInterface interface {
	Create(models.Department, *models.AuthUser) (*models.Department, *resterr.RestErr)
	FindAll(*models.AuthUser) (models.Departments, *resterr.RestErr)
	GetByID(string, *models.AuthUser) (*models.Department, *resterr.RestErr)
	Update(models.Department, *models.AuthUser) (*models.Department, *resterr.RestErr)
	Delete(string, *models.AuthUser) *resterr.RestErr
}

type departmentService struct{}

// DepartmentService variable
var (
	DepartmentService DepartmentServiceInterface = &departmentService{}
)

// Create department
func (s *departmentService) Create(dept models.Department, au *models.AuthUser) (*models.Department, *resterr.RestErr) {
	// Validate request
	if err := dept.Validate(); err != nil {
		return nil, err
	}

	// Verify permission --> IsGranted
	if !au.IsOrgAdmin {
		if !helpers.IsGranted("CanCreateDepartment", *au) {
			return nil, resterr.NewUnauthorizedError("Permission not granted")
		}
	}

	dept.ID = "DEPT" + encrypt.GenerateID(17)
	dept.Organization = au.Organization
	dept.Status = models.StatusActive
	dept.IsActive = true
	dept.CreatedAt = datetime.GetDateTimeString()
	dept.UpdatedAt = datetime.GetDateTimeString()

	newDept, err := dao.DepartmentDao.Create(dept)
	if err != nil {
		return nil, err
	}
	return newDept, nil
}

// FindAll department
func (s *departmentService) FindAll(au *models.AuthUser) (models.Departments, *resterr.RestErr) {
	// Verify permission --> IsGranted
	if !au.IsOrgAdmin {
		if !helpers.IsGranted("CanReadDepartment", *au) {
			return nil, resterr.NewUnauthorizedError("Permission not granted")
		}
	}

	departments, err := dao.DepartmentDao.FindAll(au.Organization)
	if err != nil {
		return nil, err
	}

	return departments, nil
}

// GetByID department
func (s *departmentService) GetByID(id string, au *models.AuthUser) (*models.Department, *resterr.RestErr) { // Verify permission --> IsGranted
	if !au.IsOrgAdmin {
		if !helpers.IsGranted("CanReadDepartment", *au) {
			return nil, resterr.NewUnauthorizedError("Permission not granted")
		}
	}

	// Get department
	department, err := dao.DepartmentDao.GetByID(id, au.Organization)
	if err != nil {
		return nil, err
	}

	return department, nil
}

// Update department
func (s *departmentService) Update(department models.Department, au *models.AuthUser) (*models.Department, *resterr.RestErr) { // Verify permission --> IsGranted
	if !au.IsOrgAdmin {
		if !helpers.IsGranted("CanUpdateDepartment", *au) {
			return nil, resterr.NewUnauthorizedError("Permission not granted")
		}
	}

	current, err := s.GetByID(department.ID, au)
	if err != nil {
		return nil, err
	}

	// Verify auth user's organization
	if au.Organization != current.Organization {
		return nil, resterr.NewUnauthorizedError("Unauthorized request")
	}

	if department.Name != "" {
		current.Name = department.Name
	}

	current.UpdatedAt = datetime.GetDateTimeString()

	if updateErr := dao.DepartmentDao.Update(*current); updateErr != nil {
		return nil, updateErr
	}

	return current, nil
}

// Delete department
func (s *departmentService) Delete(id string, au *models.AuthUser) *resterr.RestErr { // Verify permission --> IsGranted
	if !au.IsOrgAdmin {
		if !helpers.IsGranted("CanDeleteDepartment", *au) {
			return resterr.NewUnauthorizedError("Permission not granted")
		}
	}

	return dao.DepartmentDao.Delete(id)
}
