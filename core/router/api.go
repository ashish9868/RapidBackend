package router

import (
	"github.com/ashish9868/rapidbackend/core"
	"github.com/ashish9868/rapidbackend/handlers"
	"github.com/gin-gonic/gin"
)

type ApiRouter struct {
	App *core.App
}

func (r *ApiRouter) RegisterRoutes() {
	app := r.App
	api_base_group := app.Gin.Group("/api/v1")
	publicMiddlewares := []gin.HandlerFunc{app.NewAuthMiddleWare(false, true), app.NewPublicMiddleware()}
	authMiddleware := app.NewAuthMiddleWare(true, false)
	loginHandler := handlers.LoginHandler("superadmins")
	app.ResourceRoutes("login", api_base_group, *loginHandler, publicMiddlewares...)
	app.ResourceRoutes("superadmins", api_base_group, *loginHandler, publicMiddlewares...)
	app.ResourceRoutes("logout", api_base_group, *handlers.LogoutHandler())
	app.ResourceRoutes("reset-password", api_base_group, *handlers.ResetPasswordHandler(), publicMiddlewares...)
	app.ResourceRoutes("me", api_base_group, *handlers.MeHandler(), authMiddleware)
	app.ResourceRoutes("projects", api_base_group, *handlers.NewProjectsHandler(), authMiddleware)
	app.ResourceRoutes("collections", api_base_group, *handlers.NewCollectionsHandler(), authMiddleware)
	app.ResourceRoutes("collections/:collection_id/fields", api_base_group, *handlers.NewFieldsHandler(), authMiddleware)
	app.ResourceRoutes("collections/:collection_id/records", api_base_group, *handlers.NewRecordsHandler(), authMiddleware)
	app.ResourceRoutes("users", api_base_group, *handlers.NewUsersHandler(), authMiddleware)

}
