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
	"io"
	"log"
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

// Read will read the content of all messages from the cloud or local file.
func (s *Storage) Read() ([]string, error) {
	ctx := context.Background()

	if cloudConfigPresent() {
		return readFromCloud(ctx)
	}

	return readFromFiles(ctx)
}

func readFromCloud(ctx context.Context) ([]string, error) {
	// share this and context?
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("S3_BUCKET_REGION")),
	})
	if err != nil {
		return []string{}, err
	}

	bucket, err := s3blob.OpenBucket(ctx, sess, os.Getenv("S3_BUCKET_NAME"), nil)
	if err != nil {
		return []string{}, err
	}
	defer bucket.Close()

	return readFromBucket(ctx, bucket)
}

func readFromFiles(ctx context.Context) ([]string, error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		return []string{}, err
	}
	dir := homeDir + "/quack"
	if err := os.MkdirAll(dir, 0777); err != nil {
		return []string{}, err
	}

	bucket, err := fileblob.OpenBucket(dir, nil)
	if err != nil {
		return []string{}, err
	}
	defer bucket.Close()

	return readFromBucket(ctx, bucket)
}

func readFromBucket(ctx context.Context, bucket *blob.Bucket) ([]string, error) {
	iter := bucket.List(nil)
	var results []string
	for {
		obj, err := iter.Next(ctx)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err) // this should return error and exit
		}
		res, err := bucket.ReadAll(ctx, obj.Key)
		if err != nil {
			log.Fatal(err) // this should return error and exit
		}

		results = append(results, string(res))
	}

	return results, nil
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
