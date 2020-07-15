package helpers

import (
	"github.com/abhidarbey/user-service/pkg/models"
	"github.com/abhidarbey/user-service/pkg/repository/dao"
	"github.com/abhidarbey/user-service/pkg/utils/resterr"
)

// AssignUserDepartments to the user
func AssignUserDepartments(user models.User) (*[]models.UserDepartment, *resterr.RestErr) {
	// Get departments list from db
	deptList, err := dao.DepartmentDao.FindAll(user.Organization)
	if err != nil {
		return nil, err
	}

	// Factor out invalid departments from user.Departments
	validList := ValidateDepartments(user, deptList)

	// factor out duplicate entries
	userDeptList := UniqueDepartments(validList)

	return &userDeptList, nil
}

// AssignUserRoles to the user
func AssignUserRoles(user models.User) (*[]models.UserRole, *[]models.UserDepartment, *[]models.Permission, *resterr.RestErr) {
	// Get departments list from db
	roleList, err := dao.RoleDao.FindAll(user.Organization)
	if err != nil {
		return nil, nil, nil, err
	}

	// Factor out invalid departments from user.Roles
	validList := ValidateRoles(user, roleList)

	// factor out duplicate entries
	userRoleList := UniqueRoles(validList)

	// set role departments to user departments
	roleDeptList := AssignRoleDeptToUser(user, roleList)

	// assign roles permissions to user
	rolePermList, err := AssignRolesPermToUser(user)
	if err != nil {
		return nil, nil, nil, err
	}

	return &userRoleList, &roleDeptList, rolePermList, nil
}

// AssignUserPermissions to the user
func AssignUserPermissions(user models.User) (*[]models.Permission, *resterr.RestErr) {
	// Get permissions list from db
	permList, err := dao.PermissionDao.FindAll()
	if err != nil {
		return nil, err
	}
	// Factor out invalid permissions from user.Permissions
	validList := ValidatePermissions(user.Permissions, permList)

	// factor out duplicate entries
	userPermList := UniquePermissions(validList)

	return &userPermList, nil
}

// UniqueDepartments verifies the validity of request
func UniqueDepartments(list []models.UserDepartment) []models.UserDepartment {
	encountered := make(map[models.UserDepartment]bool)
	result := []models.UserDepartment{}

	for v := range list {
		encountered[list[v]] = true
	}

	for k := range encountered {
		result = append(result, k)
	}

	return result
}

// UniqueRoles verifies the validity of request
func UniqueRoles(list []models.UserRole) []models.UserRole {
	encountered := make(map[models.UserRole]bool)
	result := []models.UserRole{}

	for v := range list {
		encountered[list[v]] = true
	}

	for k := range encountered {
		result = append(result, k)
	}

	return result
}

// ValidateDepartments verifies the validity of request
func ValidateDepartments(user models.User, deptList []models.Department) []models.UserDepartment {
	validList := []models.UserDepartment{}

	// Factor out invalid departments from user.Departments
	for i := 0; i < len(user.Departments); i++ {
		for j := 0; j < len(deptList); j++ {
			if user.Departments[i].DepartmentID == deptList[j].ID {
				user.Departments[i].DepartmentName = deptList[j].Name
				validList = append(validList, user.Departments[i])
			}
		}
	}

	return validList
}

// ValidateRoles verifies the validity of request
func ValidateRoles(user models.User, roleList []models.Role) []models.UserRole {
	validList := []models.UserRole{}

	if len(user.Roles) > 0 {
		// Factor out invalid departments from user.Roles
		for i := 0; i < len(user.Roles); i++ {
			for j := 0; j < len(roleList); j++ {
				if user.Roles[i].RoleID == roleList[j].ID {
					user.Roles[i].RoleName = roleList[j].Name
					validList = append(validList, user.Roles[i])
				}
			}
		}
	}

	return validList
}

// AssignRolesPermToUser sets roles permissions as user permissions
func AssignRolesPermToUser(user models.User) (*[]models.Permission, *resterr.RestErr) {
	// Get all role permissions
	rp, err := dao.RoleDao.FindAllRolePermissions(user.Organization)
	if err != nil {
		return nil, err
	}

	rolePermList := []models.Permission{}

	for i := 0; i < len(user.Roles); i++ {
		for j := 0; j < len(rp); j++ {
			if user.Roles[i].RoleID == rp[j].RoleID {
				rolePermList = append(rolePermList, rp[j].Permissions...)
			}
		}

	}

	userPermList := []models.Permission{}
	userPermList = append(userPermList, rolePermList...)

	return &userPermList, nil
}

// AssignRoleDeptToUser from roles
func AssignRoleDeptToUser(user models.User, roleList []models.Role) []models.UserDepartment {
	userDeptList := []models.UserDepartment{}
	newDept := models.UserDepartment{}

	for i := 0; i < len(user.Roles); i++ {
		for j := 0; j < len(roleList); j++ {
			if user.Roles[i].RoleID == roleList[j].ID {
				newDept.DepartmentID = roleList[j].Department
				userDeptList = append(userDeptList, newDept)
			}
		}
	}

	return userDeptList
}
