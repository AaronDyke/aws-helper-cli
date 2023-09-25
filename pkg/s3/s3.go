package s3

import (
	"log"
	"os/exec"
	"strings"
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
