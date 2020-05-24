package storage

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/mitchellh/go-homedir"
	"gocloud.dev/blob"
	"gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/gcsblob"
	"gocloud.dev/blob/s3blob"
	"io"
	"os"
	"time"
)

var amazonVars = []string{
	"S3_BUCKET_REGION",
	"S3_BUCKET_NAME",
	"AWS_ACCESS_KEY_ID",
	"AWS_SECRET_ACCESS_KEY",
}

var googleVars = []string{
	"GOOGLE_APPLICATION_CREDENTIALS",
	"GOOGLE_BUCKET_NAME",
}

const amazon = "amazon"
const google = "google"

type Storage struct{}

type cloudEnv struct {
	name string
}

func (c *cloudEnv) amazon() bool {
	return c.name == amazon
}

func (c *cloudEnv) google() bool {
	return c.name == google
}

var cloud cloudEnv

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

// Read will read a single message from the cloud or local file, selected by
// key.
func (s *Storage) ReadByKey(key string) (Entry, error) {
	ctx := context.Background()

	if cloudConfigPresent() {
		return readByKeyFromCloud(ctx, key)
	}

	return readByKeyFromFiles(ctx, key)
}

func readByKeyFromCloud(ctx context.Context, key string) (Entry, error) {
	bucket, err := openCloudBucket(ctx)
	defer bucket.Close()
	if err != nil {
		return Entry{}, err
	}

	return readFromBucketByKey(ctx, bucket, key)
}

func readFromBucketByKey(ctx context.Context, bucket *blob.Bucket, key string) (Entry, error) {
	res, err := bucket.ReadAll(ctx, key)
	if err != nil {
		return Entry{}, err
	}

	entry := Entry{
		// would be nice to get modtime
		Content: string(res),
		Key:     key,
	}

	return entry, nil
}

func readByKeyFromFiles(ctx context.Context, key string) (Entry, error) {
	bucket, err := openFileBucket()
	defer bucket.Close()
	if err != nil {
		return Entry{}, err
	}

	return readFromBucketByKey(ctx, bucket, key)
}

// Read will delete an entry by its unique key, from either the cloud or a local
// file.
func (s *Storage) Delete(key string) error {
	ctx := context.Background()

	if cloudConfigPresent() {
		return deleteFromCloud(key, ctx)
	}

	return deleteFromFile(key, ctx)
}

func deleteFromCloud(key string, ctx context.Context) error {
	bucket, err := openCloudBucket(ctx)
	defer bucket.Close()
	if err != nil {
		return err
	}

	err = bucket.Delete(ctx, key)
	if err != nil {
		return err
	}

	return nil
}

func deleteFromFile(key string, ctx context.Context) error {
	bucket, err := openFileBucket()
	defer bucket.Close()
	if err != nil {
		return err
	}

	err = bucket.Delete(ctx, key)
	if err != nil {
		return err
	}

	return nil
}

func openCloudBucket(ctx context.Context) (*blob.Bucket, error) {
	if cloud.amazon() {
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(os.Getenv("S3_BUCKET_REGION")),
		})
		if err != nil {
			return new(blob.Bucket), err
		}

		return s3blob.OpenBucket(ctx, sess, os.Getenv("S3_BUCKET_NAME"), nil)
	}

	return blob.OpenBucket(ctx, "gs://"+os.Getenv("GOOGLE_BUCKET_NAME"))
}

func openFileBucket() (*blob.Bucket, error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		return new(blob.Bucket), err
	}
	dir := homeDir + "/.quack"
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
			Key:     obj.Key,
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
	sum := sha256.Sum256([]byte(time.Now().String()))
	key := fmt.Sprintf("%x", sum)
	w, err := bucket.NewWriter(ctx, key, nil)
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

func allVarsPresent(params []string) bool {
	result := true
	for i := 0; i < len(params); i++ {
		if os.Getenv(params[i]) == "" {
			result = false
			break
		}
	}

	return result
}

func (c *cloudEnv) setName() {
	if allVarsPresent(amazonVars) {
		c.name = amazon
	} else if allVarsPresent(googleVars) {
		c.name = google
	}
}

//would be nice to memoize?
func cloudConfigPresent() bool {
	if cloud.name == "" {
		cloud.setName()
	}

	return cloud.name != ""
}
