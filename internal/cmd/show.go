package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gicket/gicket/internal/i18n"
	"github.com/gicket/gicket/internal/store"
	"github.com/spf13/cobra"
)

var showJSON bool

var showCmd = &cobra.Command{
	Use:   "show <id>",
	Short: i18n.T("show.short"),
	Args:  cobra.ExactArgs(1),
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

		ticket, err := s.Load(args[0])
		if err != nil {
			return err
		}

		if showJSON {
			data, err := json.MarshalIndent(ticket, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("ID:        %s\n", ticket.ID)
		fmt.Printf("Title:     %s\n", ticket.Title)
		fmt.Printf("Status:    %s\n", ticket.Status)
		fmt.Printf("Priority:  %s\n", ticket.Priority)
		fmt.Printf("Author:    %s\n", ticket.Author)
		fmt.Printf("Assignee:  %s\n", ticket.Assignee)
		if len(ticket.Labels) > 0 {
			fmt.Printf("Labels:    %s\n", strings.Join(ticket.Labels, ", "))
		}
		fmt.Printf("Created:   %s\n", ticket.Created.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated:   %s\n", ticket.Updated.Format("2006-01-02 15:04:05"))

		if ticket.Description != "" {
			fmt.Printf("\n--- Description ---\n%s\n", ticket.Description)
		}

		if len(ticket.Comments) > 0 {
			fmt.Printf("\n--- Comments (%d) ---\n", len(ticket.Comments))
			for i, c := range ticket.Comments {
				fmt.Printf("\n[%d] %s (%s)\n%s\n",
					i+1, c.Author, c.Date.Format("2006-01-02 15:04:05"), c.Body)
			}
		}

		return nil
	},
}

func init() {
	showCmd.Flags().BoolVar(&showJSON, "json", false, i18n.T("flag.json"))
}
