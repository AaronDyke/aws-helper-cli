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
