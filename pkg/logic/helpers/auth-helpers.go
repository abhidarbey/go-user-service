package helpers

import "github.com/abhidarbey/user-service/pkg/models"

// IsGranted middlewares verifies the permission of the user
func IsGranted(permission string, user models.AuthUser) bool {
	for i := 0; i < len(user.Permissions); i++ {
		if permission == user.Permissions[i].Name {
			return true
		}
	}
	return false
}
