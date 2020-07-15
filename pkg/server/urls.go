package server

import (
	"github.com/abhidarbey/user-service/pkg/controllers/routes"
)

func mapUrls() {
	routes.Ping(router)
	routes.Auth(router)
	routes.Users(router)
	routes.Organization(router)
	routes.Department(router)
	routes.Role(router)
}
