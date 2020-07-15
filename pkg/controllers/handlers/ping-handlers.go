package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping test
func Ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "pong")
}

// King test
func King(ctx *gin.Context) {
	ctx.String(http.StatusOK, "kong")
}
