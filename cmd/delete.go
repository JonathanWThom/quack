package cmd

import (
	"fmt"

	"github.com/jonathanwthom/quack/secure"
	"github.com/spf13/cobra"
)

const (
	deleteSuccessMsg    = "Successfully deleted entry."
	unableToDeleteError = "Unable to delete entry."
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: DeleteRunner,
}

func DeleteRunner(cmd *cobra.Command, args []string) {
	result := Delete(args...)
	fmt.Println(result)
}

func Delete(args ...string) string {
	key := args[0]
	entry, err := store.ReadByKey(key)
	if err != nil {
		return unableToDeleteError
	}

	_, err = secure.Decrypt(entry.Content)
	if err != nil {
		return err.Error()
	}

	err = store.Delete(key)
	if err != nil {
		return unableToDeleteError
	}

	return "Entry deleted."
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
