package models

import (
	"github.com/abhidarbey/user-service/pkg/utils/resterr"
)

// Organization Structure (Model)
type Organization struct {
	ID        string `json:"id" bson:"id"`
	Name      string `json:"name" bson:"name"`
	Website   string `json:"website" bson:"website"`
	Status    string `json:"status" bson:"status"`
	IsActive  bool   `json:"is_active" bson:"is_active"`
	CreatedAt string `json:"created_at" bson:"created_at"`
	UpdatedAt string `json:"updated_at" bson:"updated_at"`
}

// Organizations array
type Organizations []Organization

// Validate function
func (org *Organization) Validate() *resterr.RestErr {
	if org.Name == "" {
		return resterr.NewBadRequestError("Organization name is required")
	}
	return nil
}
