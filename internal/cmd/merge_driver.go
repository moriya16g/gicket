package cmd

import (
	"fmt"
	"os"

	gitutil "github.com/gicket/gicket/internal/git"
	"github.com/spf13/cobra"
)

var mergeDriverCmd = &cobra.Command{
	Use:    "merge-driver <ancestor> <ours> <theirs>",
	Short:  "カスタムマージドライバ（git が内部的に呼び出す）",
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
