package handlers

import (
	"net/http"

	"github.com/abhidarbey/user-service/pkg/logic/services"
	"github.com/abhidarbey/user-service/pkg/middlewares/jwt"
	"github.com/abhidarbey/user-service/pkg/models"
	"github.com/abhidarbey/user-service/pkg/utils/resterr"
	"github.com/gin-gonic/gin"
)

// OrganizationHandlerInterface type
type OrganizationHandlerInterface interface {
	FindAll(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

// organizationHandler struct
type organizationHandler struct{}

// OrganizationHandler variable
var (
	OrganizationHandler OrganizationHandlerInterface = &organizationHandler{}
)

// FindAll Handler
func (ctrl *organizationHandler) FindAll(ctx *gin.Context) {
	// Get JWT --> authUser from request.Header
	authUser, err := jwt.DecodeToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	org, err := services.OrganizationService.FindAll(authUser)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	response := gin.H{
		"list":    org,
		"message": "List of all active organizations",
	}

	ctx.JSON(http.StatusOK, response)
}

// GetByID organization
func (ctrl *organizationHandler) GetByID(ctx *gin.Context) {
	// Get JWT --> authUser from request.Header
	authUser, err := jwt.DecodeToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	// Get id from request.Param
	id := ctx.Param("id")

	org, err := services.OrganizationService.GetByID(id, authUser)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	response := gin.H{
		"object":  org,
		"message": "Organization object",
	}

	ctx.JSON(http.StatusOK, response)
}

// Update organization
func (ctrl *organizationHandler) Update(ctx *gin.Context) {
	// Get JWT --> authUser from request.Header
	authUser, err := jwt.DecodeToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	// Get id from request.Param
	id := ctx.Param("id")

	// Verify body
	var request models.Organization
	if err := ctx.ShouldBindJSON(&request); err != nil {
		reqErr := resterr.NewBadRequestError("Invalid JSON body")
		ctx.JSON(reqErr.Status, reqErr)
		return
	}

	request.ID = id

	org, updateErr := services.OrganizationService.Update(request, authUser)
	if updateErr != nil {
		ctx.JSON(updateErr.Status, updateErr)
		return
	}

	response := gin.H{
		"object":  org,
		"message": "Organization updated",
	}

	ctx.JSON(http.StatusOK, response)
}

// Delete Handler
func (ctrl *organizationHandler) Delete(ctx *gin.Context) {
	// Get JWT --> authUser from request.Header
	authUser, err := jwt.DecodeToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	// Verify ID
	id := ctx.Param("id")

	if err := services.OrganizationService.Delete(id, authUser); err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	response := gin.H{
		"object":  map[string]string{"Status": "Deleted"},
		"message": "Organization successfully deleted",
	}

	ctx.JSON(http.StatusOK, response)
}
