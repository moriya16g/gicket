package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/gicket/gicket/internal/i18n"
	"github.com/gicket/gicket/internal/model"
	"github.com/gicket/gicket/internal/store"
	"github.com/spf13/cobra"
)

var reopenCmd = &cobra.Command{
	Use:   "reopen <id>",
	Short: i18n.T("reopen.short"),
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

		if ticket.Status == model.StatusOpen {
			return errors.New(i18n.Tf("reopen.already.open", ticket.ID))
		}

		ticket.Status = model.StatusOpen
		if err := s.Save(ticket); err != nil {
			return err
		}
		fmt.Println(i18n.Tf("reopen.success", ticket.ID, ticket.Title))
		return nil
	},
}
