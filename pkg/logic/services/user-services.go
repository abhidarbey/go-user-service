package services

import (
	"github.com/abhidarbey/user-service/pkg/logic/helpers"
	"github.com/abhidarbey/user-service/pkg/models"
	"github.com/abhidarbey/user-service/pkg/repository/dao"
	"github.com/abhidarbey/user-service/pkg/utils/datetime"
	"github.com/abhidarbey/user-service/pkg/utils/encrypt"
	"github.com/abhidarbey/user-service/pkg/utils/resterr"
)

// UserServiceInterface interface
type UserServiceInterface interface {
	Create(models.User, *models.AuthUser) (*models.User, *resterr.RestErr)
	FindAll(*models.AuthUser) (models.Users, *resterr.RestErr)
	GetByID(string, *models.AuthUser) (*models.User, *resterr.RestErr)
	Update(models.User, *models.AuthUser) (*models.User, *resterr.RestErr)
	UpdatePassword(models.User, *models.AuthUser) (*models.User, *resterr.RestErr)
	Delete(string, *models.AuthUser) *resterr.RestErr
}

type userService struct{}

// UserService variable
var (
	UserService UserServiceInterface = &userService{}
)

// Create user
func (s *userService) Create(user models.User, au *models.AuthUser) (*models.User, *resterr.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	// Verify permission --> IsGranted
	if !au.IsOrgAdmin {
		if !helpers.IsGranted("CanCreateUser", *au) {
			return nil, resterr.NewUnauthorizedError("Permission not granted")
		}
	}

	// Verify unique email
	_, emailErr := dao.UserDao.GetByEmail(user.Email)
	if emailErr == nil {
		return nil, resterr.NewBadRequestError("Email already registered")
	}

	user.ID = "U" + encrypt.GenerateID(20)
	user.Organization = au.Organization
	user.Password = encrypt.GetMd5(user.Password)
	user.Status = models.StatusActive
	user.IsActive = true
	user.IsOrgAdmin = false
	user.CreatedAt = datetime.GetDateTimeString()
	user.UpdatedAt = datetime.GetDateTimeString()

	// Add user roles
	if len(user.Roles) > 0 {
		// validate roles
		userRoleList, roleDeptList, rolePermList, err := helpers.AssignUserRoles(user)
		if err != nil {
			return nil, err
		}
		user.Roles = *userRoleList
		user.Departments = *roleDeptList
		user.Permissions = *rolePermList
		// user.Permissions = append(user.Permissions, *rolePermList...)
	}

	// Add user departments
	if len(user.Departments) > 0 {
		userDeptList, err := helpers.AssignUserDepartments(user)
		if err != nil {
			return nil, err
		}
		user.Departments = *userDeptList
	}

	// Add user permissions
	if len(user.Permissions) > 0 {
		userPermList, err := helpers.AssignUserPermissions(user)
		if err != nil {
			return nil, err
		}
		user.Permissions = *userPermList
	}

	// Create new user
	newUser, err := dao.UserDao.Create(user)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

// FindAll active users
func (s *userService) FindAll(au *models.AuthUser) (models.Users, *resterr.RestErr) {
	// Verify permission --> IsGranted
	if !au.IsOrgAdmin {
		if !helpers.IsGranted("CanReadUser", *au) {
			return nil, resterr.NewUnauthorizedError("Permission not granted")
		}
	}

	org := au.Organization
	users, err := dao.UserDao.FindAll(org)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *userService) GetByID(id string, au *models.AuthUser) (*models.User, *resterr.RestErr) {
	// Verify permission --> IsGranted
	if !au.IsOrgAdmin {
		if !helpers.IsGranted("CanReadUser", *au) {
			return nil, resterr.NewUnauthorizedError("Permission not granted")
		}
	}

	user, err := dao.UserDao.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Verify organization
	if !au.IsSuperuser && user.Organization != au.Organization {
		return nil, resterr.NewUnauthorizedError("Unauthorized request")
	}

	return user, nil
}

func (s *userService) Update(user models.User, au *models.AuthUser) (*models.User, *resterr.RestErr) {
	// Verify permission --> IsGranted
	if !au.IsOrgAdmin {
		if !helpers.IsGranted("CanUpdateUser", *au) {
			return nil, resterr.NewUnauthorizedError("Permission not granted")
		}
	}

	current, err := s.GetByID(user.ID, au)
	if err != nil {
		return nil, err
	}

	// set request user organization
	user.Organization = current.Organization

	// Verify user fields
	if user.Firstname != "" {
		current.Firstname = user.Firstname
	}

	if user.Lastname != "" {
		current.Lastname = user.Lastname
	}

	if user.Email != "" && user.Email != current.Email {
		// verify unique email
		_, emailErr := dao.UserDao.GetByEmail(user.Email)
		if emailErr == nil {
			return nil, resterr.NewBadRequestError("Email already registered")
		}
		current.Email = user.Email
	}

	if len(user.Roles) > 0 {
		// validate roles
		userRoleList, roleDeptList, rolePermList, err := helpers.AssignUserRoles(user)
		if err != nil {
			return nil, err
		}
		user.Roles = *userRoleList
		user.Departments = *roleDeptList
		user.Permissions = *rolePermList

		// set current user roles
		current.Roles = user.Roles
	}

	// Add user departments
	if len(user.Departments) > 0 {
		userDeptList, err := helpers.AssignUserDepartments(user)
		if err != nil {
			return nil, err
		}
		current.Departments = *userDeptList
	}

	if len(user.Permissions) > 0 {
		// Validate permission request
		userPermList, err := helpers.AssignUserPermissions(user)
		if err != nil {
			return nil, err
		}

		current.Permissions = *userPermList
	}

	current.UpdatedAt = datetime.GetDateTimeString()

	// Update user
	if updateErr := dao.UserDao.Update(*current); updateErr != nil {
		return nil, updateErr
	}

	return current, nil
}

func (s *userService) UpdatePassword(user models.User, au *models.AuthUser) (*models.User, *resterr.RestErr) {
	// Get current user from db
	current, err := s.GetByID(user.ID, au)
	if err != nil {
		return nil, err
	}

	// Verify permission --> IsGranted
	if !au.IsOrgAdmin && au.ID != current.ID {
		return nil, resterr.NewUnauthorizedError("Unauthorized request")
	}

	// Verify password field
	if user.Password != "" {
		user.Password = encrypt.GetMd5(user.Password)
		current.Password = user.Password
	}

	current.UpdatedAt = datetime.GetDateTimeString()

	// Update user
	if updateErr := dao.UserDao.Update(*current); updateErr != nil {
		return nil, updateErr
	}

	return current, nil
}

func (s *userService) Delete(id string, au *models.AuthUser) *resterr.RestErr {
	// Verify permission --> IsGranted
	if !au.IsOrgAdmin {
		if !helpers.IsGranted("CanDeleteUser", *au) {
			return resterr.NewUnauthorizedError("Permission not granted")
		}
	}

	return dao.UserDao.Delete(id)
}
