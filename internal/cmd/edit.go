package cmd

import (
	"fmt"
	"os"

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
	Short: "チケットを編集する",
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
			ticket.Priority = model.Priority(editPriority)
		}
		if cmd.Flags().Changed("status") {
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
		fmt.Printf("チケットを更新しました: %s\n", ticket.ID)
		return nil
	},
}

func init() {
	editCmd.Flags().StringVarP(&editTitle, "title", "t", "", "タイトル")
	editCmd.Flags().StringVarP(&editPriority, "priority", "p", "", "優先度 (low/medium/high)")
	editCmd.Flags().StringVarP(&editStatus, "status", "s", "", "ステータス (open/in-progress/closed)")
	editCmd.Flags().StringVarP(&editAssignee, "assignee", "a", "", "担当者")
	editCmd.Flags().StringSliceVarP(&editLabels, "label", "l", nil, "ラベル")
	editCmd.Flags().StringVarP(&editDesc, "description", "d", "", "説明")
}
