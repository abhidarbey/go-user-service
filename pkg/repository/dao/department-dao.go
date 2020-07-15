package dao

import (
	"context"
	"time"

	"github.com/abhidarbey/user-service/pkg/models"
	"github.com/abhidarbey/user-service/pkg/settings/db/mongodb"
	"github.com/abhidarbey/user-service/pkg/utils/resterr"
	"go.mongodb.org/mongo-driver/bson"
)

// DepartmentDaoInterface type
type DepartmentDaoInterface interface {
	Create(models.Department) (*models.Department, *resterr.RestErr)
	FindAll(string) (models.Departments, *resterr.RestErr)
	GetByID(string, string) (*models.Department, *resterr.RestErr)
	Update(models.Department) *resterr.RestErr
	Delete(string) *resterr.RestErr
}

type departmentDao struct{}

// DepartmentDao variable
var (
	DepartmentDao DepartmentDaoInterface = &departmentDao{}
)

// Create department
func (d *departmentDao) Create(department models.Department) (*models.Department, *resterr.RestErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userDB := mongodb.Client.Database("erp-user-service")

	deptCollection := userDB.Collection("department")

	_, err := deptCollection.InsertOne(ctx, bson.M{
		"id":           department.ID,
		"organization": department.Organization,
		"name":         department.Name,
		"status":       department.Status,
		"is_active":    department.IsActive,
		"created_at":   department.CreatedAt,
		"updated_at":   department.UpdatedAt,
	})
	if err != nil {
		return nil, resterr.NewInternalServerError(err.Error())
	}
	return &department, nil
}

// FindAll department
func (d *departmentDao) FindAll(org string) (models.Departments, *resterr.RestErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userDB := mongodb.Client.Database("erp-user-service")

	departments := []models.Department{}
	deptCollection := userDB.Collection("department")

	filter := bson.M{"status": models.StatusActive, "organization": org}
	cursor, err := deptCollection.Find(ctx, filter)
	if err != nil {
		return nil, resterr.NewInternalServerError(err.Error())
	}

	if err = cursor.All(ctx, &departments); err != nil {
		return nil, resterr.NewInternalServerError(err.Error())
	}

	return departments, nil
}

// GetByID department
func (d *departmentDao) GetByID(id string, org string) (*models.Department, *resterr.RestErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userDB := mongodb.Client.Database("erp-user-service")

	department := models.Department{}
	deptCollection := userDB.Collection("department")

	filter := bson.M{"id": id, "organization": org}
	err := deptCollection.FindOne(ctx, filter).Decode(&department)
	if err != nil {
		return nil, resterr.NewInternalServerError(err.Error())
	}

	return &department, nil
}

// Update  department
func (d *departmentDao) Update(department models.Department) *resterr.RestErr {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userDB := mongodb.Client.Database("erp-user-service")
	deptCollection := userDB.Collection("department")

	filter := bson.M{"id": department.ID}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "name", Value: department.Name},
			{Key: "status", Value: department.Status},
			{Key: "is_active", Value: department.IsActive},
			{Key: "updated_at", Value: department.UpdatedAt},
		}},
	}

	_, err := deptCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return resterr.NewInternalServerError(err.Error())
	}
	return nil
}

// Delete User
func (d *departmentDao) Delete(id string) *resterr.RestErr {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userDB := mongodb.Client.Database("erp-user-service")
	deptCollection := userDB.Collection("department")

	filter := bson.M{"id": id}

	_, err := deptCollection.DeleteOne(ctx, filter)
	if err != nil {
		return resterr.NewInternalServerError(err.Error())
	}
	return nil
}
