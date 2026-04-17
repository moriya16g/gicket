package cmd

import (
	"fmt"
	"os"

	"github.com/gicket/gicket/internal/store"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "現在のディレクトリに gicket を初期化する",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		s, err := store.NewStore(cwd)
		if err != nil {
			return err
		}
		if err := s.Init(); err != nil {
			return err
		}
		fmt.Println("gicket を初期化しました (.gicket/)")
		return nil
	},
}
