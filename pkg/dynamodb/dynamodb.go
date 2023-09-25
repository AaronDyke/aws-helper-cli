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
	TableName []string `json:"TableName"`
}

func ListTables(aws aws.Aws) []string {
	cmd := exec.Command("aws", "dynamodb", "list-tables", "--profile", aws.Profile, "--region", aws.Region)
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	ListTableResponse := ListTable{}
	json.Unmarshal(out, &ListTableResponse)
	return ListTableResponse.TableName
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
