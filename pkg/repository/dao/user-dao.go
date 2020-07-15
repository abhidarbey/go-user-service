package dao

import (
	"context"
	"time"

	"github.com/abhidarbey/user-service/pkg/models"
	"github.com/abhidarbey/user-service/pkg/settings/db/mongodb"
	"github.com/abhidarbey/user-service/pkg/utils/resterr"
	"go.mongodb.org/mongo-driver/bson"
)

// UserDaoInterface type
type UserDaoInterface interface {
	Create(models.User) (*models.User, *resterr.RestErr)
	FindAll(string) ([]models.User, *resterr.RestErr)
	GetByID(string) (*models.User, *resterr.RestErr)
	GetByEmail(string) (*models.User, *resterr.RestErr)
	Update(models.User) *resterr.RestErr
	Delete(string) *resterr.RestErr
}

type userDao struct{}

// UserDao variable
var (
	UserDao UserDaoInterface = &userDao{}
)

// Create User
func (d *userDao) Create(user models.User) (*models.User, *resterr.RestErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userDB := mongodb.Client.Database("erp-user-service")

	userCollection := userDB.Collection("user")

	_, err := userCollection.InsertOne(ctx, bson.M{
		"id":           user.ID,
		"first_name":   user.Firstname,
		"last_name":    user.Lastname,
		"email":        user.Email,
		"password":     user.Password,
		"organization": user.Organization,
		"departments":  user.Departments,
		"roles":        user.Roles,
		"status":       user.Status,
		"is_active":    user.IsActive,
		"is_superuser": user.IsSuperuser,
		"is_org_admin": user.IsOrgAdmin,
		"created_at":   user.CreatedAt,
		"updated_at":   user.UpdatedAt,
	})
	if err != nil {
		return nil, resterr.NewInternalServerError(err.Error())
	}

	// add permissions
	if upErr := d.addPremissions(user); upErr != nil {
		return nil, upErr
	}

	return &user, nil
}

// FindAll Users
func (d *userDao) FindAll(org string) ([]models.User, *resterr.RestErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userDB := mongodb.Client.Database("erp-user-service")

	users := []models.User{}
	userCollection := userDB.Collection("user")

	filter := bson.M{"status": models.StatusActive, "organization": org}
	cursor, err := userCollection.Find(ctx, filter)
	if err != nil {
		return nil, resterr.NewInternalServerError(err.Error())
	}

	if err = cursor.All(ctx, &users); err != nil {
		return nil, resterr.NewInternalServerError(err.Error())
	}

	return users, nil
}

// GetByID User
func (d *userDao) GetByID(id string) (*models.User, *resterr.RestErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userDB := mongodb.Client.Database("erp-user-service")

	user := models.User{}
	userCollection := userDB.Collection("user")

	filter := bson.M{"id": id}
	err := userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, resterr.NewInternalServerError(err.Error())
	}

	// get user permissions
	up, upERr := d.getPermissonsByUserID(id)
	if upERr != nil {
		return nil, upERr
	}

	user.Permissions = up.Permissions

	return &user, nil
}

// GetByEmail User
func (d *userDao) GetByEmail(email string) (*models.User, *resterr.RestErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userDB := mongodb.Client.Database("erp-user-service")

	user := models.User{}
	userCollection := userDB.Collection("user")

	filter := bson.M{"email": email}
	err := userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, resterr.NewInternalServerError(err.Error())
	}

	return &user, nil
}

// Update User
func (d *userDao) Update(user models.User) *resterr.RestErr {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userDB := mongodb.Client.Database("erp-user-service")
	userCollection := userDB.Collection("user")

	filter := bson.M{"id": user.ID}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "first_name", Value: user.Firstname},
			{Key: "last_name", Value: user.Lastname},
			{Key: "email", Value: user.Email},
			{Key: "password", Value: user.Password},
			{Key: "departments", Value: user.Departments},
			{Key: "roles", Value: user.Roles},
			{Key: "status", Value: user.Status},
			{Key: "is_active", Value: user.IsActive},
			{Key: "is_superuser", Value: user.IsSuperuser},
			{Key: "is_org_admin", Value: user.IsOrgAdmin},
			{Key: "updated_at", Value: user.UpdatedAt},
		}},
	}

	_, err := userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return resterr.NewInternalServerError(err.Error())
	}
	return nil
}

// Delete User
func (d *userDao) Delete(id string) *resterr.RestErr {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userDB := mongodb.Client.Database("erp-user-service")
	userCollection := userDB.Collection("user")

	filter := bson.M{"id": id}

	_, err := userCollection.DeleteOne(ctx, filter)
	if err != nil {
		return resterr.NewInternalServerError(err.Error())
	}
	return nil
}

// addPremissions to user
func (d *userDao) addPremissions(user models.User) *resterr.RestErr {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userDB := mongodb.Client.Database("erp-user-service")

	userCollection := userDB.Collection("user-permissions")

	_, err := userCollection.InsertOne(ctx, bson.M{
		"user_id":     user.ID,
		"permissions": user.Permissions,
	})
	if err != nil {
		return resterr.NewInternalServerError(err.Error())
	}
	return nil
}

// getPermissonsByUserID User
func (d *userDao) getPermissonsByUserID(userID string) (*models.UserPermissions, *resterr.RestErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userDB := mongodb.Client.Database("erp-user-service")

	userPermList := models.UserPermissions{}
	userCollection := userDB.Collection("user-permissions")

	filter := bson.M{"user_id": userID}
	err := userCollection.FindOne(ctx, filter).Decode(&userPermList)
	if err != nil {
		return nil, resterr.NewInternalServerError(err.Error())
	}

	return &userPermList, nil
}

// updatePermissions of user
func (d *userDao) updatePermissions(user models.User) *resterr.RestErr {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userDB := mongodb.Client.Database("erp-user-service")
	userCollection := userDB.Collection("user")

	filter := bson.M{"user_id": user.ID}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "permissions", Value: user.Permissions},
		}},
	}

	_, err := userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return resterr.NewInternalServerError(err.Error())
	}
	return nil
}
