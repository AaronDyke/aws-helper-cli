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

// dynamodbPutCmd represents the dynamodbPut command
var dynamodbPutCmd = &cobra.Command{
	Use:   "put",
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

		region := cmd.Flag("region").Value.String()
		if region == "" {
			region = aws.PromptRegion(profile)
		}

		aws := aws.Aws{
			Profile: profile,
			Region:  region,
		}

		table := cmd.Flag("table").Value.String()
		if table == "" {
			table = dynamodb.PromptTables(aws, "Select DynamoDB Table", []string{})
		} else if !dynamodb.TableExists(aws, table) {
			fmt.Println("Table does not exist")
			return
		}

		item := cmd.Flag("item").Value.String()
		dynamodb.PutItem(aws, table, item)

		cmd.Annotations["commandString"] = utils.CommandString("dynamodb put", map[string]string{"profile": profile, "region": region, "table": table, "item": item}, args)
	},
}

func init() {
	dynamodbCmd.AddCommand(dynamodbPutCmd)
	dynamodbPutCmd.Flags().String("table", "", "DynamoDB table name")
	dynamodbPutCmd.Flags().String("item", "", "Path to item to put")
	dynamodbPutCmd.MarkFlagRequired("item")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dynamodbPutCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dynamodbPutCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
