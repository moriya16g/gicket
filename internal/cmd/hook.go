package cmd

import (
	"fmt"
	"os"

	gitutil "github.com/gicket/gicket/internal/git"
	"github.com/spf13/cobra"
)

var hookCmd = &cobra.Command{
	Use:   "hook",
	Short: "Git フックを管理する",
	Long:  "commit-msg フックとカスタムマージドライバをインストール/アンインストールします。",
}

var hookInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Git フックとマージドライバをインストールする",
	Long: `以下をインストールします:
  - commit-msg フック: コミットメッセージ内のチケットID参照を検証
  - カスタムマージドライバ: .gicket/issues/*.yml の3-wayマージを自動処理
  - .gitattributes: マージドライバの適用ルール`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		gitRoot, err := gitutil.FindGitRoot(cwd)
		if err != nil {
			return err
		}

		if err := gitutil.InstallHooks(gitRoot); err != nil {
			return err
		}

		fmt.Println("Git フックをインストールしました:")
		fmt.Println("  ✓ commit-msg フック")
		fmt.Println("  ✓ カスタムマージドライバ (merge.gicket)")
		fmt.Println("  ✓ .gitattributes")
		fmt.Println()
		fmt.Println("チケットID参照を必須にするには:")
		fmt.Println("  export GICKET_HOOK_REQUIRE_ID=1")
		return nil
	},
}

var hookUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Git フックとマージドライバをアンインストールする",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		gitRoot, err := gitutil.FindGitRoot(cwd)
		if err != nil {
			return err
		}

		if err := gitutil.UninstallHooks(gitRoot); err != nil {
			return err
		}

		fmt.Println("Git フックをアンインストールしました")
		return nil
	},
}

func init() {
	hookCmd.AddCommand(hookInstallCmd)
	hookCmd.AddCommand(hookUninstallCmd)
}
