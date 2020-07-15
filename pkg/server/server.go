package server

import (
	"log"
	"os"

	"github.com/abhidarbey/user-service/pkg/settings/db/mongodb"
	"github.com/abhidarbey/user-service/pkg/settings/seed"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

// StartApplication is ...
func StartApplication() {
	// Initiate mongodb
	mongodb.InitMongoClient()

	seedData := os.Getenv("SEED")

	if seedData == "true" {
		// Seed initial data
		seed.AddPermissions()
	}

	// Map all urls
	mapUrls()

	stage := os.Getenv("STAGE")

	if stage == "prod" {
		log.Printf("Listening and serving on port 8080")
		router.Run(":8080")
	} else {
		log.Printf("Listening and serving on port 5011")
		router.Run(":5011")
	}
}
