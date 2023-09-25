package dynamodb

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

	"github.com/AaronDyke/aws-helper-cli/pkg/aws"
	"github.com/manifoldco/promptui"
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
	cmd := exec.Command("aws", "dynamodb", "list-tables", "--profile", aws.Profile, "--region", aws.Region)
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	ListTableResponse := ListTable{}
	json.Unmarshal(out, &ListTableResponse)
	return ListTableResponse.TableNames
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

func PromptTables(aws aws.Aws) string {
	tables := ListTables(aws)
	if len(tables) == 0 {
		fmt.Println("No tables found")
		return ""
	}

	tablePrompt := promptui.Select{
		Label: "Select DynamoDB Table",
		Items: ListTables(aws),
	}
	_, table, err := tablePrompt.Run()
	if err != nil {
		log.Fatal(err)
	}
	return table
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
