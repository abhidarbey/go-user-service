package models

import (
	"github.com/abhidarbey/user-service/pkg/utils/resterr"
)

// Department Structure (Model)
type Department struct {
	ID           string `json:"id" bson:"id"`
	Organization string `json:"organization" bson:"organization"`
	Name         string `json:"name" bson:"name"`
	Status       string `json:"status" bson:"status"`
	IsActive     bool   `json:"is_active" bson:"is_active"`
	CreatedAt    string `json:"created_at" bson:"created_at"`
	UpdatedAt    string `json:"updated_at" bson:"updated_at"`
}

// Departments array
type Departments []Department

// UserDepartment Structure
type UserDepartment struct {
	DepartmentID   string `json:"department_id" bson:"department_id"`
	DepartmentName string `json:"department_name" bson:"department_name"`
}

// UserDepartments array
type UserDepartments []UserDepartment

// Validate function
func (department *Department) Validate() *resterr.RestErr {
	// if department.Organization == "" {
	// 	return resterr.NewBadRequestError("Organization ID is required")
	// }
	if department.Name == "" {
		return resterr.NewBadRequestError("Department name is required")
	}
	return nil
}
