package cmd

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ashish9868/rapidbackend/core"
	"github.com/ashish9868/rapidbackend/handlers"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

func NewServeCommand(app *core.App) *cobra.Command {

	template_data := map[string]any{
		"Version": app.Version,
		"Year":    time.Now().Year(),
		"AppName": "FastBackend",
	}
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start RapidBackend server",
		RunE: func(cmd *cobra.Command, args []string) error {

			if app.FeFs != nil {
				app.Gin.LoadHTMLFS(http.FS(*app.FeFs), "templates/**/*")
				app.BaseUtil.PrintFiles(*app.FeFs)
			}
			api_base_group := app.Gin.Group("/api/v1")
			// routes
			app.ResourceRoutes("health", api_base_group, *handlers.HealthHandler())
			app.ResourceRoutes("projects", api_base_group, *handlers.NewProjectsHandler())
			app.ResourceRoutes("collections", api_base_group, *handlers.NewCollectionsHandler())
			app.ResourceRoutes("collections/:collection_id/fields", api_base_group, *handlers.NewFieldsHandler())
			app.ResourceRoutes("collections/:collection_id/records", api_base_group, *handlers.NewRecordsHandler())

			// other handling

			if app.FeFs != nil {
				template404 := "templates/pages/not_found"
				// static
				staticSub := app.BaseUtil.SubFs(*app.FeFs, "static")
				if staticSub != nil {
					app.Gin.StaticFS("/static", http.FS(*staticSub))
				}
				// web
				app.Gin.NoRoute(func(ctx *gin.Context) {
					if app.FeFs != nil {
						path := strings.TrimPrefix(strings.TrimSpace(ctx.Request.URL.Path), "/")
						if len(path) < 2 {
							path = "templates/pages/home"
						} else if strings.HasPrefix(path, "/partials") {
							path = "templates" + path
						} else {
							path = "templates/pages/" + path
						}
						fmt.Println("loading template", path)
						if app.BaseUtil.FileExists(*app.FeFs, path+".html") {
							ctx.HTML(http.StatusOK, path, template_data)
						} else {
							ctx.HTML(http.StatusNotFound, template404, template_data)
						}
						return
					}
					ctx.HTML(http.StatusNotFound, template404, template_data)
				})
			}
			port := app.BaseUtil.SafeEnvGet("PORT", "7007")
			return app.Gin.Run(":" + port)

		},
	}
	return serveCmd
}
