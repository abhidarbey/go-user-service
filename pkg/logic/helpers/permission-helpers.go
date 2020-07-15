package helpers

import (
	"github.com/abhidarbey/user-service/pkg/models"
	"github.com/abhidarbey/user-service/pkg/repository/dao"
	"github.com/abhidarbey/user-service/pkg/utils/resterr"
)

// AssignRolePermissions to the user
func AssignRolePermissions(role models.Role) (*[]models.Permission, *resterr.RestErr) {
	// Get permissions list from db
	permList, err := dao.PermissionDao.FindAll()
	if err != nil {
		return nil, err
	}
	// Factor out invalid permissions from role.Permissions
	validList := ValidatePermissions(role.Permissions, permList)

	// factor out duplicate entries
	rolePermList := UniquePermissions(validList)

	return &rolePermList, nil
}

// ValidatePermissions verifies the validity of request
func ValidatePermissions(request []models.Permission, permList []models.Permission) []models.Permission {
	validList := []models.Permission{}

	if len(request) > 0 {
		// Factor out invalid permissions from request
		for i := 0; i < len(request); i++ {
			for j := 0; j < len(permList); j++ {
				if request[i].Name == permList[j].Name {
					validList = append(validList, request[i])
				}
			}
		}
	}

	return validList
}

// UniquePermissions verifies the validity of request
func UniquePermissions(list []models.Permission) []models.Permission {
	encountered := make(map[models.Permission]bool)
	result := []models.Permission{}

	for v := range list {
		encountered[list[v]] = true
	}

	for k := range encountered {
		result = append(result, k)
	}

	return result
}
