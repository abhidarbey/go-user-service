package dao

import (
	"context"
	"time"

	"github.com/abhidarbey/user-service/pkg/models"
	"github.com/abhidarbey/user-service/pkg/settings/db/mongodb"
	"github.com/abhidarbey/user-service/pkg/utils/resterr"
	"go.mongodb.org/mongo-driver/bson"
)

// RoleDaoInterface type
type RoleDaoInterface interface {
	Create(models.Role) (*models.Role, *resterr.RestErr)
	FindAll(string) (models.Roles, *resterr.RestErr)
	GetByID(string, string) (*models.Role, *resterr.RestErr)
	Update(models.Role) *resterr.RestErr
	Delete(string) *resterr.RestErr
	FindAllRolePermissions(string) ([]models.RolePermissions, *resterr.RestErr)
}

type roleDao struct{}

// RoleDao variable
var (
	RoleDao RoleDaoInterface = &roleDao{}
)

// Create role
func (d *roleDao) Create(role models.Role) (*models.Role, *resterr.RestErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userDB := mongodb.Client.Database("erp-user-service")

	roleCollection := userDB.Collection("role")

	_, err := roleCollection.InsertOne(ctx, bson.M{
		"id":           role.ID,
		"organization": role.Organization,
		"department":   role.Department,
		"name":         role.Name,
		"status":       role.Status,
		"is_active":    role.IsActive,
		"created_at":   role.CreatedAt,
		"updated_at":   role.UpdatedAt,
	})
	if err != nil {
		return nil, resterr.NewInternalServerError(err.Error())
	}

	// Add role permissions
	if rpErr := d.addPermissions(role); rpErr != nil {
		return nil, rpErr
	}

	return &role, nil
}

// FindAll role
func (d *roleDao) FindAll(org string) (models.Roles, *resterr.RestErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userDB := mongodb.Client.Database("erp-user-service")

	roles := []models.Role{}
	roleCollection := userDB.Collection("role")

	filter := bson.M{"status": models.StatusActive, "organization": org}
	cursor, err := roleCollection.Find(ctx, filter)
	if err != nil {
		return nil, resterr.NewInternalServerError(err.Error())
	}

	if err = cursor.All(ctx, &roles); err != nil {
		return nil, resterr.NewInternalServerError(err.Error())
	}

	return roles, nil
}

// GetByID role
func (d *roleDao) GetByID(id string, org string) (*models.Role, *resterr.RestErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userDB := mongodb.Client.Database("erp-user-service")

	role := models.Role{}
	roleCollection := userDB.Collection("role")

	filter := bson.M{"id": id, "organization": org}
	err := roleCollection.FindOne(ctx, filter).Decode(&role)
	if err != nil {
		return nil, resterr.NewInternalServerError(err.Error())
	}

	// Get role permisions
	rp, rpErr := d.getPermissionsByRoleID(id)
	if rpErr != nil {
		return nil, rpErr
	}

	role.Permissions = rp.Permissions

	return &role, nil
}

// Update  role
func (d *roleDao) Update(role models.Role) *resterr.RestErr {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userDB := mongodb.Client.Database("erp-user-service")
	roleCollection := userDB.Collection("role")

	filter := bson.M{"id": role.ID}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "name", Value: role.Name},
			{Key: "status", Value: role.Status},
			{Key: "is_active", Value: role.IsActive},
			{Key: "updated_at", Value: role.UpdatedAt},
		}},
	}

	_, err := roleCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return resterr.NewInternalServerError(err.Error())
	}

	// Update role permissions
	if err := d.updatePermissions(role); err != nil {
		return err
	}
	return nil
}

// Delete role
func (d *roleDao) Delete(id string) *resterr.RestErr {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userDB := mongodb.Client.Database("erp-user-service")
	roleCollection := userDB.Collection("role")

	filter := bson.M{"id": id}

	_, err := roleCollection.DeleteOne(ctx, filter)
	if err != nil {
		return resterr.NewInternalServerError(err.Error())
	}

	// Delete permissions assoicated with this role
	if rpErr := d.deletePermissions(id); rpErr != nil {
		return rpErr
	}

	return nil
}

// FindAllRolePermissions role
func (d *roleDao) FindAllRolePermissions(org string) ([]models.RolePermissions, *resterr.RestErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userDB := mongodb.Client.Database("erp-user-service")

	rp := []models.RolePermissions{}
	roleCollection := userDB.Collection("role-permissions")

	filter := bson.M{"organization": org}
	cursor, err := roleCollection.Find(ctx, filter)
	if err != nil {
		return nil, resterr.NewInternalServerError(err.Error())
	}

	if err = cursor.All(ctx, &rp); err != nil {
		return nil, resterr.NewInternalServerError(err.Error())
	}

	return rp, nil
}

// addPermissions to role
func (d *roleDao) addPermissions(role models.Role) *resterr.RestErr {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userDB := mongodb.Client.Database("erp-user-service")

	roleCollection := userDB.Collection("role-permissions")

	_, err := roleCollection.InsertOne(ctx, bson.M{
		"role_id":      role.ID,
		"organization": role.Organization,
		"permissions":  role.Permissions,
	})
	if err != nil {
		return resterr.NewInternalServerError(err.Error())
	}
	return nil
}

// getPermissionsByRoleID
func (d *roleDao) getPermissionsByRoleID(roleID string) (*models.RolePermissions, *resterr.RestErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userDB := mongodb.Client.Database("erp-user-service")

	rp := models.RolePermissions{}
	roleCollection := userDB.Collection("role-permissions")

	filter := bson.M{"role_id": roleID}
	err := roleCollection.FindOne(ctx, filter).Decode(&rp)
	if err != nil {
		return nil, resterr.NewInternalServerError(err.Error())
	}

	return &rp, nil
}

// updatePermissions  role
func (d *roleDao) updatePermissions(role models.Role) *resterr.RestErr {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userDB := mongodb.Client.Database("erp-user-service")
	roleCollection := userDB.Collection("role-permissions")

	filter := bson.M{"role_id": role.ID}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "permissions", Value: role.Permissions},
		}},
	}

	_, err := roleCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return resterr.NewInternalServerError(err.Error())
	}
	return nil
}

// deletePermissions role permissions
func (d *roleDao) deletePermissions(roleID string) *resterr.RestErr {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userDB := mongodb.Client.Database("erp-user-service")
	roleCollection := userDB.Collection("role-permissions")

	filter := bson.M{"role_id": roleID}

	_, err := roleCollection.DeleteOne(ctx, filter)
	if err != nil {
		return resterr.NewInternalServerError(err.Error())
	}
	return nil
}
