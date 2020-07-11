package storage

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestTransform(t *testing.T) {
	zone, _ := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Now().Location()).Zone()
	os.Setenv("QUACKWORD", "password")
	tests := []struct {
		entry    Entry
		verbose  bool
		search   string
		date     string
		expected string
	}{
		{
			entry: Entry{
				CreatedAt: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Now().Location()),
				Content:   "7ruS7L8Ksk8bHCtpWp1+OOJ0N9z92Xr5fFUJHARiTWwXpQwaJ6iBLQ==",
			},
			verbose:  false,
			search:   "",
			date:     "",
			expected: fmt.Sprintf("%s - %s %s\n%s", "November 10, 2009", "11:00 PM", zone, "Hello World!"),
		},
	}

	for i := 0; i < len(tests); i++ {
		test := tests[i]
		entry := test.entry
		verbose := test.verbose
		search := test.search
		date := test.date
		expected := test.expected

		actual, _ := entry.Transform(verbose, search, date)

		if actual != expected {
			t.Errorf(
				"entry.Transform(%v, %s, %s) with %v returned %s, expected %s",
				verbose,
				search,
				date,
				entry,
				actual,
				expected,
			)
		}
	}
}

func TestFilter(t *testing.T) {
	tests := []struct {
		entry    Entry
		search   string
		date     string
		expected bool
	}{
		{
			entry: Entry{
				DecryptedContent: "Foo",
				CreatedAt:        time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Now().Location()),
			},
			search:   "",
			date:     "",
			expected: true,
		},
		{
			entry: Entry{
				DecryptedContent: "Foo",
				CreatedAt:        time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Now().Location()),
			},
			search:   "Foo",
			date:     "",
			expected: true,
		},
		{
			entry: Entry{
				DecryptedContent: "Foo",
				CreatedAt:        time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Now().Location()),
			},
			search:   "Bar",
			date:     "",
			expected: false,
		},
		{
			entry: Entry{
				DecryptedContent: "Foo",
				CreatedAt:        time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Now().Location()),
			},
			search:   "",
			date:     "November 10, 2009",
			expected: true,
		},
		{
			entry: Entry{
				DecryptedContent: "Foo",
				CreatedAt:        time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Now().Location()),
			},
			search:   "",
			date:     "November 11, 2009",
			expected: false,
		},
	}

	for i := 0; i < len(tests); i++ {
		test := tests[i]
		entry := test.entry
		search := test.search
		date := test.date
		expected := test.expected

		_, actual := entry.Filter(search, date)
		if actual != expected {
			t.Errorf(
				"entry.Filter(%s, %s) with %v returned %v, expected %v",
				search,
				date,
				entry,
				actual,
				expected,
			)
		}
	}
}

func TestSetDecryptedContent(t *testing.T) {
}

func TestFormat(t *testing.T) {
	zone, _ := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Now().Location()).Zone()
	tests := []struct {
		entry    Entry
		verbose  bool
		expected string
	}{
		{
			entry: Entry{
				CreatedAt:        time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Now().Location()),
				DecryptedContent: "Oh hey there",
				Key:              "obfuscatedkey",
			},
			verbose:  false,
			expected: fmt.Sprintf("%s - %s %s\n%s", "November 10, 2009", "11:00 PM", zone, "Oh hey there"),
		},
		{
			entry: Entry{
				CreatedAt:        time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Now().Location()),
				DecryptedContent: "Oh hey there",
				Key:              "obfuscatedkey",
			},
			verbose:  true,
			expected: fmt.Sprintf("%s - %s %s - %s\n%s", "November 10, 2009", "11:00 PM", zone, "obfuscatedkey", "Oh hey there"),
		},
	}

	for i := 0; i < len(tests); i++ {
		test := tests[i]
		entry := test.entry
		verbose := test.verbose
		expected := test.expected

		actual, _ := entry.Format(verbose)

		if actual != expected {
			t.Errorf("entry.Format(%v) with %v returned %s, expected %s", verbose, entry, actual, expected)
		}
	}
}
