package routes

import (
	"github.com/abhidarbey/user-service/pkg/controllers/handlers"
	"github.com/gin-gonic/gin"
)

// Organization Routes function
func Organization(r *gin.Engine) {
	h := handlers.OrganizationHandler

	router := r.Group("/api/org")

	router.GET("", h.FindAll)
	router.GET(":id", h.GetByID)
	router.PUT(":id", h.Update)
	router.DELETE(":id", h.Delete)
}
