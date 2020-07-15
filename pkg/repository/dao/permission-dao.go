package dao

import (
	"context"
	"time"

	"github.com/abhidarbey/user-service/pkg/models"
	"github.com/abhidarbey/user-service/pkg/settings/db/mongodb"
	"github.com/abhidarbey/user-service/pkg/utils/resterr"
	"go.mongodb.org/mongo-driver/bson"
)

// PermissionDaoInterface type
type PermissionDaoInterface interface {
	Create(models.Permission) (*models.Permission, *resterr.RestErr)
	FindAll() (models.Permissions, *resterr.RestErr)
	GetByName(string) (*models.Permission, *resterr.RestErr)
	Update(models.Permission) *resterr.RestErr
	Delete(string) *resterr.RestErr
}

type permissionDao struct{}

// PermissionDao variable
var (
	PermissionDao PermissionDaoInterface = &permissionDao{}
)

// Create permission
func (d *permissionDao) Create(permission models.Permission) (*models.Permission, *resterr.RestErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userDB := mongodb.Client.Database("erp-user-service")

	permissionCollection := userDB.Collection("permission")

	_, err := permissionCollection.InsertOne(ctx, bson.M{
		"name": permission.Name,
	})
	if err != nil {
		return nil, resterr.NewInternalServerError(err.Error())
	}
	return &permission, nil
}

// FindAll permission
func (d *permissionDao) FindAll() (models.Permissions, *resterr.RestErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userDB := mongodb.Client.Database("erp-user-service")

	permissions := []models.Permission{}
	permissionCollection := userDB.Collection("permission")

	filter := bson.M{}
	cursor, err := permissionCollection.Find(ctx, filter)
	if err != nil {
		return nil, resterr.NewInternalServerError(err.Error())
	}

	if err = cursor.All(ctx, &permissions); err != nil {
		return nil, resterr.NewInternalServerError(err.Error())
	}

	return permissions, nil
}

// GetByName permission
func (d *permissionDao) GetByName(name string) (*models.Permission, *resterr.RestErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userDB := mongodb.Client.Database("erp-user-service")

	permission := models.Permission{}
	permissionCollection := userDB.Collection("permission")

	filter := bson.M{"name": name}
	err := permissionCollection.FindOne(ctx, filter).Decode(&permission)
	if err != nil {
		return nil, resterr.NewInternalServerError(err.Error())
	}

	return &permission, nil
}

// Update  permission
func (d *permissionDao) Update(permission models.Permission) *resterr.RestErr {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userDB := mongodb.Client.Database("erp-user-service")
	permissionCollection := userDB.Collection("permission")

	filter := bson.M{"name": permission.Name}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "name", Value: permission.Name},
		}},
	}

	_, err := permissionCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return resterr.NewInternalServerError(err.Error())
	}
	return nil
}

// Delete User
func (d *permissionDao) Delete(name string) *resterr.RestErr {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userDB := mongodb.Client.Database("erp-user-service")
	permissionCollection := userDB.Collection("permission")

	filter := bson.M{"name": name}

	_, err := permissionCollection.DeleteOne(ctx, filter)
	if err != nil {
		return resterr.NewInternalServerError(err.Error())
	}
	return nil
}
