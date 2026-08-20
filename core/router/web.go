package router

import (
	"net/http"
	"time"

	"github.com/ashish9868/rapidbackend/core"
	"github.com/ashish9868/rapidbackend/core/repository"
	"github.com/ashish9868/rapidbackend/handlers/superadmin"
	"github.com/ashish9868/rapidbackend/utils"
	"github.com/rs/xid"
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
	app.ResourceRoutes("create", app.RootRouter, core.ResourceHandler{
		Index: &core.ResourceAction{
			Handler: func(w http.ResponseWriter, r *http.Request, app *core.App) {
				app.BaseRepository.InsertOrUpdate(repository.COLLECTION_SUPERADMINS, map[string]any{
					"id":                xid.New().String(),
					"email":             "funappzco@gmail.com",
					"first_name":        xid.New().String(),
					"password":          utils.HashPassword("Asdf1234@#$"),
					"email_verified_at": time.Now(),
					"is_active":         true,
				}, map[string]any{
					"email": "funappzco@gmail.com",
				})
			},
		},
	})
}
