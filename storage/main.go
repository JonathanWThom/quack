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
	"os"
	"time"
)

type Storage struct{}
type Entry struct {
	ModTime time.Time
	Content string
}

// Create will save a message to the cloud, or a local file.
func (s *Storage) Create(msg string) error {
	ctx := context.Background()

	if cloudConfigPresent() {
		return writeToCloud(ctx, msg)
	}

	return writeToFile(ctx, msg)
}

// Read will read the content of all messages from the cloud or local file.
func (s *Storage) Read() ([]Entry, error) {
	ctx := context.Background()

	if cloudConfigPresent() {
		return readFromCloud(ctx)
	}

	return readFromFiles(ctx)
}

func openCloudBucket(ctx context.Context) (*blob.Bucket, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("S3_BUCKET_REGION")),
	})
	if err != nil {
		return new(blob.Bucket), err
	}

	bucket, err := s3blob.OpenBucket(ctx, sess, os.Getenv("S3_BUCKET_NAME"), nil)
	if err != nil {
		return bucket, err
	}

	return bucket, nil
}

func openFileBucket() (*blob.Bucket, error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		return new(blob.Bucket), err
	}
	dir := homeDir + "/quack"
	if err := os.MkdirAll(dir, 0777); err != nil {
		return new(blob.Bucket), err
	}

	bucket, err := fileblob.OpenBucket(dir, nil)
	if err != nil {
		return bucket, err
	}

	return bucket, nil
}

func readFromCloud(ctx context.Context) ([]Entry, error) {
	bucket, err := openCloudBucket(ctx)
	defer bucket.Close()
	if err != nil {
		return []Entry{}, err
	}

	return readFromBucket(ctx, bucket)
}

func readFromFiles(ctx context.Context) ([]Entry, error) {
	bucket, err := openFileBucket()
	defer bucket.Close()
	if err != nil {
		return []Entry{}, err
	}

	return readFromBucket(ctx, bucket)
}

func readFromBucket(ctx context.Context, bucket *blob.Bucket) ([]Entry, error) {
	var entries []Entry
	iter := bucket.List(nil)

	for {
		obj, err := iter.Next(ctx)
		if err == io.EOF {
			break
		}
		if err != nil {
			return []Entry{}, err
		}
		res, err := bucket.ReadAll(ctx, obj.Key)
		if err != nil {
			return []Entry{}, err
		}

		entry := Entry{
			ModTime: obj.ModTime,
			Content: string(res),
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func writeToCloud(ctx context.Context, msg string) error {
	bucket, err := openCloudBucket(ctx)
	defer bucket.Close()
	if err != nil {
		return err
	}

	return writeToBucket(ctx, msg, bucket)
}

func writeToFile(ctx context.Context, msg string) error {
	bucket, err := openFileBucket()
	defer bucket.Close()
	if err != nil {
		return err
	}

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
