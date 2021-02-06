package cmd

import (
	"errors"
	"os"
	"testing"

	"github.com/jonathanwthom/quack/secure"
	"github.com/jonathanwthom/quack/storage"
)

func TestDelete(t *testing.T) {
	store = new(fakeStorage)
	os.Setenv("QUACKWORD", "password")

	tests := []struct {
		args               string
		expected           string
		readByKeyMock      string
		readByKeyErrorMock error
		description        string
	}{
		{
			args:               "found-key",
			expected:           deleteSuccessMsg,
			readByKeyMock:      "found entry content",
			readByKeyErrorMock: nil,
			description:        "when key to delete is found",
		}, {
			args:               "not-found-key",
			expected:           unableToDeleteError,
			readByKeyMock:      "",
			readByKeyErrorMock: errors.New("can't find that"),
			description:        "when key to delete is not found",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			expected := test.expected
			args := test.args
			encrypted, _ := secure.Encrypt(test.readByKeyMock)
			readByKeyMock = storage.Entry{
				Content: encrypted,
			}
			readByKeyErrorMock = test.readByKeyErrorMock

			actual := Delete(args)

			if actual != expected {
				t.Errorf("cmd.Delete(%v) returned %s, expected %s", args, actual, expected)
			}
		})
	}
}
