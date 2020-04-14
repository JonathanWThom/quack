package cmd

import "testing"
import "fmt"

// stub out storage too
// stub out printing

func fakeLog(...interface{}) {
	fmt.Println("it worked")
}

type fakeStorage struct{}

func (s *fakeStorage) Create(msg string) error {
	fmt.Println(msg)

	return nil
}

func TestNew(t *testing.T) {
	logFunc = fakeLog
	store = new(fakeStorage)

	var tests = []struct {
		args   []string
		output string
	}{
		{
			args:   []string{"valid", "entry"},
			output: "Entry saved.",
		},
	}

	for i := 0; i < len(tests); i++ {
		test := tests[i]
		args := test.args
		output := test.output
		var result string
		New(args...)

		if result != output {
			t.Errorf("cmd.New(%v) returned %s, expected %s", args, result, output)
		}
	}
}
