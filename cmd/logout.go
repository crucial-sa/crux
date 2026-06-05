package cmd

import (
	"github.com/crucial-sa/crux/internal/auth"
	"github.com/crucial-sa/crux/internal/ui"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout of the current logged in user",
	Run: func(cmd *cobra.Command, args []string) {
		secret, err := auth.GetSecret()
		if err != nil {
			ui.Panic("Failed to check login status", err)
		}

		if secret != "" {
			err := auth.ClearSecret()
			if err != nil {
				ui.Panic("Failed to clear secret", err)
			}
		}

		ui.Say("Logged out successfully!")
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
