package cmd

import (
	"fmt"
	"time"

	"github.com/ashish9868/rapidbackend/core"
	"github.com/ashish9868/rapidbackend/utils"
	"github.com/spf13/cobra"
)

func CreateNewMigration(name string) error {
	now := time.Now().UTC()
	prefix := now.Format("20060102150405")

	upFile := fmt.Sprintf("%s/%s_%s.sql", "database/migrations", prefix, name)

	err := utils.SafeCreateFile(upFile, fmt.Sprintf("---%s---", upFile))

	if err != nil {
		return err
	}

	fmt.Println("Created:")
	fmt.Println(" ", upFile)

	return nil

}

var migration_name string

func NewMigrationGenerateCommand(app *core.App) *cobra.Command {
	testCmd := &cobra.Command{
		Use:   "migrate:make",
		Short: "Make New Migration File",
		RunE: func(cmd *cobra.Command, args []string) error {
			CreateNewMigration(migration_name)
			return nil
		},
	}
	testCmd.Flags().StringVar(&migration_name, "name", "", "name is (required)")
	testCmd.MarkFlagRequired("name")
	return testCmd
}

func MigrateDatabaseCommand(app *core.App) *cobra.Command {
	testCmd := &cobra.Command{
		Use:   "migrate:latest",
		Short: "Migrate database migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			app.DBMigrate()

			utils.Log("\nSuccessfully Migrated")
			return nil
		},
	}
	return testCmd
}
