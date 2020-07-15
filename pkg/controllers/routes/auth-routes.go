package routes

import (
	"github.com/abhidarbey/user-service/pkg/controllers/handlers"
	"github.com/gin-gonic/gin"
)

// Auth Routes function
func Auth(r *gin.Engine) {
	h := handlers.AuthHandler

	router := r.Group("api")

	router.POST("login", h.Login)
	router.POST("register/org", h.RegisterOrg)
	router.POST("register/superuser", h.RegisterSuperuser)
}
