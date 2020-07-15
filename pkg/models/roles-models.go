package models

import (
	"github.com/abhidarbey/user-service/pkg/utils/resterr"
)

// Role Structure (Model)
type Role struct {
	ID           string       `json:"id" bson:"id"`
	Organization string       `json:"organization" bson:"organization"`
	Department   string       `json:"department" bson:"department"`
	Name         string       `json:"name" bson:"name"`
	Permissions  []Permission `json:"permissions" bson:"permissions"`
	Status       string       `json:"status" bson:"status"`
	IsActive     bool         `json:"is_active" bson:"is_active"`
	CreatedAt    string       `json:"created_at" bson:"created_at"`
	UpdatedAt    string       `json:"updated_at" bson:"updated_at"`
}

// Roles array
type Roles []Role

// RolePermissions structure (db)
type RolePermissions struct {
	RoleID       string       `json:"role_id" bson:"role_id"`
	Organization string       `json:"organization" bson:"organization"`
	Permissions  []Permission `json:"permissions" bson:"permissions"`
}

// UserRole Structure
type UserRole struct {
	RoleID   string `json:"role_id" bson:"role_id"`
	RoleName string `json:"role_name" bson:"role_name"`
}

// UserRoles array
type UserRoles []UserRole

// Validate function
func (role *Role) Validate() *resterr.RestErr {
	if role.Name == "" {
		return resterr.NewBadRequestError("Role Name is required")
	}
	if role.Department == "" {
		return resterr.NewBadRequestError("Department ID is required")
	}
	return nil
}
