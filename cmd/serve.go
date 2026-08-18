package cmd

import (
	"net/http"
	"time"

	"github.com/ashish9868/rapidbackend/constants"
	"github.com/ashish9868/rapidbackend/core"
	"github.com/ashish9868/rapidbackend/core/router"
	"github.com/ashish9868/rapidbackend/middlewares"
	"github.com/ashish9868/rapidbackend/utils"
	"github.com/spf13/cobra"
)

func NewServeCommand(app *core.App) *cobra.Command {

	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start RapidBackend server",
		RunE: func(cmd *cobra.Command, args []string) error {
			port := utils.SafeEnvGet("PORT", constants.DEFAULT_PORT)

			router.NewApiRouter(app)
			router.NewWebRouter(app)

			server := &http.Server{
				Addr: ":" + utils.ToString(port),
				Handler: middlewares.Chain(
					app.RootRouter,
					middlewares.Recovery,
					middlewares.Logger,
				),
				ReadHeaderTimeout: 5 * time.Second,
				ReadTimeout:       30 * time.Second,
				WriteTimeout:      30 * time.Second,
				IdleTimeout:       60 * time.Second,
			}

			return server.ListenAndServe()
		},
	}
	return serveCmd
}
