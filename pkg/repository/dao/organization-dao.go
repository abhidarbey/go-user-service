package dao

import (
	"context"
	"time"

	"github.com/abhidarbey/user-service/pkg/models"
	"github.com/abhidarbey/user-service/pkg/settings/db/mongodb"
	"github.com/abhidarbey/user-service/pkg/utils/resterr"
	"go.mongodb.org/mongo-driver/bson"
)

// OrganizationDaoInterface type
type OrganizationDaoInterface interface {
	Create(models.Organization) (*models.Organization, *resterr.RestErr)
	FindAll() (models.Organizations, *resterr.RestErr)
	GetByID(string) (*models.Organization, *resterr.RestErr)
	Update(models.Organization) *resterr.RestErr
	Delete(string) *resterr.RestErr
}

type organizationDao struct{}

// OrganizationDao variable
var (
	OrganizationDao OrganizationDaoInterface = &organizationDao{}
)

// Create organization
func (d *organizationDao) Create(organization models.Organization) (*models.Organization, *resterr.RestErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userDB := mongodb.Client.Database("erp-user-service")

	orgCollection := userDB.Collection("organization")

	_, err := orgCollection.InsertOne(ctx, bson.M{
		"id":         organization.ID,
		"name":       organization.Name,
		"website":    organization.Website,
		"status":     organization.Status,
		"is_active":  organization.IsActive,
		"created_at": organization.CreatedAt,
		"updated_at": organization.UpdatedAt,
	})
	if err != nil {
		return nil, resterr.NewInternalServerError(err.Error())
	}
	return &organization, nil
}

// FindAll organization
func (d *organizationDao) FindAll() (models.Organizations, *resterr.RestErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userDB := mongodb.Client.Database("erp-user-service")

	organizations := []models.Organization{}
	orgCollection := userDB.Collection("organization")

	filter := bson.M{"status": models.StatusActive}
	cursor, err := orgCollection.Find(ctx, filter)
	if err != nil {
		return nil, resterr.NewInternalServerError(err.Error())
	}

	if err = cursor.All(ctx, &organizations); err != nil {
		return nil, resterr.NewInternalServerError(err.Error())
	}

	return organizations, nil
}

// GetByID organization
func (d *organizationDao) GetByID(id string) (*models.Organization, *resterr.RestErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userDB := mongodb.Client.Database("erp-user-service")

	organization := models.Organization{}
	orgCollection := userDB.Collection("organization")

	filter := bson.M{"id": id}
	err := orgCollection.FindOne(ctx, filter).Decode(&organization)
	if err != nil {
		return nil, resterr.NewInternalServerError(err.Error())
	}

	return &organization, nil
}

// Update  organization
func (d *organizationDao) Update(organization models.Organization) *resterr.RestErr {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userDB := mongodb.Client.Database("erp-user-service")
	orgCollection := userDB.Collection("organization")

	filter := bson.M{"id": organization.ID}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "name", Value: organization.Name},
			{Key: "website", Value: organization.Website},
			{Key: "status", Value: organization.Status},
			{Key: "is_active", Value: organization.IsActive},
			{Key: "updated_at", Value: organization.UpdatedAt},
		}},
	}

	_, err := orgCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return resterr.NewInternalServerError(err.Error())
	}
	return nil
}

// Delete User
func (d *organizationDao) Delete(id string) *resterr.RestErr {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userDB := mongodb.Client.Database("erp-user-service")
	orgCollection := userDB.Collection("organization")

	filter := bson.M{"id": id}

	_, err := orgCollection.DeleteOne(ctx, filter)
	if err != nil {
		return resterr.NewInternalServerError(err.Error())
	}
	return nil
}
