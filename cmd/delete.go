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
	Short: "Delete an entry",
	Long: `
Delete an entry by running quack delete <unique-id>.
The unique id of an entry can be found by running quack read -v`,
	Run: deleteRunner,
}

func deleteRunner(cmd *cobra.Command, args []string) {
	result := Delete(args...)
	fmt.Println(result)
}

// Delete removes entries by key
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
