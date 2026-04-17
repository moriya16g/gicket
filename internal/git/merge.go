package git

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/gicket/gicket/internal/model"
	"gopkg.in/yaml.v3"
)

// MergeTicketFiles は3-wayマージでチケットYAMLファイルをマージする
// ancestor: 共通祖先, ours: 現在のブランチ, theirs: マージ元ブランチ
// 結果は ours ファイルに書き込む（git merge driver の規約）
func MergeTicketFiles(ancestorPath, oursPath, theirsPath string) error {
	ancestor, err := loadTicketFile(ancestorPath)
	if err != nil {
		return fmt.Errorf("ancestor の読み込みに失敗: %w", err)
	}
	ours, err := loadTicketFile(oursPath)
	if err != nil {
		return fmt.Errorf("ours の読み込みに失敗: %w", err)
	}
	theirs, err := loadTicketFile(theirsPath)
	if err != nil {
		return fmt.Errorf("theirs の読み込みに失敗: %w", err)
	}

	merged, conflict := mergeTickets(ancestor, ours, theirs)

	data, err := yaml.Marshal(merged)
	if err != nil {
		return fmt.Errorf("マージ結果のマーシャルに失敗: %w", err)
	}
	if err := os.WriteFile(oursPath, data, 0644); err != nil {
		return fmt.Errorf("マージ結果の書き込みに失敗: %w", err)
	}

	if conflict {
		return fmt.Errorf("CONFLICT: チケット %s のフィールドが両方のブランチで異なる値に変更されました", merged.ID)
	}
	return nil
}

func loadTicketFile(path string) (*model.Ticket, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var ticket model.Ticket
	if err := yaml.Unmarshal(data, &ticket); err != nil {
		return nil, err
	}
	return &ticket, nil
}

// mergeTickets は3-wayマージロジックを実装する
// conflict が true の場合、解決できないコンフリクトがあったことを示す
func mergeTickets(ancestor, ours, theirs *model.Ticket) (*model.Ticket, bool) {
	merged := *ours
	conflict := false

	// Title
	merged.Title, conflict = mergeString(ancestor.Title, ours.Title, theirs.Title, conflict)

	// Status — 新しい方を優先
	merged.Status, conflict = mergeStatus(ancestor.Status, ours.Status, theirs.Status, conflict)

	// Priority
	p, c := mergeString(string(ancestor.Priority), string(ours.Priority), string(theirs.Priority), conflict)
	merged.Priority = model.Priority(p)
	conflict = c

	// Assignee
	merged.Assignee, conflict = mergeString(ancestor.Assignee, ours.Assignee, theirs.Assignee, conflict)

	// Description
	merged.Description, conflict = mergeString(ancestor.Description, ours.Description, theirs.Description, conflict)

	// Labels — 和集合
	merged.Labels = mergeLabels(ancestor.Labels, ours.Labels, theirs.Labels)

	// Comments — すべてのコメントを統合（重複除去、日付順）
	merged.Comments = mergeComments(ancestor.Comments, ours.Comments, theirs.Comments)

	// Updated — 最新のタイムスタンプを採用
	merged.Updated = latestTime(ours.Updated, theirs.Updated)

	return &merged, conflict
}

// mergeString は文字列フィールドの3-wayマージ
func mergeString(ancestor, ours, theirs string, prevConflict bool) (string, bool) {
	if ours == theirs {
		return ours, prevConflict
	}
	if ours == ancestor {
		// ours は変更なし → theirs を採用
		return theirs, prevConflict
	}
	if theirs == ancestor {
		// theirs は変更なし → ours を採用
		return ours, prevConflict
	}
	// 両方変更された → theirs を採用しつつコンフリクトを報告
	return theirs, true
}

// mergeStatus はステータスの3-wayマージ（closed を優先）
func mergeStatus(ancestor, ours, theirs model.Status, prevConflict bool) (model.Status, bool) {
	if ours == theirs {
		return ours, prevConflict
	}
	if ours == ancestor {
		return theirs, prevConflict
	}
	if theirs == ancestor {
		return ours, prevConflict
	}
	// 両方変更 → closed > in-progress > open の優先順位
	priority := map[model.Status]int{
		model.StatusOpen:       0,
		model.StatusInProgress: 1,
		model.StatusClosed:     2,
	}
	if priority[ours] >= priority[theirs] {
		return ours, prevConflict
	}
	return theirs, prevConflict
}

// mergeLabels はラベルの和集合を返す
func mergeLabels(ancestor, ours, theirs []string) []string {
	set := make(map[string]bool)

	// ancestor から削除されたラベルを追跡
	ancestorSet := make(map[string]bool)
	for _, l := range ancestor {
		ancestorSet[l] = true
	}

	oursSet := make(map[string]bool)
	for _, l := range ours {
		oursSet[l] = true
	}
	theirsSet := make(map[string]bool)
	for _, l := range theirs {
		theirsSet[l] = true
	}

	// ours にあるラベル（theirs が ancestor から削除したものは除外）
	for _, l := range ours {
		if ancestorSet[l] && !theirsSet[l] {
			continue // theirs が削除した
		}
		set[l] = true
	}
	// theirs にあるラベル（ours が ancestor から削除したものは除外）
	for _, l := range theirs {
		if ancestorSet[l] && !oursSet[l] {
			continue // ours が削除した
		}
		set[l] = true
	}

	var result []string
	for l := range set {
		result = append(result, l)
	}
	sort.Strings(result)
	return result
}

// mergeComments はコメントを統合する（重複除去、日付順ソート）
func mergeComments(ancestor, ours, theirs []model.Comment) []model.Comment {
	type key struct {
		Author string
		Body   string
		Date   int64
	}
	seen := make(map[key]bool)

	// ancestor のコメントキーを記録
	ancestorKeys := make(map[key]bool)
	for _, c := range ancestor {
		k := key{c.Author, c.Body, c.Date.Unix()}
		ancestorKeys[k] = true
	}

	var merged []model.Comment

	addComment := func(c model.Comment) {
		k := key{c.Author, c.Body, c.Date.Unix()}
		if !seen[k] {
			seen[k] = true
			merged = append(merged, c)
		}
	}

	// ours のコメントを追加
	for _, c := range ours {
		addComment(c)
	}
	// theirs のコメントを追加（ours にないもの）
	for _, c := range theirs {
		addComment(c)
	}

	sort.Slice(merged, func(i, j int) bool {
		return merged[i].Date.Before(merged[j].Date)
	})
	return merged
}

func latestTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}
