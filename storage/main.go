package storage

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/mitchellh/go-homedir"
	"gocloud.dev/blob"
	"gocloud.dev/blob/fileblob"
	"gocloud.dev/blob/s3blob"
	"os"
	"time"
)

type Storage struct{}

// Create will save a message to the cloud, or a local file.
func (s *Storage) Create(msg string) error {
	ctx := context.Background()

	if cloudConfigPresent() {
		return writeToCloud(ctx, msg)
	}

	return writeToFile(ctx, msg)
}

func writeToCloud(ctx context.Context, msg string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("S3_BUCKET_REGION")),
	})
	if err != nil {
		return err
	}

	bucket, err := s3blob.OpenBucket(ctx, sess, os.Getenv("S3_BUCKET_NAME"), nil)
	if err != nil {
		return err
	}
	defer bucket.Close()

	return writeToBucket(ctx, msg, bucket)
}

func writeToFile(ctx context.Context, msg string) error {
	homeDir, err := homedir.Dir()
	if err != nil {
		return err
	}
	dir := homeDir + "/quack"
	if err := os.MkdirAll(dir, 0777); err != nil {
		return err
	}

	bucket, err := fileblob.OpenBucket(dir, nil)
	if err != nil {
		return err
	}
	defer bucket.Close()

	return writeToBucket(ctx, msg, bucket)
}

func writeToBucket(ctx context.Context, msg string, bucket *blob.Bucket) error {
	w, err := bucket.NewWriter(ctx, time.Now().String(), nil)
	if err != nil {
		return err
	}
	_, writeErr := fmt.Fprintln(w, msg)
	closeErr := w.Close()
	if writeErr != nil {
		return writeErr
	}
	if closeErr != nil {
		return closeErr
	}

	return nil
}

func cloudConfigPresent() bool {
	params := []string{
		"S3_BUCKET_REGION",
		"S3_BUCKET_NAME",
		"AWS_ACCESS_KEY_ID",
		"AWS_SECRET_ACCESS_KEY",
	}

	result := true
	for i := 0; i < len(params); i++ {
		if os.Getenv(params[i]) == "" {
			result = false
			break
		}
	}

	return result
}
