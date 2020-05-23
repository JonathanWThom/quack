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

var verbose bool
var search string
var date string

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Read all entries",
	Long: `
Run quack read to see all entries in normal mode.
Run quack read -v to read in verbose mode. Verbose mode
includes each entry's unique identifier, which can be passed to
quack delete`,
	Run: ReadRunner,
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

	for i := 0; i < len(entries); i++ {
		// should probably filter then format, in separate methods
		result, err := entries[i].Transform(verbose, search, date)
		if err != nil {
			return err.Error()
		}

		if result != "" {
			results = append(results, result)
		}
	}

	return strings.Join(results, "\n\n")
}

func init() {
	rootCmd.AddCommand(readCmd)
	readCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Display entries in verbose mode")
	readCmd.Flags().StringVarP(&search, "search", "s", "", "Search entries by text")
	readCmd.Flags().StringVarP(&date, "date", "d", "", "Search entries by date in format:  \"March 9, 2020\"")
}
