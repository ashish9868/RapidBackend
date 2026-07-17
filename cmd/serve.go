package cmd

import (
	"strconv"

	"github.com/ashish9868/rapidbackend/core"
	"github.com/ashish9868/rapidbackend/handlers"
	"github.com/spf13/cobra"
)

var (
	port int
)

func NewServeCommand(app *core.App) *cobra.Command {

	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start RapidBackend server",
		RunE: func(cmd *cobra.Command, args []string) error {

			api_base_group := app.Gin.Group("/api/v1")
			// routes

			app.ResourceRoutes("health", api_base_group, *handlers.HealthHandler())
			app.ResourceRoutes("projects", api_base_group, *handlers.NewProjectsHandler())
			app.ResourceRoutes("collections", api_base_group, *handlers.NewCollectionsHandler())
			app.ResourceRoutes("collections/:collection_id/fields", api_base_group, *handlers.NewFieldsHandler())
			app.ResourceRoutes("collections/:collection_id/records", api_base_group, *handlers.NewRecordsHandler())

			return app.Gin.Run(":" + strconv.Itoa(port))

		},
	}

	serveCmd.Flags().IntVar(&port, "port", 3000, "App Port")

	return serveCmd

}
