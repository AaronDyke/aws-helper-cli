package s3

import (
	"context"
	"log"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func ListBuckets(profile string) []string {
	s3Cmd := exec.Command("aws", "s3", "ls", "--profile", profile)
	out, err := s3Cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	buckets := strings.Split(string(out), "\n")
	buckets = buckets[:len(buckets)-1]
	for i, bucket := range buckets {
		buckets[i] = strings.Split(bucket, " ")[2]
	}

	return buckets
}

func ListBucketsSDK(profile string) []string {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))
	if err != nil {
		panic(err)
	}
	s3Client := s3.NewFromConfig(sdkConfig)
	result, err := s3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		panic(err)
	}
	var buckets []string
	for _, bucket := range result.Buckets {
		buckets = append(buckets, *bucket.Name)
	}
	return buckets
}
