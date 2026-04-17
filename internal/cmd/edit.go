package cmd

import (
	"fmt"
	"os"

	"github.com/gicket/gicket/internal/i18n"
	"github.com/gicket/gicket/internal/model"
	"github.com/gicket/gicket/internal/store"
	"github.com/spf13/cobra"
)

var (
	editTitle    string
	editPriority string
	editStatus   string
	editAssignee string
	editLabels   []string
	editDesc     string
)

var editCmd = &cobra.Command{
	Use:   "edit <id>",
	Short: i18n.T("edit.short"),
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

		if cmd.Flags().Changed("title") {
			ticket.Title = editTitle
		}
		if cmd.Flags().Changed("priority") {
			if err := model.ValidatePriority(editPriority); err != nil {
				return err
			}
			ticket.Priority = model.Priority(editPriority)
		}
		if cmd.Flags().Changed("status") {
			if err := model.ValidateStatus(editStatus); err != nil {
				return err
			}
			ticket.Status = model.Status(editStatus)
		}
		if cmd.Flags().Changed("assignee") {
			ticket.Assignee = editAssignee
		}
		if cmd.Flags().Changed("label") {
			ticket.Labels = editLabels
		}
		if cmd.Flags().Changed("description") {
			ticket.Description = editDesc
		}

		if err := s.Save(ticket); err != nil {
			return err
		}
		fmt.Println(i18n.Tf("edit.success", ticket.ID))
		return nil
	},
}

func init() {
	editCmd.Flags().StringVarP(&editTitle, "title", "t", "", i18n.T("edit.flag.title"))
	editCmd.Flags().StringVarP(&editPriority, "priority", "p", "", i18n.T("edit.flag.priority"))
	editCmd.Flags().StringVarP(&editStatus, "status", "s", "", i18n.T("edit.flag.status"))
	editCmd.Flags().StringVarP(&editAssignee, "assignee", "a", "", i18n.T("edit.flag.assignee"))
	editCmd.Flags().StringSliceVarP(&editLabels, "label", "l", nil, i18n.T("edit.flag.label"))
	editCmd.Flags().StringVarP(&editDesc, "description", "d", "", i18n.T("edit.flag.description"))
}
