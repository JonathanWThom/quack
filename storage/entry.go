package storage

import (
	"fmt"
	"github.com/jonathanwthom/quack/secure"
	"strings"
	"time"
)

const layout = "Mon Jan 2 15:04:05 -0700 MST 2006"

// Entry stores entries with metadata
type Entry struct {
	CreatedAt        time.Time
	Content          string
	Key              string
	DecryptedContent string
}

// SetDecryptedContent decrypts an entry's content and sets the plain value on the object
func (entry *Entry) SetDecryptedContent() error {
	content, err := secure.Decrypt(entry.Content)
	if err != nil {
		return err
	}

	entry.DecryptedContent = content
	return nil
}

// Filter filters entries down by search term and date
func (entry *Entry) Filter(search, date string) (*Entry, bool) {
	if date != "" && entry.CreatedAt.Format("January 2, 2006") != date {
		return entry, false
	}

	if search != "" && !strings.Contains(strings.ToLower(entry.DecryptedContent), strings.ToLower(search)) {
		return entry, false
	}

	return entry, true
}

// Transform prepares and entry for display
func (entry *Entry) Transform(verbose bool, search string, date string) (string, error) {
	err := entry.SetDecryptedContent()
	if err != nil {
		return "", err
	}

	entry, ok := entry.Filter(search, date)
	if !ok {
		return "", err
	}

	return entry.Format(verbose)
}

// Format pretty-prints an entry
func (entry *Entry) Format(verbose bool) (string, error) {
	loc := time.Now().Location()
	formatted := entry.CreatedAt.In(loc).Format("January 2, 2006 - 3:04 PM MST")
	var result string
	if verbose {
		key := entry.Key
		result = fmt.Sprintf("%v - %s\n%s", formatted, key, entry.DecryptedContent)
	} else {
		result = fmt.Sprintf("%v\n%s", formatted, entry.DecryptedContent)
	}

	return result, nil
}
