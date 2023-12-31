package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/AaronDyke/aws-helper-cli/pkg/utils"
	"github.com/aws/aws-sdk-go-v2/config"
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

func IsCliInstalled() bool {
	cmd := exec.Command("aws", "--version")
	err := cmd.Run()
	return err == nil
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
	profileChoice := utils.PromptItems("Select Profile", ListProfiles())
	return profileChoice
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
	if utils.PromptYesNo(fmt.Sprintf("Use default region (%s)?", DefaultRegion(profile))) {
		return DefaultRegion(profile)
	}

	regionChoice := utils.PromptItems("Select Region", ListRegions(profile))
	return regionChoice
}

func ConfigSDK(aws Aws) config.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(aws.Profile), config.WithRegion(aws.Region))
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	return cfg
}
