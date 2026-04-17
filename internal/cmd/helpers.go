package cmd

import (
	"os/exec"
	"strings"
)

// getGitUser は git config からユーザー情報を取得する
func getGitUser() string {
	name, err := exec.Command("git", "config", "user.name").Output()
	if err != nil {
		return "unknown"
	}
	email, err := exec.Command("git", "config", "user.email").Output()
	if err != nil {
		return strings.TrimSpace(string(name))
	}
	return strings.TrimSpace(string(name)) + " <" + strings.TrimSpace(string(email)) + ">"
}
