package cmd

import (
	"fmt"
	"os"

	"github.com/gicket/gicket/internal/i18n"
	"github.com/gicket/gicket/internal/model"
	"github.com/gicket/gicket/internal/store"
	"github.com/spf13/cobra"
)

var closeCmd = &cobra.Command{
	Use:   "close <id>",
	Short: i18n.T("close.short"),
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
			return fmt.Errorf(i18n.Tf("close.already.closed", ticket.ID))
		}

		ticket.Status = model.StatusClosed
		if err := s.Save(ticket); err != nil {
			return err
		}
		fmt.Println(i18n.Tf("close.success", ticket.ID, ticket.Title))
		return nil
	},
}
