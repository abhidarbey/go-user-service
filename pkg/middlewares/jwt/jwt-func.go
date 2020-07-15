package jwt

import (
	"strings"

	"github.com/abhidarbey/user-service/pkg/models"
	"github.com/abhidarbey/user-service/pkg/utils/resterr"
)

// GenerateToken func
func GenerateToken(authUser *models.AuthUser) *models.ValueToken {
	// set TokenHeader
	header := headerEncoder()
	payload := payloadEncoder(authUser)
	secret := secretKey

	// Hash secret
	hashSec := Hash(payload, secret)
	tokenString := header + "." + payload + "." + hashSec

	valueToken := models.ValueToken{}
	valueToken.ValueToken = tokenString

	return &valueToken
}

// DecodeToken func
func DecodeToken(authHeader string) (*models.AuthUser, *resterr.RestErr) {
	// Check validity of authHeader
	if authHeader == "" {
		return nil, resterr.NewBadRequestError("Value token not provided")
	}

	authToken := strings.Split(authHeader, " ")
	if len(authToken) != 2 {
		return nil, resterr.NewBadRequestError("Malformed token")
	}

	tokenType := authToken[0]
	tokenString := authToken[1]

	// verify tokenType
	if tokenType != "JWT" {
		resterr.NewBadRequestError("Incorrect token type")
	}

	token := strings.Split(tokenString, ".")
	if len(token) != 3 {
		return nil, resterr.NewBadRequestError("Malformed token")
	}

	// header := token[0]
	payload := token[1]
	hashSec := token[2]

	secret := secretKey

	if !isValidHash(payload, hashSec, secret) {
		return nil, resterr.NewBadRequestError("Malformed token - Hash")
	}

	data, err := payloadDecoder(payload)
	if err != nil {
		return nil, err
	}

	if !isValidToken(data.Expiry) {
		return nil, resterr.NewBadRequestError("Token is expired")
	}

	user := models.AuthUser{}
	user.ID = data.ID
	user.Organization = data.Organization
	user.IsSuperuser = data.IsSuperuser
	user.IsOrgAdmin = data.IsOrgAdmin
	user.Permissions = data.Permissions

	return &user, nil
}
