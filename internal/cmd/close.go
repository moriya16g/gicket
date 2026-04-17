package cmd

import (
	"fmt"
	"os"

	"github.com/gicket/gicket/internal/model"
	"github.com/gicket/gicket/internal/store"
	"github.com/spf13/cobra"
)

var closeCmd = &cobra.Command{
	Use:   "close <id>",
	Short: "チケットをクローズする",
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

		if ticket.Status == model.StatusClosed {
			return fmt.Errorf("チケット %s は既にクローズされています", ticket.ID)
		}

		ticket.Status = model.StatusClosed
		if err := s.Save(ticket); err != nil {
			return err
		}
		fmt.Printf("チケットをクローズしました: %s - %s\n", ticket.ID, ticket.Title)
		return nil
	},
}
