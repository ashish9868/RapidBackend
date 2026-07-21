package cmd

import (
	"github.com/ashish9868/rapidbackend/core"
	"github.com/spf13/cobra"
)

func NewTestCommand(app *core.App) *cobra.Command {
	testCmd := &cobra.Command{
		Use:   "test",
		Short: "Test random code",
		RunE: func(cmd *cobra.Command, args []string) error {
			println(app.BaseUtil.CheckPassword("$2a$10$XuTjyhzB0qCMQiQiFc0yQu1XNqp3VxwkdIoKDaDpBOXHFbgQgNPme", "Asdf1234@#$"))
			return nil
		},
	}
	return testCmd

}
