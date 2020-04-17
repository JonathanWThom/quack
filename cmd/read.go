package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"sort"
	"strings"
)

const (
	unableToReadError = "Unable to read entries."
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Read all entries",
	Run:   ReadRunner,
}

// ReadRunner wraps Read for easier testing
func ReadRunner(cmd *cobra.Command, args []string) {
	result := Read(args...)
	fmt.Println(result)
}

// Read returns all entries
func Read(args ...string) string {
	entries, err := store.Read()
	if err != nil {
		return unableToReadError
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].ModTime.After(entries[j].ModTime)
	})

	var results []string

	// TODO: handle a lot of entries
	for i := 0; i < len(entries); i++ {
		entry := entries[i]
		// TODO: could this be a function on Entry type?
		// might make sense to move out of storage if so
		result := fmt.Sprintf("%v:\n%s", entry.ModTime, entry.Content)
		results = append(results, result)
	}

	return strings.Join(results, "\n")
}

func init() {
	rootCmd.AddCommand(readCmd)
}
