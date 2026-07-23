package cmd

import (
	"strconv"

	"github.com/ashish9868/rapidbackend/core"
	"github.com/ashish9868/rapidbackend/handlers"
	"github.com/ashish9868/rapidbackend/middlewares"
	"github.com/ashish9868/rapidbackend/utils"
	"github.com/spf13/cobra"
)

func NewServeCommand(app *core.App) *cobra.Command {

	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start RapidBackend server",
		RunE: func(cmd *cobra.Command, args []string) error {

			// static serve
			app.ServeStatic()
			api_base_group := app.Gin.Group("/api/v1")
			// routes
			app.ResourceRoutes("login", api_base_group, *handlers.LoginHandler())
			app.ResourceRoutes("reset-password", api_base_group, *handlers.ResetPasswordHandler())
			app.ResourceRoutes("me", api_base_group, *handlers.MeHandler(), middlewares.NewAuthMiddleWare(app, true))
			app.ResourceRoutes("projects", api_base_group, *handlers.NewProjectsHandler())
			app.ResourceRoutes("collections", api_base_group, *handlers.NewCollectionsHandler())
			app.ResourceRoutes("collections/:collection_id/fields", api_base_group, *handlers.NewFieldsHandler())
			app.ResourceRoutes("collections/:collection_id/records", api_base_group, *handlers.NewRecordsHandler())

			app.ServeNoRoute()

			port := app.BaseUtil.SafeEnvGet("PORT", strconv.Itoa(utils.DEFAULT_PORT))
			return app.Gin.Run(":" + port)

		},
	}
	return serveCmd
}
