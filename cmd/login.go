package cmd

import (
	"context"

	"github.com/crucial-sa/crux/internal/auth"
	"github.com/crucial-sa/crux/internal/ui"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login as a crucial user. You need to be authenticated to use most commands",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := auth.Login(context.Background())
		if err != nil {
			ui.Panic("Failed to login", err)
		}

		ui.Say("Logged in successfully!")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
