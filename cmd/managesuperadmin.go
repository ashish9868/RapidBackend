package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ashish9868/rapidbackend/core"
	"github.com/ashish9868/rapidbackend/core/respository"
	"github.com/ashish9868/rapidbackend/utils"
	"github.com/charmbracelet/x/term"
	"github.com/rs/xid"
	"github.com/spf13/cobra"
)

func NewCreateSuperAdminCommand(app *core.App) *cobra.Command {
	createSuperadminCmd := &cobra.Command{
		Use:   "managesuperadmin",
		Short: "Manage a superadmin",
		RunE: func(cmd *cobra.Command, args []string) error {

			reader := bufio.NewReader(os.Stdin)

			fmt.Println("Create Super Administrator")
			fmt.Println("--------------------------")

			fmt.Print("Email: ")
			email, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			email = strings.TrimSpace(email)

			if email == "" {
				return errors.New("email is required")
			}

			fmt.Print("Password: ")
			passwordBytes, err := term.ReadPassword(os.Stdin.Fd())
			if err != nil {
				return err
			}

			password := string(passwordBytes)

			fmt.Print("\nConfirm Password: ")
			confirmBytes, err := term.ReadPassword(os.Stdin.Fd())
			if err != nil {
				return err
			}
			fmt.Println()

			confirm := string(confirmBytes)

			if password != confirm {
				return errors.New("passwords do not match")
			}

			err = utils.ValidatePassword(password)
			if err != nil {
				return err
			}

			now := time.Now()

			err = app.BaseRepository.InsertOrUpdate(respository.COLLECTION_SUPERADMINS, map[string]any{
				"id":                xid.New().String(),
				"first_name":        xid.New().String(),
				"last_name":         email,
				"password":          utils.HashPassword(password),
				"email_verified_at": &now,
			}, map[string]any{
				"email": email,
			})

			if err == nil {
				utils.Log("✓ Super administrator created/updated successfully.")
			} else {
				utils.Log("x Unable to create/update superadmin.", err.Error())
			}
			return nil
		},
	}

	return createSuperadminCmd
}
