package cmd

import (
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/ashish9868/rapidbackend/core"
	"github.com/ashish9868/rapidbackend/handlers"
	"github.com/ashish9868/rapidbackend/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

func NewServeCommand(app *core.App) *cobra.Command {

	template_data := gin.H{
		"Version":     app.Version,
		"Year":        time.Now().Year(),
		"AppName":     "FastBackend",
		"DefaultPort": utils.DEFAULT_PORT,
		"Now":         time.Now().UnixMicro(),
	}
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start RapidBackend server",
		RunE: func(cmd *cobra.Command, args []string) error {

			if app.FeFs != nil {
				app.Gin.SetFuncMap(template.FuncMap{
					"coalesce": app.BaseUtil.Coalesce,
					"merge":    app.BaseUtil.Merge,
					"dict":     app.BaseUtil.Dict,
					"toJSON":   app.BaseUtil.ToJSON,
					"debug":    app.BaseUtil.Debug,
				})
				app.Gin.LoadHTMLFS(http.FS(*app.FeFs), "templates/**/*")
				// static serve
				app.ServeStatic()
			}
			api_base_group := app.Gin.Group("/api/v1")
			// routes
			app.ResourceRoutes("health", api_base_group, *handlers.HealthHandler())
			app.ResourceRoutes("projects", api_base_group, *handlers.NewProjectsHandler())
			app.ResourceRoutes("collections", api_base_group, *handlers.NewCollectionsHandler())
			app.ResourceRoutes("collections/:collection_id/fields", api_base_group, *handlers.NewFieldsHandler())
			app.ResourceRoutes("collections/:collection_id/records", api_base_group, *handlers.NewRecordsHandler())

			app.RenderPage("/", template_data)
			app.RenderPage("/dashboard", template_data)
			app.RenderPage("/settings", template_data)
			app.RenderPage("/projects", template_data)
			app.RenderPage("/collections", template_data)
			app.RenderFragments(template_data)
			app.ServeNoRoute(template_data)

			port := app.BaseUtil.SafeEnvGet("PORT", strconv.Itoa(utils.DEFAULT_PORT))
			return app.Gin.Run(":" + port)

		},
	}
	return serveCmd
}
