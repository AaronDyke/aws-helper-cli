package utils

import "fmt"

func PrintRunCommandAgain(cmd string, flags map[string]string, args []string) {
	fmt.Println("To run this exact command again, use the following:")
	finishedCmd := fmt.Sprintf("aws-helper-cli %s", cmd)
	for key, value := range flags {
		finishedCmd = finishedCmd + fmt.Sprintf(" --%s %s", key, value)
	}
	for _, arg := range args {
		finishedCmd = finishedCmd + fmt.Sprintf(" %s", arg)
	}
	fmt.Println(finishedCmd)
}
