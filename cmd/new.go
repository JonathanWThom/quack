package cmd

import (
	"fmt"
	"strings"

	"github.com/jonathanwthom/quack/secure"
	"github.com/spf13/cobra"
)

const (
	successMsg        = "Entry saved."
	tooManyCharsError = "Message must be shorter than 280 characters."
	storageError      = "Failed to create entry."
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new entry",
	Long: `
Create a new entry like this:
quack new These are my deepest darkest secrets...`,
	Run: NewRunner,
}

// NewRunner wraps New for easier testing
func NewRunner(cmd *cobra.Command, args []string) {
	result := New(args...)
	fmt.Println(result)
}

// New creates and stores a new message
func New(args ...string) string {
	msg := strings.Join(args, " ")

	if len(msg) > 280 {
		return tooManyCharsError
	}

	encrypted, err := secure.Encrypt(msg)
	if err != nil {
		return err.Error()
	}

	err = store.Create(encrypted)
	if err != nil {
		return storageError
	}

	return successMsg
}

func init() {
	rootCmd.AddCommand(newCmd)
}
