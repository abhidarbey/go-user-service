package handlers

import (
	"net/http"

	"github.com/abhidarbey/user-service/pkg/logic/services"
	"github.com/abhidarbey/user-service/pkg/models"
	"github.com/abhidarbey/user-service/pkg/utils/resterr"
	"github.com/gin-gonic/gin"
)

// AuthHandlerInterface type
type AuthHandlerInterface interface {
	Login(ctx *gin.Context)
	RegisterOrg(ctx *gin.Context)
	RegisterSuperuser(ctx *gin.Context)
}

// authHandler struct
type authHandler struct{}

// AuthHandler variable
var (
	AuthHandler AuthHandlerInterface = &authHandler{}
)

// Login Handler
func (ctrl *authHandler) Login(ctx *gin.Context) {

	var request models.LoginRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		reqErr := resterr.NewBadRequestError("Invalid JSON body")
		ctx.JSON(reqErr.Status, reqErr)
		return
	}

	user, err := services.AuthService.Login(request)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	response := gin.H{
		"object":  user.AuthMarshal(),
		"message": "User logged in successfully!",
	}

	ctx.JSON(http.StatusAccepted, response)
}

// RegisterOrg Handler
func (ctrl *authHandler) RegisterOrg(ctx *gin.Context) {

	var request models.RegistrationRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		reqErr := resterr.NewBadRequestError("Invalid JSON body")
		ctx.JSON(reqErr.Status, reqErr)
		return
	}

	organization, err := services.AuthService.RegisterOrg(request)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	response := gin.H{
		"object":  organization,
		"message": "Organization created successfully",
	}

	ctx.JSON(http.StatusCreated, response)
}

// Register Handler
func (ctrl *authHandler) RegisterSuperuser(ctx *gin.Context) {

	var request models.User
	if err := ctx.ShouldBindJSON(&request); err != nil {
		reqErr := resterr.NewBadRequestError("Invalid JSON body")
		ctx.JSON(reqErr.Status, reqErr)
		return
	}

	superuser, err := services.AuthService.RegisterSuperuser(request)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	response := gin.H{
		"object":  superuser.Marshal(),
		"message": "Superuser registration successful",
	}

	ctx.JSON(http.StatusCreated, response)
}
