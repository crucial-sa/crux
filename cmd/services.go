package cmd

import (
	"github.com/spf13/cobra"
)

var servicesCmd = &cobra.Command{
	Use:     "services",
	Short:   "Manage your crucial services",
	Aliases: []string{"service", "svc"},
}

func init() {
	rootCmd.AddCommand(servicesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// servicesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// servicesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
