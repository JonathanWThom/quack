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
	// gcsblob is needed but only called indirectly
	_ "gocloud.dev/blob/gcsblob"
	"gocloud.dev/blob/s3blob"
	"io"
	"os"
	"time"
)

var amazonVars = []string{
	"QUACK_S3_BUCKET_REGION",
	"QUACK_S3_BUCKET_NAME",
	"QUACK_AWS_ACCESS_KEY_ID",
	"QUACK_AWS_SECRET_ACCESS_KEY",
}

var googleVars = []string{
	"QUACK_GOOGLE_APPLICATION_CREDENTIALS",
	"QUACK_GOOGLE_BUCKET_NAME",
}

const amazon = "amazon"
const google = "google"

// Storage implement all CRUD methods
type Storage struct{}

type cloudEnv struct {
	name string
}

func (c *cloudEnv) amazon() bool {
	return c.name == amazon
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

// Update rewrites and entry in storage
func (s *Storage) Update(e Entry) error {
	ctx := context.Background()

	if cloudConfigPresent() {
		return updateToCloud(ctx, e)
	}

	return updateToFile(ctx, e)
}

// Read will read the content of all messages from the cloud or local file.
func (s *Storage) Read() ([]Entry, error) {
	ctx := context.Background()

	if cloudConfigPresent() {
		return readFromCloud(ctx)
	}

	return readFromFiles(ctx)
}

// ReadByKey will read a single message from the cloud or local file, selected by
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
	if err != nil {
		return Entry{}, err
	}
	defer bucket.Close()

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
	if err != nil {
		return Entry{}, err
	}
	defer bucket.Close()

	return readFromBucketByKey(ctx, bucket, key)
}

// Delete will delete an entry by its unique key, from either the cloud or a local
// file.
func (s *Storage) Delete(key string) error {
	ctx := context.Background()

	if cloudConfigPresent() {
		return deleteFromCloud(ctx, key)
	}

	return deleteFromFile(ctx, key)
}

func deleteFromCloud(ctx context.Context, key string) error {
	bucket, err := openCloudBucket(ctx)
	if err != nil {
		return err
	}
	defer bucket.Close()

	err = bucket.Delete(ctx, key)
	if err != nil {
		return err
	}

	return nil
}

func deleteFromFile(ctx context.Context, key string) error {
	bucket, err := openFileBucket()
	if err != nil {
		return err
	}
	defer bucket.Close()

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
	if err != nil {
		return []Entry{}, err
	}
	defer bucket.Close()

	return readFromBucket(ctx, bucket)
}

func readFromFiles(ctx context.Context) ([]Entry, error) {
	bucket, err := openFileBucket()
	if err != nil {
		return []Entry{}, err
	}
	defer bucket.Close()

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

		attr, err := bucket.Attributes(ctx, obj.Key)
		if err != nil {
			return []Entry{}, err
		}

		rawCreatedAt := attr.Metadata["createdat"]
		createdAt, err := time.Parse(layout, rawCreatedAt)

		if err != nil {
			return []Entry{}, err
		}

		entry := Entry{
			CreatedAt: createdAt,
			Content:   string(res),
			Key:       obj.Key,
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func writeToCloud(ctx context.Context, msg string) error {
	bucket, err := openCloudBucket(ctx)
	if err != nil {
		return err
	}
	defer bucket.Close()

	return writeToBucket(ctx, msg, bucket)
}

func updateToCloud(ctx context.Context, e Entry) error {
	bucket, err := openCloudBucket(ctx)
	if err != nil {
		return err
	}
	defer bucket.Close()

	return updateToBucket(ctx, e, bucket)
}

func writeToFile(ctx context.Context, msg string) error {
	bucket, err := openFileBucket()
	if err != nil {
		return err
	}
	defer bucket.Close()

	return writeToBucket(ctx, msg, bucket)
}

func updateToFile(ctx context.Context, e Entry) error {
	bucket, err := openFileBucket()
	if err != nil {
		return err
	}
	defer bucket.Close()

	return updateToBucket(ctx, e, bucket)
}

func updateToBucket(ctx context.Context, e Entry, bucket *blob.Bucket) error {
	metadata := map[string]string{"createdAt": e.CreatedAt.Format(layout)}
	options := blob.WriterOptions{Metadata: metadata}
	w, err := bucket.NewWriter(ctx, e.Key, &options)
	if err != nil {
		return err
	}
	_, writeErr := fmt.Fprintln(w, e.Content)
	closeErr := w.Close()
	if writeErr != nil {
		return writeErr
	}
	if closeErr != nil {
		return closeErr
	}

	return nil
}

func writeToBucket(ctx context.Context, msg string, bucket *blob.Bucket) error {
	sum := sha256.Sum256([]byte(time.Now().String()))
	key := fmt.Sprintf("%x", sum)
	metadata := map[string]string{"createdAt": time.Now().Format(layout)}
	options := blob.WriterOptions{Metadata: metadata}
	w, err := bucket.NewWriter(ctx, key, &options)
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
