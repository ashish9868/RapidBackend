package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ashish9868/rapidbackend/core"
	"github.com/ashish9868/rapidbackend/models"
	"github.com/charmbracelet/x/term"
	"github.com/rs/xid"
	"github.com/spf13/cobra"
)

var (
	email    string
	password string
)

func NewCreateSuperAdminCommand(app *core.App) *cobra.Command {
	createSuperadminCmd := &cobra.Command{
		Use:   "createsuperadmin",
		Short: "Create a superadmin",
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

			exists := models.SuperAdmin{}
			total, err := app.Bun.NewSelect().Model(&exists).Where("email = ?", email).Count(context.Background())

			if err != nil {
				return err
			}

			if total > 0 {
				return errors.New("\n\nemail already exists\n\n")
			}

			fmt.Print("Password: ")
			passwordBytes, err := term.ReadPassword(os.Stdin.Fd())
			if err != nil {
				return err
			}

			password := string(passwordBytes)

			fmt.Println("\n")
			fmt.Print("Confirm Password: ")
			confirmBytes, err := term.ReadPassword(os.Stdin.Fd())
			if err != nil {
				return err
			}
			fmt.Println()

			confirm := string(confirmBytes)

			if password != confirm {
				return errors.New("passwords do not match")
			}

			err = app.BaseUtil.ValidatePassword(password)
			if err != nil {
				return err
			}

			verifiedAt := time.Now()
			_, err = app.Bun.NewInsert().Model(&models.SuperAdmin{
				ID:              xid.New().String(),
				FirstName:       xid.New().String(),
				Email:           email,
				Password:        app.BaseUtil.HashPassword(password),
				EmailVerifiedAt: &verifiedAt,
				IsActive:        true,
			}).Exec(context.Background())

			if err != nil {
				fmt.Println("Unable to delete user : %s", err.Error())
				return err
			}

			fmt.Println("✓ Super administrator created successfully.")

			return nil
		},
	}

	return createSuperadminCmd
}

func NewDeleteSuperAdminCommand(app *core.App) *cobra.Command {

	removeSuperadminCmd := &cobra.Command{
		Use:   "removesuperadmin",
		Short: "Delete a superadmin",
		RunE: func(cmd *cobra.Command, args []string) error {

			reader := bufio.NewReader(os.Stdin)

			fmt.Println("Remove Super Administrator")
			fmt.Println("--------------------------")

			fmt.Print("Email: ")
			email, err := reader.ReadString('\n')
			if err != nil {
				return err
			}

			println(email)
			_, err = app.Bun.NewDelete().Model((*models.SuperAdmin)(nil)).Where("email = ?", strings.TrimSpace(email)).Exec(context.Background())

			if err != nil {
				fmt.Println("Unable to delete user : " + err.Error())
				return err
			}

			fmt.Println("✓ Super administrator deleted successfully.")

			return nil
		},
	}
	return removeSuperadminCmd
}
