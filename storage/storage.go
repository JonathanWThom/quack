package storage

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/jonathanwthom/quack/secure"
	"github.com/mitchellh/go-homedir"
	"gocloud.dev/blob"
	"gocloud.dev/blob/fileblob"
	"gocloud.dev/blob/s3blob"
	"io"
	"os"
	"strings"
	"time"
)

type Storage struct{}
type Entry struct {
	ModTime time.Time
	Content string
	Key     string
}

func (entry *Entry) Format(verbose bool, search string, date string) (string, error) {
	content, err := secure.Decrypt(entry.Content)

	if err != nil {
		return "", err
	}

	if date != "" && entry.ModTime.Format("January 2, 2006") != date {
		return "", nil
	}

	if search != "" && !strings.Contains(strings.ToLower(content), strings.ToLower(search)) {
		return "", nil
	}

	loc := time.Now().Location()
	formatted := entry.ModTime.In(loc).Format("January 2, 2006 - 3:04 PM MST")
	var result string
	if verbose == true {
		key := entry.Key
		result = fmt.Sprintf("%v - %s\n%s", formatted, key, content)
	} else {
		result = fmt.Sprintf("%v\n%s", formatted, content)
	}

	return result, nil
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
