package models

import (
	"encoding/json"
	"strings"

	"github.com/abhidarbey/user-service/pkg/utils/resterr"
)

// User Structure (Model)
type User struct {
	ID           string           `json:"id" bson:"id"`
	Firstname    string           `json:"first_name" bson:"first_name"`
	Lastname     string           `json:"last_name" bson:"last_name"`
	Email        string           `json:"email" bson:"email"`
	Password     string           `json:"password" bson:"password"`
	Organization string           `json:"organization" bson:"organization"`
	Departments  []UserDepartment `json:"departments" bson:"departments"`
	Roles        []UserRole       `json:"roles" bson:"roles"`
	Permissions  []Permission     `json:"permissions" bson:"permissions"`
	Status       string           `json:"status" bson:"status"`
	IsActive     bool             `json:"is_active" bson:"is_active"`
	IsSuperuser  bool             `json:"is_superuser" bson:"is_superuser"`
	IsOrgAdmin   bool             `json:"is_org_admin" bson:"is_org_admin"`
	CreatedAt    string           `json:"created_at" bson:"created_at"`
	UpdatedAt    string           `json:"updated_at" bson:"updated_at"`
}

// Users array
type Users []User

// UserPermissions structure (db)
type UserPermissions struct {
	UserID      string       `json:"user_id" bson:"user_id"`
	Permissions []Permission `json:"permissions" bson:"permissions"`
}

// PrivateUser Structure
type PrivateUser struct {
	ID           string           `json:"id"`
	Firstname    string           `json:"first_name"`
	Lastname     string           `json:"last_name"`
	Email        string           `json:"email"`
	Status       string           `json:"status"`
	Organization string           `json:"organization"`
	Departments  []UserDepartment `json:"departments"`
	Roles        []UserRole       `json:"roles"`
	Permissions  []Permission     `json:"permissions"`
	IsActive     bool             `json:"is_active"`
	IsSuperuser  bool             `json:"is_superuser"`
	IsOrgAdmin   bool             `json:"is_org_admin"`
	CreatedAt    string           `json:"created_at"`
	UpdatedAt    string           `json:"updated_at"`
}

// AuthUser Structure
type AuthUser struct {
	ID           string       `json:"id"`
	Organization string       `json:"organization"`
	IsSuperuser  bool         `json:"is_superuser"`
	IsOrgAdmin   bool         `json:"is_org_admin"`
	Permissions  []Permission `json:"permissions"`
}

// Marshal User interface
func (user *User) Marshal() interface{} {
	userJSON, _ := json.Marshal(user)
	var privateUser PrivateUser
	json.Unmarshal(userJSON, &privateUser)
	return privateUser
}

// Marshal Users array interface
func (users Users) Marshal() []interface{} {
	result := make([]interface{}, len(users))
	for index, user := range users {
		result[index] = user.Marshal()
	}
	return result
}

// AuthMarshal User interface
func (user *User) AuthMarshal() interface{} {
	userJSON, _ := json.Marshal(user)
	var auth AuthUser
	json.Unmarshal(userJSON, &auth)
	return auth
}

// Validate function
func (user *User) Validate() *resterr.RestErr {
	if user.Firstname == "" {
		return resterr.NewBadRequestError("Firstname is required")
	}
	if user.Lastname == "" {
		return resterr.NewBadRequestError("Lastname is required")
	}
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	if user.Email == "" {
		return resterr.NewBadRequestError("Invalid Email address")
	}
	user.Password = strings.TrimSpace(user.Password)
	if user.Password == "" {
		return resterr.NewBadRequestError("Password field cannot be empty")
	}
	return nil
}
