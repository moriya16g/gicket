package cmd

import (
	"github.com/gicket/gicket/internal/i18n"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gicket",
	Short: i18n.T("root.short"),
	Long:  i18n.T("root.long"),
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(commentCmd)
	rootCmd.AddCommand(closeCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(hookCmd)
	rootCmd.AddCommand(logCmd)
	rootCmd.AddCommand(mergeDriverCmd)
}
