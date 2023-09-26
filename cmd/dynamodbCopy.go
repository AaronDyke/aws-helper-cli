/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/AaronDyke/aws-helper-cli/pkg/aws"
	"github.com/AaronDyke/aws-helper-cli/pkg/dynamodb"
	"github.com/spf13/cobra"
)

// dynamodbCopyCmd represents the dynamodbCopy command
var dynamodbCopyCmd = &cobra.Command{
	Use:   "copy",
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
		fromTable := cmd.Flag("from").Value.String()
		if fromTable == "" {
			fromTable = dynamodb.PromptTables(aws, "Select DynamoDB Table to copy from", []string{})
		} else if !dynamodb.TableExists(aws, fromTable) {
			fmt.Println("Table does not exist")
			return
		}

		toTable := cmd.Flag("to").Value.String()
		if toTable == "" {
			toTable = dynamodb.PromptTables(aws, "Select DynamoDB Table to copy to", []string{fromTable})
		} else if !dynamodb.TableExists(aws, toTable) {
			fmt.Println("Table does not exist")
			return
		}

		partitionKey := cmd.Flag("partition-key").Value.String()
		if partitionKey == "DEFAULT" {
			partitionKey = dynamodb.PromptPartitionKey(aws, fromTable)
		}

		sortKey := cmd.Flag("sort-key").Value.String()
		sortKeyBeginsWith := cmd.Flag("sort-key-begins-with").Value.String()

		if partitionKey == "" {
			fmt.Println("Please specify a partition key")
			return
		} else if sortKey != "" {
			dynamodb.CopyItem(aws, fromTable, toTable, partitionKey, sortKey)
		} else {
			dynamodb.CopyItems(aws, fromTable, toTable, partitionKey, sortKeyBeginsWith)
		}

		if cmd.Flag("quiet").Value.String() == "true" {
			return
		} else {
			fmt.Println("\nTo run this exact command again, use the following:")
			finishedCmd := fmt.Sprintf("aws-helper-cli dynamodb copy --profile %s --region %s --from %s --to %s --partition-key %s", profile, region, fromTable, toTable, partitionKey)
			if sortKey != "" {
				finishedCmd = finishedCmd + fmt.Sprintf(" --sort-key %s", sortKey)
			} else if sortKeyBeginsWith != "" {
				finishedCmd = finishedCmd + fmt.Sprintf(" --sort-key-begins-with %s", sortKeyBeginsWith)
			}
			fmt.Println(finishedCmd)
		}
	},
}

func init() {
	dynamodbCmd.AddCommand(dynamodbCopyCmd)
	dynamodbCopyCmd.Flags().String("from", "", "DynamoDB table name to copy from")
	dynamodbCopyCmd.Flags().String("to", "", "DynamoDB table name to copy to")
	dynamodbCopyCmd.Flags().String("partition-key", "", "Partition key to use for copy")
	dynamodbCopyCmd.Flags().String("sort-key", "", "Sort key to use for copy")
	dynamodbCopyCmd.Flags().String("sort-key-begins-with", "", "Sort key begins with value to use for copy")
	dynamodbCopyCmd.MarkFlagsMutuallyExclusive("sort-key", "sort-key-begins-with")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dynamodbCopyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dynamodbCopyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
