package cmd

import (
	"fmt"
	"github.com/jonathanwthom/quack/secure"
	"github.com/spf13/cobra"
	"sort"
	"strings"
	"time"
)

const (
	unableToReadError = "Unable to read entries."
)

var verbose bool

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

	for i := 0; i < len(entries); i++ {
		entry := entries[i]
		content, err := secure.Decrypt(entry.Content)
		if err != nil {
			return err.Error()
		}
		loc := time.Now().Location()
		formatted := entry.ModTime.In(loc).Format("January 2nd, 2006 - 3:04 PM MST")
		var result string
		if verbose == true {
			key := entry.Key
			result = fmt.Sprintf("%v - %s\n%s", formatted, key, content)
		} else {
			result = fmt.Sprintf("%v\n%s", formatted, content)
		}
		results = append(results, result)
	}

	return strings.Join(results, "\n\n")
}

func init() {
	rootCmd.AddCommand(readCmd)
	readCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Display messages in verbose mode")
}
