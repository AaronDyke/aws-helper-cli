package dynamodb

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

	"github.com/AaronDyke/aws-helper-cli/pkg/aws"
	"github.com/AaronDyke/aws-helper-cli/pkg/utils"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type ListTable struct {
	TableNames []string `json:"TableNames"`
}

type DescribeTableResponse struct {
	Table Table `json:"Table"`
}

type Table struct {
	AttributeDefinitions  []AttributeDefinition `json:"AttributeDefinitions"`
	CreationDateTime      int64                 `json:"CreationDateTime"`
	ItemCount             int                   `json:"ItemCount"`
	KeySchema             []KeySchemaElement    `json:"KeySchema"`
	ProvisionedThroughput ProvisionedThroughput `json:"ProvisionedThroughput"`
	TableName             string                `json:"TableName"`
	TableSizeBytes        int                   `json:"TableSizeBytes"`
	TableStatus           string                `json:"TableStatus"`
}

type AttributeDefinition struct {
	AttributeName string `json:"AttributeName"`
	AttributeType string `json:"AttributeType"`
}

type KeySchemaElement struct {
	AttributeName string `json:"AttributeName"`
	KeyType       string `json:"KeyType"`
}

type ProvisionedThroughput struct {
	LastDecreaseDateTime   int64 `json:"LastDecreaseDateTime"`
	LastIncreaseDateTime   int64 `json:"LastIncreaseDateTime"`
	NumberOfDecreasesToday int   `json:"NumberOfDecreasesToday"`
	ReadCapacityUnits      int   `json:"ReadCapacityUnits"`
	WriteCapacityUnits     int   `json:"WriteCapacityUnits"`
}

func ListTables(aws aws.Aws) []string {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(aws.Profile), config.WithRegion(aws.Region))
	if err != nil {
		panic(err)
	}
	dynamoClient := dynamodb.NewFromConfig(sdkConfig)
	result, err := dynamoClient.ListTables(context.TODO(), &dynamodb.ListTablesInput{})
	if err != nil {
		panic(err)
	}
	tables := make([]string, 0, len(result.TableNames))
	tables = append(tables, result.TableNames...)
	return tables
}

func TableExists(aws aws.Aws, table string) bool {
	tables := ListTables(aws)
	for _, t := range tables {
		if t == table {
			return true
		}
	}
	return false
}

func PromptTables(aws aws.Aws, label string, excludeTables []string) string {
	tables := ListTables(aws)
	if len(tables) == 0 {
		fmt.Println("No tables found")
		return ""
	}
	// remove all tables in excludeTables
	for _, excludeTable := range excludeTables {
		for i, table := range tables {
			if table == excludeTable {
				tables = append(tables[:i], tables[i+1:]...)
			}
		}
	}

	if len(tables) == 0 {
		fmt.Println("No tables found")
		return ""
	}

	tableChoice := utils.PromptItems(label, tables)
	return tableChoice
}

func PutItem(aws aws.Aws, table string, pathToItem string) {
	// fmt.Println("aws ", "dynamodb ", "put-item ", "--table-name ", table, " --item ", fmt.Sprintf("file://%s", pathToItem), " --profile ", aws.Profile, " --region ", aws.Region)
	cmd := exec.Command("aws", "dynamodb", "put-item", "--table-name", table, "--item", fmt.Sprintf("file://%s", pathToItem), "--profile", aws.Profile, "--region", aws.Region)
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))
}

func DescribeTable(aws aws.Aws, table string) DescribeTableResponse {
	cmd := exec.Command("aws", "dynamodb", "describe-table", "--table-name", table, "--profile", aws.Profile, "--region", aws.Region)
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	DescribeTableResponse := DescribeTableResponse{}
	json.Unmarshal(out, &DescribeTableResponse)
	return DescribeTableResponse
}

func TableKeysFromTableDescription(tableDescription DescribeTableResponse) (string, string) {
	partitionKey := ""
	sortKey := ""
	for _, key := range tableDescription.Table.KeySchema {
		if partitionKey != "" && sortKey != "" {
			break
		}
		if key.KeyType == "HASH" {
			partitionKey = key.AttributeName
		} else if key.KeyType == "RANGE" {
			sortKey = key.AttributeName
		}
	}
	return partitionKey, sortKey
}

func TableKeys(aws aws.Aws, table string) (string, string) {
	tableDescription := DescribeTable(aws, table)
	partitionKey := ""
	sortKey := ""
	for _, key := range tableDescription.Table.KeySchema {
		if partitionKey != "" && sortKey != "" {
			break
		}
		if key.KeyType == "HASH" {
			partitionKey = key.AttributeName
		} else if key.KeyType == "RANGE" {
			sortKey = key.AttributeName
		}
	}
	return partitionKey, sortKey
}

func TablePartitionKey(aws aws.Aws, table string) string {
	tableDescription := DescribeTable(aws, table)
	for _, key := range tableDescription.Table.KeySchema {
		if key.KeyType == "HASH" {
			return key.AttributeName
		}
	}
	return ""
}

func TableSortKey(aws aws.Aws, table string) string {
	tableDescription := DescribeTable(aws, table)
	for _, key := range tableDescription.Table.KeySchema {
		if key.KeyType == "RANGE" {
			return key.AttributeName
		}
	}
	return ""
}

func PromptKey(keyName string) string {
	key := utils.PromptText(fmt.Sprintf("Query where %s = ", keyName))
	return key
}

func PromptPartitionKey(aws aws.Aws, table string) string {
	partitionKeyName := TablePartitionKey(aws, table)
	if partitionKeyName == "" {
		fmt.Println("No partition key found")
		return ""
	}

	return PromptKey(partitionKeyName)
}

func PromptSortKey(aws aws.Aws, table string) string {
	sortKeyName := TableSortKey(aws, table)
	if sortKeyName == "" {
		fmt.Println("No sort key found")
		return ""
	}

	return PromptKey(sortKeyName)
}

func CopyItems(aws aws.Aws, fromTable string, toTable string, partitionKey string, sortKeyBeginsWith string) {
	fmt.Println("Copy Items ", " --profile ", aws.Profile, " --region ", aws.Region, " --from-table ", fromTable, " --to-table ", toTable, " --partition-key ", partitionKey, " --sort-key-begins-with ", sortKeyBeginsWith)
}

func CopyItem(aws aws.Aws, fromTable string, toTable string, partitionKey string, sortKey string) {
	fmt.Println("Copy Item ", " --profile ", aws.Profile, " --region ", aws.Region, " --from-table ", fromTable, " --to-table ", toTable, " --partition-key ", partitionKey, " --sort-key ", sortKey)
}
