package aws

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
)

type Aws struct {
	Profile string
	Region  string
}

type AwsRegion struct {
	RegionName      string `json:"RegionName"`
	RegionOptStatus string `json:"RegionOptStatus"`
}
type AwsRegions struct {
	Regions []AwsRegion `json:"Regions"`
}

func ListProfiles() []string {
	cmd := exec.Command("aws", "configure", "list-profiles", "--output", "json")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	profiles := strings.Split(string(out), "\n")
	profiles = profiles[:len(profiles)-1]
	return profiles
}

func ProfileExists(profile string) bool {
	profiles := ListProfiles()
	for _, p := range profiles {
		if p == profile {
			return true
		}
	}
	return false
}

func PromptProfile() string {
	profilePrompt := promptui.Select{
		Label: "Select Profile",
		Items: ListProfiles(),
	}
	_, profile, err := profilePrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		panic(err)
	}
	return profile
}

func ListRegions(profile string) []string {
	cmd := exec.Command("aws", "account", "list-regions", "--output", "json", "--profile", profile)
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	awsRegions := AwsRegions{}
	json.Unmarshal(out, &awsRegions)

	enabledRegions := []string{}
	for _, region := range awsRegions.Regions {
		if region.RegionOptStatus == "ENABLED" || region.RegionOptStatus == "ENABLED_BY_DEFAULT" {
			enabledRegions = append(enabledRegions, region.RegionName)
		}
	}

	return enabledRegions
}

func DefaultRegion(profile string) string {
	cmd := exec.Command("aws", "configure", "get", "region", "--profile", profile)
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(string(out))
}

func RegionExists(profile string, region string) bool {
	regions := ListRegions(profile)
	for _, r := range regions {
		if r == region {
			return true
		}
	}
	return false
}

func PromptRegion(profile string) string {
	defaultPrompt := promptui.Prompt{
		Label:   "The default region is " + DefaultRegion(profile) + ". Do you want to use this region? (y/n)",
		Default: "n",
	}
	result, err := defaultPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		panic(err)
	}
	result = strings.ToLower(result)
	if result == "y" || result == "yes" {
		return DefaultRegion(profile)
	}

	regionPrompt := promptui.Select{
		Label: "Select Region",
		Items: ListRegions(profile),
	}
	_, region, err := regionPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		panic(err)
	}
	return region
}
