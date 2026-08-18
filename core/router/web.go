package router

import (
	"net/http"

	"github.com/ashish9868/rapidbackend/core"
	"github.com/ashish9868/rapidbackend/handlers/superadmin"
)

type webRouter struct {
	App    *core.App
	Router *http.ServeMux
}

func NewWebRouter(app *core.App) {
	r := &webRouter{App: app, Router: &http.ServeMux{}}
	r.registerRoutes()
}

func (r *webRouter) registerRoutes() {
	app := r.App
	app.ResourceRoutes("login", app.RootRouter, *superadmin.SuperAdminLoginHandler())
}
