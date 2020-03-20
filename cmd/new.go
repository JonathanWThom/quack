package cmd

import (
	"fmt"
	"log"
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
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 280 {
			log.Fatal("Message must be shorter than 280 characters.")
		}

		msg := strings.Join(args, " ")
		err := storage.Create(msg)
		if err != nil {
			log.Fatal("Failed to create entry.")
		}

		fmt.Println("Entry saved.")
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
