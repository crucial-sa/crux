package cmd

import (
	"github.com/crucial-sa/crux/internal/auth"
	"github.com/crucial-sa/crux/internal/services"
	"github.com/crucial-sa/crux/internal/ui"
	"github.com/spf13/cobra"
)

var servicesListCmd = &cobra.Command{
	Use:     "list",
	Short:   "lists all your services",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("servicesList called")
		session, canEnter := auth.CheckAndPromptLogin(cmd.Context())

		if !canEnter {
			return
		}

		svcs, err := services.GetServices(session)
		if err != nil {
			ui.Panic("Failed to get services: %v", err)
		}

		services.PrintServicesTable(svcs)
	},
}

func init() {
	servicesCmd.AddCommand(servicesListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// servicesListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// servicesListCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
