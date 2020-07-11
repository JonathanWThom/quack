package cmd

import (
	"fmt"
	"github.com/jonathanwthom/quack/storage"
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
var number int

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Read all entries",
	Long: `
Run quack read to see last 10 entries in normal mode.
See more entries by passing the -n flag, e.g quack read -n 30 for last 30 entries.
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
		return entries[i].CreatedAt.After(entries[j].CreatedAt)
	})

	var results []string

	for i := 0; i < count(entries); i++ {
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

func count(entries []storage.Entry) int {
	if number != 0 && number <= len(entries) {
		return number
	} else if len(entries) < 10 {
		return len(entries)
	}

	return 10
}

func init() {
	rootCmd.AddCommand(readCmd)
	readCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Display entries in verbose mode")
	readCmd.Flags().StringVarP(&search, "search", "s", "", "Search entries by text")
	readCmd.Flags().StringVarP(&date, "date", "d", "", "Search entries by date in format:  \"March 9, 2020\"")
	readCmd.Flags().IntVarP(&number, "number", "n", 0, "Return last n entries")
}
