package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/gicket/gicket/internal/i18n"
	"github.com/gicket/gicket/internal/model"
	"github.com/gicket/gicket/internal/store"
	"github.com/spf13/cobra"
)

var (
	searchJSON bool
)

var searchCmd = &cobra.Command{
	Use:   "search <keyword>",
	Short: i18n.T("search.short"),
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		keyword := strings.ToLower(args[0])

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		repoPath, err := store.FindRoot(cwd)
		if err != nil {
			return err
		}
		s, err := store.NewStore(repoPath)
		if err != nil {
			return err
		}

		tickets, err := s.List("")
		if err != nil {
			return err
		}

		var results []*model.Ticket
		for _, t := range tickets {
			if containsKeyword(t.Title, keyword) ||
				containsKeyword(t.Description, keyword) ||
				containsKeyword(string(t.Status), keyword) ||
				containsKeyword(string(t.Priority), keyword) ||
				containsKeyword(t.Assignee, keyword) ||
				containsKeywordInLabels(t.Labels, keyword) ||
				containsKeywordInComments(t.Comments, keyword) {
				results = append(results, t)
			}
		}

		if len(results) == 0 {
			if searchJSON {
				fmt.Println("[]")
			} else {
				fmt.Println(i18n.T("search.no.results"))
			}
			return nil
		}

		if searchJSON {
			data, err := json.MarshalIndent(results, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(data))
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "ID\tSTATUS\tPRIORITY\tTITLE\tASSIGNEE\tLABELS\n")
		fmt.Fprintf(w, "--\t------\t--------\t-----\t--------\t------\n")
		for _, t := range results {
			labels := ""
			if len(t.Labels) > 0 {
				labels = strings.Join(t.Labels, ", ")
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
				t.ID, t.Status, t.Priority, t.Title, t.Assignee, labels)
		}
		w.Flush()
		return nil
	},
}

func containsKeyword(s, keyword string) bool {
	return strings.Contains(strings.ToLower(s), keyword)
}

func containsKeywordInLabels(labels []string, keyword string) bool {
	for _, l := range labels {
		if containsKeyword(l, keyword) {
			return true
		}
	}
	return false
}

func containsKeywordInComments(comments []model.Comment, keyword string) bool {
	for _, c := range comments {
		if containsKeyword(c.Body, keyword) || containsKeyword(c.Author, keyword) {
			return true
		}
	}
	return false
}

func init() {
	searchCmd.Flags().BoolVar(&searchJSON, "json", false, i18n.T("flag.json"))
}
