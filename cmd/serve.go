package cmd

import (
	"context"
	"strconv"

	"github.com/ashish9868/rapidbackend/core"
	"github.com/ashish9868/rapidbackend/handlers"
	"github.com/ashish9868/rapidbackend/templates/pages"
	"github.com/ashish9868/rapidbackend/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

func NewServeCommand(app *core.App) *cobra.Command {

	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start RapidBackend server",
		RunE: func(cmd *cobra.Command, args []string) error {

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

			app.Gin.GET("/", func(ctx *gin.Context) {
				pages.LoginPage().Render(context.Background(), ctx.Writer)
			})

			app.ServeStatic()
			app.ServeNoRoute()

			port := app.BaseUtil.SafeEnvGet("PORT", strconv.Itoa(utils.DEFAULT_PORT))
			return app.Gin.Run(":" + port)
		},
	}
	return serveCmd
}
