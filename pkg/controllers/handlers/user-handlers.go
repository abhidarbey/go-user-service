package handlers

import (
	"net/http"

	"github.com/abhidarbey/user-service/pkg/logic/services"
	"github.com/abhidarbey/user-service/pkg/middlewares/jwt"
	"github.com/abhidarbey/user-service/pkg/models"
	"github.com/abhidarbey/user-service/pkg/utils/resterr"
	"github.com/gin-gonic/gin"
)

// UserHandlerInterface type
type UserHandlerInterface interface {
	Create(ctx *gin.Context)
	FindAll(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	Update(ctx *gin.Context)
	UpdatePassword(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

// userHandler struct
type userHandler struct{}

// UserHandler variable
var (
	UserHandler UserHandlerInterface = &userHandler{}
)

// Create Handler
func (ctrl *userHandler) Create(ctx *gin.Context) {
	// Get JWT from request.Header
	authHeader := ctx.GetHeader("Authorization")
	authUser, err := jwt.DecodeToken(authHeader)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}
	var request models.User
	if err := ctx.ShouldBindJSON(&request); err != nil {
		restErr := resterr.NewBadRequestError("Invalid JSON body")
		ctx.JSON(restErr.Status, restErr)
		return
	}

	user, err := services.UserService.Create(request, authUser)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	response := gin.H{
		"object":  user.Marshal(),
		"message": "User successfully registered",
	}

	ctx.JSON(http.StatusOK, response)
}

// FindAll Handler
func (ctrl *userHandler) FindAll(ctx *gin.Context) {
	// Get JWT from request.Header
	authHeader := ctx.GetHeader("Authorization")
	authUser, err := jwt.DecodeToken(authHeader)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	users, err := services.UserService.FindAll(authUser)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	response := gin.H{
		"list":    users.Marshal(),
		"message": "List of all active users",
	}

	ctx.JSON(http.StatusOK, response)
}

// GetByID Handler
func (ctrl *userHandler) GetByID(ctx *gin.Context) {
	// Get JWT from request.Header
	authHeader := ctx.GetHeader("Authorization")
	authUser, err := jwt.DecodeToken(authHeader)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	// Get id from request.Param
	id := ctx.Param("id")

	user, err := services.UserService.GetByID(id, authUser)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	response := gin.H{
		"object":  user.Marshal(),
		"message": "User object",
	}

	ctx.JSON(http.StatusOK, response)
}

// Update Handler
func (ctrl *userHandler) Update(ctx *gin.Context) {
	// Get JWT from request.Header
	authHeader := ctx.GetHeader("Authorization")
	authUser, err := jwt.DecodeToken(authHeader)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	// Verify ID
	id := ctx.Param("id")

	// Verify body
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		restErr := resterr.NewBadRequestError("Invalid JSON body")
		ctx.JSON(restErr.Status, restErr)
		return
	}

	user.ID = id

	result, updateErr := services.UserService.Update(user, authUser)
	if updateErr != nil {
		ctx.JSON(updateErr.Status, updateErr)
		return
	}

	response := gin.H{
		"object":  result.Marshal(),
		"message": "User successfully updated",
	}

	ctx.JSON(http.StatusOK, response)
}

// UpdatePassword Handler
func (ctrl *userHandler) UpdatePassword(ctx *gin.Context) {
	// Get JWT from request.Header
	authHeader := ctx.GetHeader("Authorization")
	authUser, err := jwt.DecodeToken(authHeader)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	// Verify ID
	id := ctx.Param("id")

	// Verify body
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		restErr := resterr.NewBadRequestError("Invalid JSON body")
		ctx.JSON(restErr.Status, restErr)
		return
	}

	user.ID = id

	result, updateErr := services.UserService.UpdatePassword(user, authUser)
	if updateErr != nil {
		ctx.JSON(updateErr.Status, updateErr)
		return
	}

	response := gin.H{
		"object":  result.Marshal(),
		"message": "User successfully updated",
	}

	ctx.JSON(http.StatusOK, response)
}

// Delete Handler
func (ctrl *userHandler) Delete(ctx *gin.Context) {
	// Get JWT from request.Header
	authHeader := ctx.GetHeader("Authorization")
	authUser, err := jwt.DecodeToken(authHeader)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}
	// Verify ID
	id := ctx.Param("id")

	if err := services.UserService.Delete(id, authUser); err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	response := gin.H{
		"object":  map[string]string{"Status": "Deleted"},
		"message": "User successfully deleted",
	}

	ctx.JSON(http.StatusOK, response)
}
