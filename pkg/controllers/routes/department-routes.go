package routes

import (
	"github.com/abhidarbey/user-service/pkg/controllers/handlers"
	"github.com/gin-gonic/gin"
)

// Department Routes function
func Department(r *gin.Engine) {
	h := handlers.DepartmentHandler

	router := r.Group("/api/department")

	router.POST("", h.Create)
	router.GET("", h.FindAll)
	router.GET(":id", h.GetByID)
	router.PUT(":id", h.Update)
	router.DELETE(":id", h.Delete)
}
