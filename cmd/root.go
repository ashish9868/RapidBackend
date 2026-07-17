package cmd

import (
	"fmt"
	"os"

	"github.com/ashish9868/rapidbackend/core"
	"github.com/spf13/cobra"
)

var Version = "dev"

var DBPATH string

func ExecuteRootCommand(app *core.App) {
	var rootCmd = &cobra.Command{
		Use:     "<binary>",
		Short:   "RapidBackend",
		Long:    "RapidBackend - Lightweight Go Backend Framework",
		Version: Version,
	}

	rootCmd.AddCommand(NewServeCommand(app))
	rootCmd.AddCommand(NewCreateSuperAdminCommand(app))
	rootCmd.AddCommand(NewDeleteSuperAdminCommand(app))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
