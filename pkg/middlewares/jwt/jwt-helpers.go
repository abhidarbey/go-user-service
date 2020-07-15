package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/abhidarbey/user-service/pkg/models"
	"github.com/abhidarbey/user-service/pkg/utils/resterr"
)

const (
	secretKey = "supersecretkey"
)

// Hash generates a Hmac256 hash of a string using a secret
func Hash(value string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(value))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

// isValidHash validates a hash againt a value
func isValidHash(value string, hash string, secret string) bool {
	return hash == Hash(value, secret)
}

// base64Encoder generates a base64 encoded string
func base64Encoder(value []byte) string {
	return base64.URLEncoding.EncodeToString(value)
}

// Base64Decode takes in a base 64 encoded string and returns the //actual string or an error of it fails to decode the string
func base64Decode(payload string) ([]byte, *resterr.RestErr) {
	decoded, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		return nil, resterr.NewInternalServerError(fmt.Sprintf("Decoding Error %s", err))
	}
	return decoded, nil
}

// headerEncoder func
func headerEncoder() string {
	// set TokenHeader
	header := models.TokenHeader{}
	header.ALG = "H256"
	header.TYP = "JWT"

	jsonHeader, _ := json.Marshal(header)
	headerString := base64Encoder(jsonHeader)

	return headerString
}

// payloadEncoder func
func payloadEncoder(authUser *models.AuthUser) string {
	// set TokenPayload
	payload := models.TokenPayload{}
	payload.ID = authUser.ID
	payload.Organization = authUser.Organization
	payload.IsSuperuser = authUser.IsSuperuser
	payload.IsOrgAdmin = authUser.IsOrgAdmin
	payload.Permissions = authUser.Permissions
	payload.Authorized = true
	payload.Expiry = time.Now().Add(time.Minute * 15).UTC().Unix()

	jsonPayload, _ := json.Marshal(payload)
	payloadString := base64Encoder(jsonPayload)

	return payloadString
}

// payloadDecoder func
func payloadDecoder(payload string) (*models.TokenPayload, *resterr.RestErr) {
	decoded, err := base64Decode(payload)
	if err != nil {
		return nil, err
	}

	data := models.TokenPayload{}
	if err := json.Unmarshal(decoded, &data); err != nil {
		return nil, resterr.NewInternalServerError("Error when trying to unmarshal data response")
	}

	return &data, nil
}

// isValidToken validates the expiration of token
func isValidToken(exp int64) bool {
	return exp > time.Now().UTC().Unix()
}
