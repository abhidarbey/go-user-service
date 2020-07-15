package handlers

import (
	"net/http"

	"github.com/abhidarbey/user-service/pkg/logic/services"
	"github.com/abhidarbey/user-service/pkg/middlewares/jwt"
	"github.com/abhidarbey/user-service/pkg/models"
	"github.com/abhidarbey/user-service/pkg/utils/resterr"
	"github.com/gin-gonic/gin"
)

// RoleHandlerInterface type
type RoleHandlerInterface interface {
	Create(ctx *gin.Context)
	FindAll(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

// roleHandler struct
type roleHandler struct{}

// RoleHandler variable
var (
	RoleHandler RoleHandlerInterface = &roleHandler{}
)

// Create Handler
func (ctrl *roleHandler) Create(ctx *gin.Context) {
	// Get JWT from request.Header
	authHeader := ctx.GetHeader("Authorization")
	authUser, err := jwt.DecodeToken(authHeader)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}
	var request models.Role
	if err := ctx.ShouldBindJSON(&request); err != nil {
		restErr := resterr.NewBadRequestError("Invalid JSON body")
		ctx.JSON(restErr.Status, restErr)
		return
	}

	role, err := services.RoleService.Create(request, authUser)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	response := gin.H{
		"object":  role,
		"message": "Role successfully created",
	}

	ctx.JSON(http.StatusOK, response)
}

// FindAll Handler
func (ctrl *roleHandler) FindAll(ctx *gin.Context) {
	// Get JWT --> authUser from request.Header
	authUser, err := jwt.DecodeToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	roles, err := services.RoleService.FindAll(authUser)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	response := gin.H{
		"list":    roles,
		"message": "List of all active roles",
	}

	ctx.JSON(http.StatusOK, response)
}

// GetByID role
func (ctrl *roleHandler) GetByID(ctx *gin.Context) {
	// Get JWT --> authUser from request.Header
	authUser, err := jwt.DecodeToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	// Get id from request.Param
	id := ctx.Param("id")

	role, err := services.RoleService.GetByID(id, authUser)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	response := gin.H{
		"object":  role,
		"message": "Role object",
	}

	ctx.JSON(http.StatusOK, response)
}

// Update role
func (ctrl *roleHandler) Update(ctx *gin.Context) {
	// Get JWT --> authUser from request.Header
	authUser, err := jwt.DecodeToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	// Get id from request.Param
	id := ctx.Param("id")

	// Verify body
	var request models.Role
	if err := ctx.ShouldBindJSON(&request); err != nil {
		reqErr := resterr.NewBadRequestError("Invalid JSON body")
		ctx.JSON(reqErr.Status, reqErr)
		return
	}

	request.ID = id

	role, updateErr := services.RoleService.Update(request, authUser)
	if updateErr != nil {
		ctx.JSON(updateErr.Status, updateErr)
		return
	}

	response := gin.H{
		"object":  role,
		"message": "Role updated",
	}

	ctx.JSON(http.StatusOK, response)
}

// Delete Handler
func (ctrl *roleHandler) Delete(ctx *gin.Context) {
	// Get JWT --> authUser from request.Header
	authUser, err := jwt.DecodeToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	// Verify ID
	id := ctx.Param("id")

	if err := services.RoleService.Delete(id, authUser); err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	response := gin.H{
		"object":  map[string]string{"Status": "Deleted"},
		"message": "Role successfully deleted",
	}

	ctx.JSON(http.StatusOK, response)
}
