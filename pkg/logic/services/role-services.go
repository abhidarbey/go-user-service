package services

import (
	"github.com/abhidarbey/user-service/pkg/logic/helpers"
	"github.com/abhidarbey/user-service/pkg/models"
	"github.com/abhidarbey/user-service/pkg/repository/dao"
	"github.com/abhidarbey/user-service/pkg/utils/datetime"
	"github.com/abhidarbey/user-service/pkg/utils/encrypt"
	"github.com/abhidarbey/user-service/pkg/utils/resterr"
)

// RoleServiceInterface interface
type RoleServiceInterface interface {
	Create(models.Role, *models.AuthUser) (*models.Role, *resterr.RestErr)
	FindAll(*models.AuthUser) (models.Roles, *resterr.RestErr)
	GetByID(string, *models.AuthUser) (*models.Role, *resterr.RestErr)
	Update(models.Role, *models.AuthUser) (*models.Role, *resterr.RestErr)
	Delete(string, *models.AuthUser) *resterr.RestErr
}

type roleService struct{}

// RoleService variable
var (
	RoleService RoleServiceInterface = &roleService{}
)

// Create role
func (s *roleService) Create(role models.Role, au *models.AuthUser) (*models.Role, *resterr.RestErr) {
	// Validate role
	if err := role.Validate(); err != nil {
		return nil, err
	}

	// Verify permission --> IsGranted
	if !au.IsOrgAdmin {
		if !helpers.IsGranted("CanCreateRole", *au) {
			return nil, resterr.NewUnauthorizedError("Permission not granted")
		}
	}

	// Get department
	dept, err := DepartmentService.GetByID(role.Department, au)
	if err != nil {
		return nil, err
	}

	role.ID = "ROLE" + encrypt.GenerateID(17)
	role.Organization = dept.Organization
	role.Department = dept.ID
	role.Status = models.StatusActive
	role.IsActive = true
	role.CreatedAt = datetime.GetDateTimeString()
	role.UpdatedAt = datetime.GetDateTimeString()

	// Add role permissions
	if len(role.Permissions) > 0 {
		rolePermList, err := helpers.AssignRolePermissions(role)
		if err != nil {
			return nil, err
		}
		role.Permissions = *rolePermList
	}

	// Create new role
	newRole, err := dao.RoleDao.Create(role)
	if err != nil {
		return nil, err
	}

	return newRole, nil
}

// FindAll role
func (s *roleService) FindAll(au *models.AuthUser) (models.Roles, *resterr.RestErr) {
	// Verify permission --> IsGranted
	if !au.IsOrgAdmin {
		if !helpers.IsGranted("CanReadRole", *au) {
			return nil, resterr.NewUnauthorizedError("Permission not granted")
		}
	}

	roles, err := dao.RoleDao.FindAll(au.Organization)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// GetByID role
func (s *roleService) GetByID(id string, au *models.AuthUser) (*models.Role, *resterr.RestErr) {
	// Verify permission --> IsGranted
	if !au.IsOrgAdmin {
		if !helpers.IsGranted("CanReadRole", *au) {
			return nil, resterr.NewUnauthorizedError("Permission not granted")
		}
	}

	// Get role
	role, err := dao.RoleDao.GetByID(id, au.Organization)
	if err != nil {
		return nil, err
	}

	return role, nil
}

// Update role
func (s *roleService) Update(role models.Role, au *models.AuthUser) (*models.Role, *resterr.RestErr) {
	// Verify permission --> IsGranted
	if !au.IsOrgAdmin {
		if !helpers.IsGranted("CanUpdateRole", *au) {
			return nil, resterr.NewUnauthorizedError("Permission not granted")
		}
	}

	current, err := s.GetByID(role.ID, au)
	if err != nil {
		return nil, err
	}

	if role.Name != "" {
		current.Name = role.Name
	}

	if len(role.Permissions) > 0 {
		// role.Permissions = append(role.Permissions, current.Permissions...)
		// Validate permission request
		rolePermList, err := helpers.AssignRolePermissions(role)
		if err != nil {
			return nil, err
		}

		current.Permissions = *rolePermList
	}

	current.UpdatedAt = datetime.GetDateTimeString()

	if updateErr := dao.RoleDao.Update(*current); updateErr != nil {
		return nil, updateErr
	}

	return current, nil
}

// Delete role
func (s *roleService) Delete(id string, au *models.AuthUser) *resterr.RestErr {
	// Verify permission --> IsGranted
	if !au.IsOrgAdmin {
		if !helpers.IsGranted("CanDeleteRole", *au) {
			return resterr.NewUnauthorizedError("Permission not granted")
		}
	}

	return dao.RoleDao.Delete(id)
}
