package models

import (
	"github.com/abhidarbey/user-service/pkg/utils/resterr"
)

// LoginRequest structure
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegistrationRequest Structure
type RegistrationRequest struct {
	OrgName   string `json:"org_name"`
	Website   string `json:"website"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// ValueToken struct
type ValueToken struct {
	ValueToken string `json:"value_token"`
}

// TokenHeader struct
type TokenHeader struct {
	TYP string `json:"typ"`
	ALG string `json:"alg"`
}

// TokenPayload struct
type TokenPayload struct {
	ID           string       `json:"id"`
	Organization string       `json:"organization"`
	IsSuperuser  bool         `json:"is_superuser"`
	IsOrgAdmin   bool         `json:"is_org_admin"`
	Permissions  []Permission `json:"permissions"`
	Authorized   bool         `json:"authorized"`
	Expiry       int64        `json:"exp"`
}

// Validate LoginRequest
func (r *LoginRequest) Validate() *resterr.RestErr {
	if r.Email == "" {
		return resterr.NewBadRequestError("Email is required")
	}
	if r.Password == "" {
		return resterr.NewBadRequestError("Password is required")
	}
	return nil
}

// Validate RegistrationRequest
func (r *RegistrationRequest) Validate() *resterr.RestErr {
	if r.OrgName == "" {
		return resterr.NewBadRequestError("Organization name is required")
	}
	if r.Firstname == "" {
		return resterr.NewBadRequestError("Firstname is required")
	}
	if r.Lastname == "" {
		return resterr.NewBadRequestError("Lastname is required")
	}
	if r.Email == "" {
		return resterr.NewBadRequestError("Email is required")
	}
	if r.Password == "" {
		return resterr.NewBadRequestError("Password is required")
	}
	return nil
}
