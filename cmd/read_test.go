package cmd

import (
	"errors"
	"fmt"
	"github.com/jonathanwthom/quack/storage"
	"os"
	"testing"
	"time"
)

func TestRead(t *testing.T) {
	store = new(fakeStorage)
	os.Setenv("QUACKWORD", "password")

	tests := []struct {
		entries  []storage.Entry
		err      error
		expected string
	}{
		{
			entries:  []storage.Entry{},
			err:      nil,
			expected: "",
		},
		{
			entries: []storage.Entry{
				{
					ModTime: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
					Content: "7ruS7L8Ksk8bHCtpWp1+OOJ0N9z92Xr5fFUJHARiTWwXpQwaJ6iBLQ==",
				},
			},
			err:      nil,
			expected: fmt.Sprintf("%s\n%s", "2009-11-10 23:00:00 +0000 UTC:", "Hello World!"),
		},
		{
			entries:  []storage.Entry{},
			err:      errors.New("AWS is down, time to panic"),
			expected: unableToReadError,
		},
	}

	for i := 0; i < len(tests); i++ {
		test := tests[i]
		expected := test.expected
		entriesMock = test.entries
		errorMock = test.err
		actual := Read()

		if actual != expected {
			t.Errorf("cmd.Read() returned %s, expected %s", actual, expected)
		}
	}
}
