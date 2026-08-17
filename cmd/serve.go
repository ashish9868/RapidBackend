package cmd

import (
	"strconv"

	"github.com/ashish9868/rapidbackend/core"
	"github.com/ashish9868/rapidbackend/core/router"
	"github.com/ashish9868/rapidbackend/utils"
	"github.com/spf13/cobra"
)

func NewServeCommand(app *core.App) *cobra.Command {

	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start RapidBackend server",
		RunE: func(cmd *cobra.Command, args []string) error {
			port := app.BaseUtil.SafeEnvGet("PORT", strconv.Itoa(utils.DEFAULT_PORT))
			apiRouter := &router.ApiRouter{App: app}
			webRouter := &router.WebRouter{App: app}
			apiRouter.RegisterRoutes()
			webRouter.RegisterRoutes()
			return app.Gin.Run(":" + port)
		},
	}
	return serveCmd
}
