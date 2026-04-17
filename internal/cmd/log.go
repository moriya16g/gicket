package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/gicket/gicket/internal/i18n"
	gitutil "github.com/gicket/gicket/internal/git"
	"github.com/gicket/gicket/internal/store"
	"github.com/spf13/cobra"
)

var logCount int

var logCmd = &cobra.Command{
	Use:   "log <ticket-id>",
	Short: i18n.T("log.short"),
	Long:  i18n.T("log.long"),
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		// チケットの存在確認
		repoPath, err := store.FindRoot(cwd)
		if err != nil {
			return err
		}
		s, err := store.NewStore(repoPath)
		if err != nil {
			return err
		}
		ticket, err := s.Load(args[0])
		if err != nil {
			return err
		}

		// Git リポジトリの確認
		gitRoot, err := gitutil.FindGitRoot(cwd)
		if err != nil {
			return err
		}

		// git log でチケットIDを含むコミットを検索
		// フルIDとファイル名の両方で検索
		logArgs := []string{
			"log", "--all", "--oneline",
			fmt.Sprintf("--max-count=%d", logCount),
			fmt.Sprintf("--grep=%s", ticket.ID),
		}
		out, err := gitutil.RunGit(gitRoot, logArgs...)
		if err != nil {
			// grep にマッチしない場合も出力が空になるだけ
			out = ""
		}

		// gicket:<id> パターンでも検索
		gicketPattern := fmt.Sprintf("gicket:%s", ticket.ID)
		out2, _ := gitutil.RunGit(gitRoot, "log", "--all", "--oneline",
			fmt.Sprintf("--max-count=%d", logCount),
			fmt.Sprintf("--grep=%s", gicketPattern))

		// チケットファイルを変更したコミットも検索
		ticketFile := fmt.Sprintf(".gicket/issues/%s.yml", ticket.ID)
		out3, _ := gitutil.RunGit(gitRoot, "log", "--all", "--oneline",
			fmt.Sprintf("--max-count=%d", logCount),
			"--", ticketFile)

		// 結果を統合（重複除去）
		seen := make(map[string]bool)
		var lines []string
		for _, block := range []string{out, out2, out3} {
			for _, line := range strings.Split(block, "\n") {
				line = strings.TrimSpace(line)
				if line != "" && !seen[line] {
					seen[line] = true
					lines = append(lines, line)
				}
			}
		}

		fmt.Println(i18n.Tf("log.header", ticket.ID, ticket.Title))
		fmt.Println(strings.Repeat("─", 60))

		if len(lines) == 0 {
			fmt.Println(i18n.T("log.no.commits"))
			return nil
		}

		for _, line := range lines {
			fmt.Println(line)
		}
		fmt.Println(i18n.Tf("log.count", len(lines)))
		return nil
	},
}

func init() {
	logCmd.Flags().IntVarP(&logCount, "count", "n", 50, i18n.T("log.flag.count"))
}
