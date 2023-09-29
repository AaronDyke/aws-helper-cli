package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func ListBuckets(profile string) []string {
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
