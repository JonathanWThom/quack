package cmd

import (
	"fmt"
	"strings"
	"testing"
)

type fakeStorage struct{}

func (s *fakeStorage) Create(msg string) error {
	fmt.Println(msg)

	return nil
}

func TestNew(t *testing.T) {
	store = new(fakeStorage)

	var tests = []struct {
		args     string
		expected string
	}{
		{
			args:     "valid entry",
			expected: successMsg,
		},
		{
			args: `
				morethan280charactersmorethan280charactersmorethan280characters
				morethan280charactersmorethan280charactersmorethan280characters
				morethan280charactersmorethan280charactersmorethan280characters
				morethan280charactersmorethan280charactersmorethan280characters
				morethan280charactersmorethan280charactersmorethan280characters
			`,
			expected: tooManyCharsError,
		},
	}

	for i := 0; i < len(tests); i++ {
		test := tests[i]
		args := strings.Split(test.args, " ")
		expected := test.expected
		actual := New(args...)

		if actual != expected {
			t.Errorf("cmd.New(%v) returned %s, expected %s", args, actual, expected)
		}
	}
}
