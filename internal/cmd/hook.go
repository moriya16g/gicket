package cmd

import (
	"fmt"
	"os"

	"github.com/gicket/gicket/internal/i18n"
	gitutil "github.com/gicket/gicket/internal/git"
	"github.com/spf13/cobra"
)

var hookCmd = &cobra.Command{
	Use:   "hook",
	Short: i18n.T("hook.short"),
	Long:  i18n.T("hook.long"),
}

var hookInstallCmd = &cobra.Command{
	Use:   "install",
	Short: i18n.T("hook.install.short"),
	Long:  i18n.T("hook.install.long"),
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

		fmt.Println(i18n.T("hook.install.success"))
		fmt.Println(i18n.T("hook.install.commitmsg"))
		fmt.Println(i18n.T("hook.install.mergedriver"))
		fmt.Println(i18n.T("hook.install.gitattr"))
		fmt.Println()
		fmt.Println(i18n.T("hook.install.require.id"))
		fmt.Println("  export GICKET_HOOK_REQUIRE_ID=1")
		return nil
	},
}

var hookUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: i18n.T("hook.uninstall.short"),
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

		fmt.Println(i18n.T("hook.uninstall.success"))
		return nil
	},
}

func init() {
	hookCmd.AddCommand(hookInstallCmd)
	hookCmd.AddCommand(hookUninstallCmd)
}
