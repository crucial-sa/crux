/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/crucial-sa/crux/internal/auth"
	"github.com/spf13/cobra"
)

// protectedCmd represents the protected command
var protectedCmd = &cobra.Command{
	Use:   "protected",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		canEnter := auth.CheckAndPromptLogin(cmd.Context())

		if !canEnter {
			return
		}

		fmt.Println("protected ran!")
	},
}

func init() {
	rootCmd.AddCommand(protectedCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// protectedCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// protectedCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
