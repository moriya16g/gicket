package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// FindGitRoot は .git ディレクトリを探して Git リポジトリのルートを返す
func FindGitRoot(startDir string) (string, error) {
	dir := startDir
	for {
		gitPath := filepath.Join(dir, ".git")
		if info, err := os.Stat(gitPath); err == nil && (info.IsDir() || info.Mode().IsRegular()) {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("Git リポジトリが見つかりません")
		}
		dir = parent
	}
}

// HooksDir は .git/hooks ディレクトリのパスを返す
func HooksDir(gitRoot string) string {
	return filepath.Join(gitRoot, ".git", "hooks")
}

// RunGit は git コマンドを実行し、stdout を返す
func RunGit(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}
	return strings.TrimSpace(string(out)), nil
}

// IsInGitRepo は指定ディレクトリが Git リポジトリ内かどうかを返す
func IsInGitRepo(dir string) bool {
	_, err := FindGitRoot(dir)
	return err == nil
}

// GitExecutable は git コマンドのパスを返す（存在しない場合エラー）
func GitExecutable() (string, error) {
	path, err := exec.LookPath("git")
	if err != nil {
		return "", fmt.Errorf("git コマンドが見つかりません。Git をインストールしてください")
	}
	return path, nil
}
