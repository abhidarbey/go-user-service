package models

// Permission struct
type Permission struct {
	Name string `json:"name" bson:"name"`
}

// Permissions list
type Permissions []Permission
