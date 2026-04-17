package git

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gicket/gicket/internal/i18n"
)

const commitMsgHookScript = `#!/bin/sh
# gicket commit-msg hook
# Validates that commit messages contain a ticket ID reference (optional).
#
# Pattern: gicket:<ID> or gicket:<prefix>
# Example: "Fix login bug gicket:20260416-200633"
#
# Set GICKET_HOOK_REQUIRE_ID=1 to make ticket ID references mandatory.

msg=$(cat "$1")

if [ -n "$GICKET_HOOK_REQUIRE_ID" ] && [ "$GICKET_HOOK_REQUIRE_ID" = "1" ]; then
    if ! echo "$msg" | grep -qE 'gicket:[0-9]{8}-[0-9]{6}'; then
        echo "ERROR: Commit message does not contain a ticket ID"
        echo "  Format: gicket:<ticket-id>"
        echo "  Example: Fix bug gicket:20260416-200633"
        echo ""
        echo "  To disable this check: GICKET_HOOK_REQUIRE_ID=0 git commit ..."
        exit 1
    fi
fi
`

const commitMsgHookScriptWindows = `#!/bin/sh
# gicket commit-msg hook
# Validates that commit messages contain a ticket ID reference (optional).
#
# Pattern: gicket:<ID> or gicket:<prefix>
# Example: "Fix login bug gicket:20260416-200633"
#
# Set GICKET_HOOK_REQUIRE_ID=1 to make ticket ID references mandatory.

msg=$(cat "$1")

if [ -n "$GICKET_HOOK_REQUIRE_ID" ] && [ "$GICKET_HOOK_REQUIRE_ID" = "1" ]; then
    if ! echo "$msg" | grep -qE 'gicket:[0-9]{8}-[0-9]{6}'; then
        echo "ERROR: Commit message does not contain a ticket ID"
        echo "  Format: gicket:<ticket-id>"
        echo "  Example: Fix bug gicket:20260416-200633"
        echo ""
        echo "  To disable this check: GICKET_HOOK_REQUIRE_ID=0 git commit ..."
        exit 1
    fi
fi
`

const hookMarker = "# gicket commit-msg hook"

// InstallHooks は .git/hooks にフックスクリプトをインストールする
// また、.gitattributes にカスタムマージドライバの設定を追加する
func InstallHooks(gitRoot string) error {
	if _, err := GitExecutable(); err != nil {
		return err
	}

	// 1. commit-msg フック
	if err := installCommitMsgHook(gitRoot); err != nil {
		return fmt.Errorf(i18n.T("git.hook.install.fail"), err)
	}

	// 2. カスタムマージドライバの登録
	if err := installMergeDriver(gitRoot); err != nil {
		return fmt.Errorf(i18n.T("git.merge.driver.fail"), err)
	}

	// 3. .gitattributes の設定
	if err := installGitAttributes(gitRoot); err != nil {
		return fmt.Errorf(i18n.T("git.gitattr.fail"), err)
	}

	return nil
}

func UninstallHooks(gitRoot string) error {
	if _, err := GitExecutable(); err != nil {
		return err
	}

	if err := uninstallCommitMsgHook(gitRoot); err != nil {
		return fmt.Errorf(i18n.T("git.hook.remove.fail"), err)
	}

	RunGit(gitRoot, "config", "--local", "--remove-section", "merge.gicket")

	if err := uninstallGitAttributes(gitRoot); err != nil {
		return fmt.Errorf(i18n.T("git.gitattr.update"), err)
	}

	return nil
}

func installCommitMsgHook(gitRoot string) error {
	hooksDir := HooksDir(gitRoot)
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		return err
	}

	hookPath := filepath.Join(hooksDir, "commit-msg")

	// 既存のフックがある場合、gicket のものかチェック
	if data, err := os.ReadFile(hookPath); err == nil {
		content := string(data)
		if strings.Contains(content, hookMarker) {
			return nil // 既にインストール済み
		}
		// 他のフックが存在する → 追記しない（上書きは危険）
		return errors.New(i18n.Tf("git.hook.exists", hookPath))
	}

	script := commitMsgHookScript
	if runtime.GOOS == "windows" {
		script = commitMsgHookScriptWindows
	}

	return os.WriteFile(hookPath, []byte(script), 0755)
}

func uninstallCommitMsgHook(gitRoot string) error {
	hookPath := filepath.Join(HooksDir(gitRoot), "commit-msg")
	data, err := os.ReadFile(hookPath)
	if err != nil {
		return nil // 存在しない
	}
	if !strings.Contains(string(data), hookMarker) {
		return nil // gicket のものではない
	}
	return os.Remove(hookPath)
}

func installMergeDriver(gitRoot string) error {
	// gicket の実行ファイルパスを取得
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	// git config にカスタムマージドライバを登録
	driverCmd := fmt.Sprintf("%s merge-driver %%O %%A %%B", exe)
	if _, err := RunGit(gitRoot, "config", "--local", "merge.gicket.name", "gicket ticket merge driver"); err != nil {
		return err
	}
	if _, err := RunGit(gitRoot, "config", "--local", "merge.gicket.driver", driverCmd); err != nil {
		return err
	}
	return nil
}

const gitAttributesLine = ".gicket/issues/*.yml merge=gicket"
const gitAttributesMarker = "# gicket merge driver"

func installGitAttributes(gitRoot string) error {
	attrPath := filepath.Join(gitRoot, ".gitattributes")
	entry := gitAttributesMarker + "\n" + gitAttributesLine + "\n"

	data, err := os.ReadFile(attrPath)
	if err != nil {
		// ファイルが存在しない → 新規作成
		return os.WriteFile(attrPath, []byte(entry), 0644)
	}

	content := string(data)
	if strings.Contains(content, gitAttributesLine) {
		return nil // 既に設定済み
	}

	// 末尾に追加
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	content += entry
	return os.WriteFile(attrPath, []byte(content), 0644)
}

func uninstallGitAttributes(gitRoot string) error {
	attrPath := filepath.Join(gitRoot, ".gitattributes")
	data, err := os.ReadFile(attrPath)
	if err != nil {
		return nil
	}

	lines := strings.Split(string(data), "\n")
	var kept []string
	for _, line := range lines {
		if line == gitAttributesMarker || line == gitAttributesLine {
			continue
		}
		kept = append(kept, line)
	}

	result := strings.Join(kept, "\n")
	// 空ファイルになった場合は削除
	trimmed := strings.TrimSpace(result)
	if trimmed == "" {
		return os.Remove(attrPath)
	}
	return os.WriteFile(attrPath, []byte(result), 0644)
}
