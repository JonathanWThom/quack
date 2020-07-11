package cmd

import (
	"errors"
	"github.com/jonathanwthom/quack/secure"
	"github.com/jonathanwthom/quack/storage"
	"os"
	"testing"
)

func TestDelete(t *testing.T) {
	store = new(fakeStorage)
	os.Setenv("QUACKWORD", "password")

	tests := []struct {
		args               string
		expected           string
		readByKeyMock      string
		readByKeyErrorMock error
	}{
		{
			args:               "found-key",
			expected:           deleteSuccessMsg,
			readByKeyMock:      "found entry content",
			readByKeyErrorMock: nil,
		}, {
			args:               "not-found-key",
			expected:           unableToDeleteError,
			readByKeyMock:      "",
			readByKeyErrorMock: errors.New("can't find that"),
		},
	}

	for i := 0; i < len(tests); i++ {
		test := tests[i]
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
	}
}
