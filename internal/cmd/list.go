package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/gicket/gicket/internal/model"
	"github.com/gicket/gicket/internal/store"
	"github.com/spf13/cobra"
)

var listAll bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "チケットの一覧を表示する",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
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

		var filter model.Status
		if !listAll {
			filter = model.StatusOpen
		}

		tickets, err := s.List(filter)
		if err != nil {
			return err
		}

		if len(tickets) == 0 {
			fmt.Println("チケットはありません")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "ID\tSTATUS\tPRIORITY\tTITLE\tASSIGNEE\tLABELS\n")
		fmt.Fprintf(w, "--\t------\t--------\t-----\t--------\t------\n")
		for _, t := range tickets {
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

func init() {
	listCmd.Flags().BoolVarP(&listAll, "all", "a", false, "すべてのステータスのチケットを表示")
}
