package router

import (
	"github.com/ashish9868/rapidbackend/core"
	"github.com/ashish9868/rapidbackend/handlers/superadmin"
)

type WebRouter struct {
	App *core.App
}

func (r *WebRouter) RegisterRoutes() {
	app := r.App

	suGroup := app.Gin.Group("/")
	app.ResourceRoutes("login", suGroup, *superadmin.SuperAdminLoginHandler())

	app.ServeStatic()
	app.ServeNoRoute()

}
