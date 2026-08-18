package router

import (
	"net/http"

	"github.com/ashish9868/rapidbackend/core"
)

type apiRouter struct {
	App *core.App
}

func NewApiRouter(app *core.App) {
	r := &apiRouter{App: app}
	r.registerRoutes()
}

func (r *apiRouter) registerRoutes() {

	mux := http.NewServeMux()

	r.App.ResourceRoutes("health", mux, core.ResourceHandler{
		Index: &core.ResourceAction{
			Handler: func(w http.ResponseWriter, r *http.Request, app *core.App) {
				w.Write([]byte("hello"))
			},
		},
	})

	r.App.RootRouter.Handle("/api/v1/", http.StripPrefix("/api/v1", mux))
	// app := r.App

	// api_base_group := &http.ServeMux{}

	// loginHandler := superadmin.SuperAdminLoginHandler()
	// app.ResourceRoutes("login", api_base_group, *loginHandler)
	// app.ResourceRoutes("superadmins", api_base_group, *loginHandler, publicMiddlewares...)
	// app.ResourceRoutes("logout", api_base_group, *handlers.LogoutHandler())
	// app.ResourceRoutes("reset-password", api_base_group, *handlers.ResetPasswordHandler(), publicMiddlewares...)
	// app.ResourceRoutes("me", api_base_group, *handlers.MeHandler(), authMiddleware)
	// app.ResourceRoutes("projects", api_base_group, *handlers.NewProjectsHandler(), authMiddleware)
	// app.ResourceRoutes("collections", api_base_group, *handlers.NewCollectionsHandler(), authMiddleware)
	// app.ResourceRoutes("collections/:collection_id/fields", api_base_group, *handlers.NewFieldsHandler(), authMiddleware)
	// app.ResourceRoutes("collections/:collection_id/records", api_base_group, *handlers.NewRecordsHandler(), authMiddleware)
	// app.ResourceRoutes("users", api_base_group, *handlers.NewUsersHandler(), authMiddleware)

}
