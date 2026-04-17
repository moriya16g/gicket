package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gicket/gicket/internal/i18n"
	"github.com/gicket/gicket/internal/model"
	"github.com/gicket/gicket/internal/store"
	"github.com/spf13/cobra"
)

var statsJSON bool

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: i18n.T("stats.short"),
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

		tickets, err := s.List("")
		if err != nil {
			return err
		}

		statusCount := map[model.Status]int{}
		priorityCount := map[model.Priority]int{}

		for _, t := range tickets {
			statusCount[t.Status]++
			priorityCount[t.Priority]++
		}

		if statsJSON {
			data := map[string]interface{}{
				"total":    len(tickets),
				"status":   statusCount,
				"priority": priorityCount,
			}
			out, err := json.MarshalIndent(data, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(out))
			return nil
		}

		fmt.Print(i18n.Tf("stats.total", len(tickets)))
		fmt.Println()
		fmt.Println()

		fmt.Println(i18n.T("stats.by.status"))
		for _, st := range model.ValidStatuses {
			fmt.Printf("  %-15s %d\n", st, statusCount[st])
		}
		fmt.Println()

		fmt.Println(i18n.T("stats.by.priority"))
		for _, p := range model.ValidPriorities {
			fmt.Printf("  %-15s %d\n", p, priorityCount[p])
		}

		return nil
	},
}

func init() {
	statsCmd.Flags().BoolVar(&statsJSON, "json", false, i18n.T("flag.json"))
}
