package cmd

import (
	"fmt"
	"os"

	"github.com/gicket/gicket/internal/api"
	"github.com/gicket/gicket/internal/store"
	"github.com/spf13/cobra"
)

var servePort int

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Web UI を起動する",
	Long:  "チケット管理用の Web UI サーバーを起動します。ブラウザでチケットの閲覧・作成・編集ができます。",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		repoPath, err := store.FindRoot(cwd)
		if err != nil {
			return err
		}

		fmt.Printf("Opening http://localhost:%d in your browser...\n", servePort)
		return api.StartServer(repoPath, servePort)
	},
}

func init() {
	serveCmd.Flags().IntVarP(&servePort, "port", "p", 8080, "サーバーのポート番号")
}
