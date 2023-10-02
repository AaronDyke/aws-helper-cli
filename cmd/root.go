/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "aws-helper-cli",
	Short: "CLI to interact with AWS",
	Long:  `A CLI to interact with AWS that allows for interactions beyond the AWS CLI.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cmd.Annotations = make(map[string]string)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("quiet").Value.String() != "true" {
			fmt.Println("To run this exact command again, use the following:")
			fmt.Println(cmd.Annotations["commandString"])
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().String("profile", "", "AWS Profile to use")
	rootCmd.PersistentFlags().String("region", "", "AWS Profile to use")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Quiet mode")
}
