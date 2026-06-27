package cmd

import (
	"context"

	"github.com/crucial-sa/crux/internal/auth"
	"github.com/crucial-sa/crux/internal/services"
	"github.com/crucial-sa/crux/internal/ui"
	"github.com/spf13/cobra"
)

var serviceInfo services.Service

var servicesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new service",
	Run: func(cmd *cobra.Command, args []string) {
		session, canEnter := auth.CheckAndPromptLogin(cmd.Context())

		if !canEnter {
			return
		}
		err := services.PromptMissingServiceFields(&serviceInfo)
		if err != nil {
			ui.Panic("Failed to create service: %v", err)
		}

		if !services.ConfirmCreate(&serviceInfo) {
			ui.Say("Aborted.")
			return
		}

		err = ui.Spinner("Creating service...", cmd.Context(), func(ctx context.Context) error {
			return services.CreateService(ctx, session, &serviceInfo)
		})
		if err != nil {
			ui.Panic("Failed to create service: %v", err)
		}

		ui.Say("Service created.")
	},
}

func init() {
	servicesCmd.AddCommand(servicesCreateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// servicesCreateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	servicesCreateCmd.Flags().StringVar(&serviceInfo.Name, "name", "", "The name of the service")
	servicesCreateCmd.Flags().StringVar(&serviceInfo.Image, "image", "", "The image of the service")
}
