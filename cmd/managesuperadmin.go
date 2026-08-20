package cmd

import (
	"time"

	"github.com/ashish9868/rapidbackend/core"
	"github.com/ashish9868/rapidbackend/core/repository"
	"github.com/ashish9868/rapidbackend/utils"
	"github.com/spf13/cobra"
)

func NewCreateSuperAdminCommand(app *core.App) *cobra.Command {
	var email string
	createSuperadminCmd := &cobra.Command{
		Use:   "managesuperadmin",
		Short: "Manage a superadmin",
		RunE: func(cmd *cobra.Command, args []string) error {

			now := time.Now()
			password := utils.GenerateRandomPassword(10)
			err := app.BaseRepository.InsertOrUpdate(repository.COLLECTION_SUPERADMINS, map[string]any{
				"email":             email,
				"password":          utils.HashPassword("Asdf1234@#$"),
				"email_verified_at": &now,
				"is_active":         true,
			}, map[string]any{
				"email": email,
			})

			if err == nil {
				utils.Log("✓ Super administrator created/updated successfully.")
				utils.LogF("Generated Password (Copy it) is : \n\n%s\n\n", password)
			} else {
				utils.Log("x Unable to create/update superadmin.", err.Error())
			}
			return nil
		},
	}

	createSuperadminCmd.Flags().StringVar(&email, "email", "", "--email something@example.com")
	createSuperadminCmd.MarkFlagRequired("email")
	return createSuperadminCmd
}
