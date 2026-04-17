package cmd

import (
	"fmt"
	"os"

	"github.com/gicket/gicket/internal/i18n"
	gitutil "github.com/gicket/gicket/internal/git"
	"github.com/spf13/cobra"
)

var mergeDriverCmd = &cobra.Command{
	Use:    "merge-driver <ancestor> <ours> <theirs>",
	Short:  i18n.T("merge_driver.short"),
	Hidden: true,
	Args:   cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		ancestorPath := args[0]
		oursPath := args[1]
		theirsPath := args[2]

		err := gitutil.MergeTicketFiles(ancestorPath, oursPath, theirsPath)
		if err != nil {
			// コンフリクトの場合、git に非ゼロ終了コードを返す
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		return nil
	},
}
