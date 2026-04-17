package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gicket",
	Short: "Git リポジトリ内で動作する分散チケット管理ツール",
	Long: `gicket は Git リポジトリ内にテキスト(YAML)ベースでチケットを管理するツールです。
WEBサーバ不要で、Git の push/pull だけで開発者間のチケット共有が可能です。`,
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
}
