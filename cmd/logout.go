package cmd

import (
	"fmt"

	"github.com/crucial-sa/crux/internal/auth"
	"github.com/crucial-sa/crux/internal/ui"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout of the current logged in user",
	Run: func(cmd *cobra.Command, args []string) {
		err := auth.ClearSecret()
		if err != nil {
			fmt.Printf("Failed to clear secret")
			return
		}

		ui.Say("Logged in successfully!")
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
