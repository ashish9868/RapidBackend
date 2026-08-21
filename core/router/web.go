package router

import (
	"net/http"

	"github.com/ashish9868/rapidbackend/core"
	"github.com/ashish9868/rapidbackend/handlers"
	"github.com/ashish9868/rapidbackend/handlers/superadmin"
	"github.com/ashish9868/rapidbackend/middlewares"
	"github.com/ashish9868/rapidbackend/models"
	"github.com/ashish9868/rapidbackend/utils"
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
	authMiddleware := middlewares.AuthMiddleware(app.AccessTokenRepository, true, false)
	app.ResourceRoutes("login", app.RootRouter, *superadmin.SuperAdminLoginHandler())
	app.ResourceRoutes("logout", app.RootRouter, core.ResourceHandler{
		Index: &core.ResourceAction{
			Handler: func(w http.ResponseWriter, r *http.Request, app *core.App) {

				app.Logout(w, r)
			},
		},
	})
	app.ResourceRoutes("dashboard", app.RootRouter, *handlers.DashboardHandler(), authMiddleware)
	app.ResourceRoutes("create", app.RootRouter, core.ResourceHandler{
		Index: &core.ResourceAction{
			Handler: func(w http.ResponseWriter, r *http.Request, app *core.App) {
				err := app.BaseRepository.GetByColumn("superadmins", "email", "funappzco@gmail.com", &models.User{})
				if err != nil {
					utils.LogF("Error: %s", err)
				}
			},
		},
	})
}
