package cmd

import (
	"fmt"
	"os"

	"github.com/gicket/gicket/internal/api"
	"github.com/gicket/gicket/internal/i18n"
	"github.com/gicket/gicket/internal/store"
	"github.com/spf13/cobra"
)

var servePort int

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: i18n.T("serve.short"),
	Long:  i18n.T("serve.long"),
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		repoPath, err := store.FindRoot(cwd)
		if err != nil {
			return err
		}

		fmt.Println(i18n.Tf("serve.open", servePort))
		return api.StartServer(repoPath, servePort)
	},
}

func init() {
	serveCmd.Flags().IntVarP(&servePort, "port", "p", 8080, i18n.T("serve.flag.port"))
}
