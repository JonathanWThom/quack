package cmd

import (
	"fmt"
	"strings"

	"github.com/jonathanwthom/quack/storage"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new entry",
	Long: `Create a new entry like this:
quack new These are my deepest darkest secrets...`,
	Run: NewRunner,
}

// NewRunner wraps New for easier testing
func NewRunner(cmd *cobra.Command, args []string) {
	New(args...)
}

// New creates and stores a new message
func New(args ...string) {
	msg := strings.Join(args, " ")

	if len(msg) > 280 {
		Log(logFunc, "Message must be shorter than 280 characters.")
	}

	err := storage.Create(msg)
	if err != nil {
		Log(logFunc, "Failed to create entry.")
	}

	fmt.Println("Entry saved.")
}

func init() {
	rootCmd.AddCommand(newCmd)
}
