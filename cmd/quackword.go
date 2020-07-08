package cmd

import (
	"fmt"

	"github.com/jonathanwthom/quack/secure"
	"github.com/spf13/cobra"
	"os"
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
	entries, err := store.Read()
	if err != nil {
		return err.Error()
	}

	// order this?
	for i := 0; i < len(entries); i++ {
		entry := entries[i]
		decrypted, err := secure.Decrypt(entry.Content)
		if err != nil {
			// should probably be a friendly message
			return err.Error()
		}

		encrypted, err := secure.EncryptWithNewQuackword(decrypted, newQuackword)
		if err != nil {
			return err.Error()
		}

		// This will work, but will lose all timing data
		store.Create(encrypted)
	}

	// Read old entries and rewrite them with new quackword
	// Make sure no other metadata changes, if possible

	return "Add a note about setting new quackword in permanent environmentenvironment"
}

func init() {
	rootCmd.AddCommand(quackwordCmd)
}
