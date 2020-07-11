package cmd

import (
	"github.com/jonathanwthom/quack/storage"
	"os"
	"testing"
	"time"
)

func TestQuackword(t *testing.T) {
	store = new(fakeStorage)
	os.Setenv("QUACKWORD", "password")

	tests := []struct {
		entries  []storage.Entry
		args     string
		expected string
	}{
		{
			entries: []storage.Entry{
				{
					Content:   "7ruS7L8Ksk8bHCtpWp1+OOJ0N9z92Xr5fFUJHARiTWwXpQwaJ6iBLQ==",
					Key:       "key",
					CreatedAt: time.Now(),
				},
			},
			args:     "new quackword",
			expected: updateSuccess,
		},
	}

	for i := 0; i < len(tests); i++ {
		test := tests[i]
		args := test.args
		expected := test.expected
		entriesMock = test.entries

		actual := Quackword(args)

		if actual != expected {
			t.Errorf("cmd.Quackword(%v) returned %s, expected %s", args, actual, expected)
		}
	}
}
