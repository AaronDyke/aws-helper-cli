package utils

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

func CommandString(cmd string, flags map[string]string, args []string) string {
	finishedCmd := fmt.Sprintf("aws-helper-cli %s", cmd)
	for key, value := range flags {
		finishedCmd = finishedCmd + fmt.Sprintf(" --%s %s", key, value)
	}
	for _, arg := range args {
		finishedCmd = finishedCmd + fmt.Sprintf(" %s", arg)
	}
	return finishedCmd
}

func PromptText(label string) string {
	defaultPrompt := promptui.Prompt{
		Label: label,
	}
	result, err := defaultPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		panic(err)
	}
	return result
}

func PromptItems(label string, items []string) string {
	itemPrompt := promptui.Select{
		Label: label,
		Items: items,
	}
	_, choice, err := itemPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		panic(err)
	}
	return choice
}

func PromptYesNo(label string) bool {
	defaultPrompt := promptui.Prompt{
		Label:   label + " (y/n)",
		Default: "n",
	}
	result, err := defaultPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		panic(err)
	}
	result = strings.ToLower(result)
	return result == "y" || result == "yes"
}
