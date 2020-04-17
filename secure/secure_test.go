package secure

import (
	"os"
	"testing"
)

func TestEncrypt(t *testing.T) {
	tests := []struct {
		input       string
		quackword   string
		expectedLen int
	}{
		{
			input:       "foo",
			quackword:   "exists",
			expectedLen: 44,
		},
		{
			input:       "foo",
			quackword:   "",
			expectedLen: 0,
		},
	}

	for i := 0; i < len(tests); i++ {
		test := tests[i]
		input := test.input
		expectedLen := test.expectedLen
		os.Setenv("QUACKWORD", test.quackword)

		actual, _ := Encrypt(input)
		actualLen := len(actual)

		if actualLen != expectedLen {
			t.Errorf("secure.Encrypt(%s) returned value of length %d, expected %d", input, actualLen, expectedLen)
		}
	}
}

func TestDecrypt(t *testing.T) {
	tests := []struct {
		expected  string
		quackword string
	}{
		{
			expected:  "foo",
			quackword: "exists",
		},
		{
			expected:  "",
			quackword: "",
		},
	}

	for i := 0; i < len(tests); i++ {
		test := tests[i]
		expected := test.expected
		os.Setenv("QUACKWORD", "exists")
		encrypted, _ := Encrypt(test.expected)
		os.Setenv("QUACKWORD", test.quackword)

		actual, _ := Decrypt(encrypted)

		if actual != expected {
			t.Errorf("Secure.Decrypt(%s) returned %s, expected %s", encrypted, actual, expected)
		}
	}
}
