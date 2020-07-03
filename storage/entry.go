package storage

import (
	"fmt"
	"github.com/jonathanwthom/quack/secure"
	"strings"
	"time"
)

type Entry struct {
	ModTime          time.Time
	Content          string
	Key              string
	DecryptedContent string
}

func (entry *Entry) SetDecryptedContent() error {
	content, err := secure.Decrypt(entry.Content)
	if err != nil {
		return err
	}

	entry.DecryptedContent = content
	return nil
}

func (entry *Entry) Filter(search, date string) (*Entry, bool) {
	if date != "" && entry.ModTime.Format("January 2, 2006") != date {
		return entry, false
	}

	if search != "" && !strings.Contains(strings.ToLower(entry.DecryptedContent), strings.ToLower(search)) {
		return entry, false
	}

	return entry, true
}

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

func (entry *Entry) Format(verbose bool) (string, error) {
	loc := time.Now().Location()
	formatted := entry.ModTime.In(loc).Format("January 2, 2006 - 3:04 PM MST")
	var result string
	if verbose == true {
		key := entry.Key
		result = fmt.Sprintf("%v - %s\n%s", formatted, key, entry.DecryptedContent)
	} else {
		result = fmt.Sprintf("%v\n%s", formatted, entry.DecryptedContent)
	}

	return result, nil
}
