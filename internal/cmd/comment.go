package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/gicket/gicket/internal/model"
	"github.com/gicket/gicket/internal/store"
	"github.com/spf13/cobra"
)

var commentBody string

var commentCmd = &cobra.Command{
	Use:   "comment <id>",
	Short: "チケットにコメントを追加する",
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

		if commentBody == "" {
			return fmt.Errorf("コメント内容は必須です (-m フラグで指定)")
		}

		ticket, err := s.Load(args[0])
		if err != nil {
			return err
		}

		comment := model.Comment{
			Author: getGitUser(),
			Date:   time.Now(),
			Body:   commentBody,
		}
		ticket.Comments = append(ticket.Comments, comment)

		if err := s.Save(ticket); err != nil {
			return err
		}
		fmt.Printf("コメントを追加しました: %s\n", ticket.ID)
		return nil
	},
}

func init() {
	commentCmd.Flags().StringVarP(&commentBody, "message", "m", "", "コメント内容 (必須)")
}
