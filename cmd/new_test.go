package cmd

import (
	"os"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	store = new(fakeStorage)
	os.Setenv("QUACKWORD", "password")

	var tests = []struct {
		args        string
		expected    string
		description string
	}{
		{
			args:        "valid entry",
			expected:    successMsg,
			description: "when new entry is valid",
		},
		{
			args: `
				morethan280charactersmorethan280charactersmorethan280characters
				morethan280charactersmorethan280charactersmorethan280characters
				morethan280charactersmorethan280charactersmorethan280characters
				morethan280charactersmorethan280charactersmorethan280characters
				morethan280charactersmorethan280charactersmorethan280characters
			`,
			expected:    tooManyCharsError,
			description: "when new entry has too many characters",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			args := strings.Split(test.args, " ")
			expected := test.expected
			actual := New(args...)

			if actual != expected {
				t.Errorf("cmd.New(%v) returned %s, expected %s", args, actual, expected)
			}
		})
	}
}
