package services

import (
	"strings"

	"github.com/abhidarbey/user-service/pkg/models"
	"github.com/abhidarbey/user-service/pkg/repository/dao"
	"github.com/abhidarbey/user-service/pkg/utils/datetime"
	"github.com/abhidarbey/user-service/pkg/utils/encrypt"
	"github.com/abhidarbey/user-service/pkg/utils/resterr"
)

// AuthServiceInterface interface
type AuthServiceInterface interface {
	Login(models.LoginRequest) (*models.User, *resterr.RestErr)
	RegisterOrg(models.RegistrationRequest) (*models.Organization, *resterr.RestErr)
	RegisterSuperuser(models.User) (*models.User, *resterr.RestErr)
}

type authService struct{}

// AuthService variable
var (
	AuthService AuthServiceInterface = &authService{}
)

// Login service
func (s *authService) Login(request models.LoginRequest) (*models.User, *resterr.RestErr) {
	email := request.Email
	password := encrypt.GetMd5(request.Password)
	user, err := dao.AuthDao.Login(email, password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// RegisterOrg organization
func (s *authService) RegisterOrg(request models.RegistrationRequest) (*models.Organization, *resterr.RestErr) {
	// Validate request
	if err := request.Validate(); err != nil {
		return nil, err
	}

	// Set organization fields
	org := models.Organization{}
	org.ID = strings.TrimSpace(strings.ToUpper(request.OrgName)) + encrypt.GenerateID(10)
	org.Name = request.OrgName
	org.Website = request.Website
	org.Status = models.StatusActive
	org.IsActive = true
	org.CreatedAt = datetime.GetDateTimeString()
	org.UpdatedAt = datetime.GetDateTimeString()

	// Verify unique email
	_, emailErr := dao.UserDao.GetByEmail(request.Email)
	if emailErr == nil {
		return nil, resterr.NewBadRequestError("Email already registered")
	}

	// Set user fields
	user := models.User{}
	user.ID = "U" + encrypt.GenerateID(20)
	user.Firstname = request.Firstname
	user.Lastname = request.Lastname
	user.Email = request.Email
	user.Password = encrypt.GetMd5(request.Password)
	user.Organization = org.ID
	user.Status = models.StatusActive
	user.IsActive = true
	user.IsOrgAdmin = true
	user.IsSuperuser = false
	user.CreatedAt = datetime.GetDateTimeString()
	user.UpdatedAt = datetime.GetDateTimeString()

	// Create organization
	newOrganization, err := dao.OrganizationDao.Create(org)
	if err != nil {
		return nil, err
	}

	// Create user
	_, err = dao.UserDao.Create(user)
	if err != nil {
		return nil, err
	}

	return newOrganization, nil
}

// RegisterSuperuser func
func (s *authService) RegisterSuperuser(user models.User) (*models.User, *resterr.RestErr) {
	// Verify unique email
	_, emailErr := dao.UserDao.GetByEmail(user.Email)
	if emailErr == nil {
		return nil, resterr.NewBadRequestError("Email already registered")
	}

	// Set user fields
	user.ID = "U" + encrypt.GenerateID(20)
	user.Password = encrypt.GetMd5(user.Password)
	user.Organization = ""
	user.Status = models.StatusActive
	user.IsActive = true
	user.IsOrgAdmin = false
	user.IsSuperuser = true
	user.CreatedAt = datetime.GetDateTimeString()
	user.UpdatedAt = datetime.GetDateTimeString()

	superuser, err := dao.UserDao.Create(user)
	if err != nil {
		return nil, err
	}
	return superuser, nil
}
