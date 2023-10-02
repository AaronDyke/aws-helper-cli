/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/AaronDyke/aws-helper-cli/pkg/aws"
	"github.com/AaronDyke/aws-helper-cli/pkg/dynamodb"
	"github.com/AaronDyke/aws-helper-cli/pkg/utils"
	"github.com/spf13/cobra"
)

// dynamodbListCmd represents the dynamodbList command
var dynamodbListCmd = &cobra.Command{
	Use:   "list",
	Short: "List DynamoDB tables",
	Long:  `List DynamoDB tables for a specific AWS profile.`,
	Run: func(cmd *cobra.Command, args []string) {
		profile := cmd.Flag("profile").Value.String()
		if profile == "" {
			profile = aws.PromptProfile()
		} else if !aws.ProfileExists(profile) {
			fmt.Println("Profile does not exist")
			return
		}

		region := cmd.Flag("region").Value.String()
		if region == "" {
			region = aws.PromptRegion(profile)
		}

		aws := aws.Aws{
			Profile: profile,
			Region:  region,
		}
		tables := dynamodb.ListTables(aws)

		if len(tables) == 0 {
			fmt.Println("No DynamoDB tables found")
			return
		}

		for _, table := range tables {
			fmt.Println(table)
		}

		cmd.Annotations["commandString"] = utils.CommandString("dynamodb list", map[string]string{"profile": profile, "region": region}, args)
	},
}

func init() {
	dynamodbCmd.AddCommand(dynamodbListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dynamodbListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dynamodbListCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
