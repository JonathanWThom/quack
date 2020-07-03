package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	newQuackwordError   = "Please specify a new QUACKWORD."
	multiQuackwordError = "Please enter multi-word QUACKWORD in quotations."
)

// quackwordCmd represents the quackword command
var quackwordCmd = &cobra.Command{
	Use:   "quackword",
	Short: "Reset your QUACKWORD",
	Long: `
Reset your QUACKWORD like this:
quack quackword "newquackword"

Be sure to change the QUACKWORD variable in your environment after the reset
is complete.
	`,
	Run: QuackwordRunner,
}

// QuackwordRunner wraps New for easier testing
func QuackwordRunner(cmd *cobra.Command, args []string) {
	result := Quackword(args...)
	fmt.Println(result)
}

// Quackword resets the QUACKWORD used to encrypt & decrypt entries
func Quackword(args ...string) string {
	if len(args) < 1 {
		return newQuackwordError
	}

	if len(args) > 1 {
		return multiQuackwordError
	}

	newQuackword := args[0]

	// Read old entries and rewrite them with new quackword
	// Make sure no other metadata changes, if possible
	return newQuackword
}

func init() {
	rootCmd.AddCommand(quackwordCmd)
}
