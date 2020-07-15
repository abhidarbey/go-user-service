package dao

import (
	"context"
	"time"

	"github.com/abhidarbey/user-service/pkg/models"
	"github.com/abhidarbey/user-service/pkg/settings/db/mongodb"
	"github.com/abhidarbey/user-service/pkg/utils/resterr"
	"go.mongodb.org/mongo-driver/bson"
)

// AuthDaoInterface type
type AuthDaoInterface interface {
	Login(string, string) (*models.User, *resterr.RestErr)
}

type authDao struct{}

// AuthDao variable
var (
	AuthDao AuthDaoInterface = &authDao{}
)

// Login auth
func (d *authDao) Login(email string, password string) (*models.User, *resterr.RestErr) {
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

	if user.Password != password {
		return nil, resterr.NewBadRequestError("Incorrect password")
	}

	return &user, nil
}
