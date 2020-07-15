package handlers

import (
	"net/http"

	"github.com/abhidarbey/user-service/pkg/logic/services"
	"github.com/abhidarbey/user-service/pkg/middlewares/jwt"
	"github.com/abhidarbey/user-service/pkg/models"
	"github.com/abhidarbey/user-service/pkg/utils/resterr"
	"github.com/gin-gonic/gin"
)

// DepartmentHandlerInterface type
type DepartmentHandlerInterface interface {
	Create(ctx *gin.Context)
	FindAll(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

// departmentHandler struct
type departmentHandler struct{}

// DepartmentHandler variable
var (
	DepartmentHandler DepartmentHandlerInterface = &departmentHandler{}
)

// Create Handler
func (ctrl *departmentHandler) Create(ctx *gin.Context) {
	// Get JWT from request.Header
	authHeader := ctx.GetHeader("Authorization")
	authUser, err := jwt.DecodeToken(authHeader)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}
	var request models.Department
	if err := ctx.ShouldBindJSON(&request); err != nil {
		restErr := resterr.NewBadRequestError("Invalid JSON body")
		ctx.JSON(restErr.Status, restErr)
		return
	}

	department, err := services.DepartmentService.Create(request, authUser)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	response := gin.H{
		"object":  department,
		"message": "Department successfully created",
	}

	ctx.JSON(http.StatusOK, response)
}

// FindAll Handler
func (ctrl *departmentHandler) FindAll(ctx *gin.Context) {
	// Get JWT --> authUser from request.Header
	authUser, err := jwt.DecodeToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	departments, err := services.DepartmentService.FindAll(authUser)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	response := gin.H{
		"list":    departments,
		"message": "List of all active departments",
	}

	ctx.JSON(http.StatusOK, response)
}

// GetByID department
func (ctrl *departmentHandler) GetByID(ctx *gin.Context) {
	// Get JWT --> authUser from request.Header
	authUser, err := jwt.DecodeToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	// Get id from request.Param
	id := ctx.Param("id")

	department, err := services.DepartmentService.GetByID(id, authUser)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	response := gin.H{
		"object":  department,
		"message": "Department object",
	}

	ctx.JSON(http.StatusOK, response)
}

// Update department
func (ctrl *departmentHandler) Update(ctx *gin.Context) {
	// Get JWT --> authUser from request.Header
	authUser, err := jwt.DecodeToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	// Get id from request.Param
	id := ctx.Param("id")

	// Verify body
	var request models.Department
	if err := ctx.ShouldBindJSON(&request); err != nil {
		reqErr := resterr.NewBadRequestError("Invalid JSON body")
		ctx.JSON(reqErr.Status, reqErr)
		return
	}

	request.ID = id

	department, updateErr := services.DepartmentService.Update(request, authUser)
	if updateErr != nil {
		ctx.JSON(updateErr.Status, updateErr)
		return
	}

	response := gin.H{
		"object":  department,
		"message": "Department updated",
	}

	ctx.JSON(http.StatusOK, response)
}

// Delete Handler
func (ctrl *departmentHandler) Delete(ctx *gin.Context) {
	// Get JWT --> authUser from request.Header
	authUser, err := jwt.DecodeToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	// Verify ID
	id := ctx.Param("id")

	if err := services.DepartmentService.Delete(id, authUser); err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	response := gin.H{
		"object":  map[string]string{"Status": "Deleted"},
		"message": "Department successfully deleted",
	}

	ctx.JSON(http.StatusOK, response)
}
