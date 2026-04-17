package cmd

import (
	"fmt"
	"os"
	"time"

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
	Short: "新しいチケットを作成する",
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
			return fmt.Errorf("タイトルは必須です (-t フラグで指定)")
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
		fmt.Printf("チケットを作成しました: %s - %s\n", ticket.ID, ticket.Title)
		return nil
	},
}

func init() {
	newCmd.Flags().StringVarP(&newTitle, "title", "t", "", "チケットのタイトル (必須)")
	newCmd.Flags().StringVarP(&newPriority, "priority", "p", "medium", "優先度 (low/medium/high)")
	newCmd.Flags().StringSliceVarP(&newLabels, "label", "l", nil, "ラベル (複数指定可)")
	newCmd.Flags().StringVarP(&newAssignee, "assignee", "a", "", "担当者")
}
