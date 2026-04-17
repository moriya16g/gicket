package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gicket/gicket/internal/i18n"
	"github.com/gicket/gicket/internal/model"
	"github.com/gicket/gicket/internal/store"
	"github.com/spf13/cobra"
)

var (
	newTitle    string
	newPriority string
	newLabels   []string
	newAssignee string
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: i18n.T("new.short"),
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

		if newTitle == "" {
			return errors.New(i18n.T("new.title.required"))
		}

		if err := model.ValidatePriority(newPriority); err != nil {
			return err
		}

		now := time.Now()
		ticket := &model.Ticket{
			ID:       store.GenerateID(),
			Title:    newTitle,
			Status:   model.StatusOpen,
			Priority: model.Priority(newPriority),
			Assignee: newAssignee,
			Labels:   newLabels,
			Created:  now,
			Updated:  now,
			Author:   getGitUser(),
		}

		if err := s.Save(ticket); err != nil {
			return err
		}
		fmt.Println(i18n.Tf("new.success", ticket.ID, ticket.Title))
		return nil
	},
}

func init() {
	newCmd.Flags().StringVarP(&newTitle, "title", "t", "", i18n.T("new.flag.title"))
	newCmd.Flags().StringVarP(&newPriority, "priority", "p", "medium", i18n.T("new.flag.priority"))
	newCmd.Flags().StringSliceVarP(&newLabels, "label", "l", nil, i18n.T("new.flag.label"))
	newCmd.Flags().StringVarP(&newAssignee, "assignee", "a", "", i18n.T("new.flag.assignee"))
}
