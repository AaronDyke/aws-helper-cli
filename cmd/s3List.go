/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/AaronDyke/aws-helper-cli/pkg/aws"
	"github.com/AaronDyke/aws-helper-cli/pkg/s3"
	"github.com/spf13/cobra"
)

// s3ListCmd represents the s3List command
var s3ListCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		profile := cmd.Flag("profile").Value.String()
		if profile == "" {
			profile = aws.PromptProfile()
		} else if !aws.ProfileExists(profile) {
			fmt.Println("Profile does not exist")
			return
		}

		fmt.Println(s3.ListBuckets(profile))
	},
}

func init() {
	s3Cmd.AddCommand(s3ListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// s3ListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// s3ListCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
