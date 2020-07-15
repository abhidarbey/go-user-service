package routes

import (
	"github.com/abhidarbey/user-service/pkg/controllers/handlers"
	"github.com/gin-gonic/gin"
)

// Role Routes function
func Role(r *gin.Engine) {
	h := handlers.RoleHandler

	router := r.Group("/api/role")

	router.POST("", h.Create)
	router.GET("", h.FindAll)
	router.GET(":id", h.GetByID)
	router.PUT(":id", h.Update)
	router.DELETE(":id", h.Delete)
}
