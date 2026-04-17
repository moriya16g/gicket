package store

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gicket/gicket/internal/i18n"
	"github.com/gicket/gicket/internal/model"
	"gopkg.in/yaml.v3"
)

const (
	GicketDir   = ".gicket"
	IssuesDir   = "issues"
	ConfigFile  = "config.yml"
)

type Store struct {
	Root string // .gicket ディレクトリの絶対パス
}

func NewStore(repoPath string) (*Store, error) {
	root := filepath.Join(repoPath, GicketDir)
	return &Store{Root: root}, nil
}

// FindRoot は現在のディレクトリから上位に向かって .gicket を探す
func FindRoot(startDir string) (string, error) {
	dir := startDir
	for {
		gicketPath := filepath.Join(dir, GicketDir)
		if info, err := os.Stat(gicketPath); err == nil && info.IsDir() {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New(i18n.T("store.gicket.not.found"))
		}
		dir = parent
	}
}

// Init は .gicket ディレクトリ構造を作成する
func (s *Store) Init() error {
	dirs := []string{
		s.Root,
		filepath.Join(s.Root, IssuesDir),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return fmt.Errorf(i18n.T("store.dir.create.failed"), err)
		}
	}

	configPath := filepath.Join(s.Root, ConfigFile)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config := map[string]string{
			"version": "1",
		}
		data, err := yaml.Marshal(config)
		if err != nil {
			return fmt.Errorf(i18n.T("store.config.marshal"), err)
		}
		if err := os.WriteFile(configPath, data, 0644); err != nil {
			return fmt.Errorf(i18n.T("store.config.write"), err)
		}
	}
	return nil
}

// Save はチケットをYAMLファイルとして保存する
func (s *Store) Save(ticket *model.Ticket) error {
	ticket.Updated = time.Now()
	data, err := yaml.Marshal(ticket)
	if err != nil {
		return fmt.Errorf(i18n.T("store.ticket.marshal"), err)
	}
	filePath := s.ticketPath(ticket.ID)
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf(i18n.T("store.ticket.save"), err)
	}
	return nil
}

// Load は指定IDのチケットを読み込む
func (s *Store) Load(id string) (*model.Ticket, error) {
	// 短縮IDでのマッチを試みる
	fullID, err := s.resolveID(id)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(s.ticketPath(fullID))
	if err != nil {
		return nil, fmt.Errorf(i18n.T("store.ticket.read"), err)
	}

	var ticket model.Ticket
	if err := yaml.Unmarshal(data, &ticket); err != nil {
		return nil, fmt.Errorf(i18n.T("store.ticket.parse"), err)
	}
	return &ticket, nil
}

// List は条件に合うチケットの一覧を返す
func (s *Store) List(statusFilter model.Status) ([]*model.Ticket, error) {
	issuesDir := filepath.Join(s.Root, IssuesDir)
	entries, err := os.ReadDir(issuesDir)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("store.ticket.list"), err)
	}

	var tickets []*model.Ticket
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yml") {
			continue
		}
		id := strings.TrimSuffix(entry.Name(), ".yml")
		ticket, err := s.Load(id)
		if err != nil {
			continue
		}
		if statusFilter != "" && ticket.Status != statusFilter {
			continue
		}
		tickets = append(tickets, ticket)
	}

	sort.Slice(tickets, func(i, j int) bool {
		return tickets[i].Created.After(tickets[j].Created)
	})

	return tickets, nil
}

// Delete はチケットファイルを削除する
func (s *Store) Delete(id string) error {
	fullID, err := s.resolveID(id)
	if err != nil {
		return err
	}
	return os.Remove(s.ticketPath(fullID))
}

func (s *Store) ticketPath(id string) string {
	return filepath.Join(s.Root, IssuesDir, id+".yml")
}

// resolveID は短縮IDからフルIDを解決する
func (s *Store) resolveID(id string) (string, error) {
	// まず完全一致を試す
	if _, err := os.Stat(s.ticketPath(id)); err == nil {
		return id, nil
	}

	// 前方一致で検索
	issuesDir := filepath.Join(s.Root, IssuesDir)
	entries, err := os.ReadDir(issuesDir)
	if err != nil {
		return "", fmt.Errorf(i18n.T("store.ticket.search"), err)
	}

	var matches []string
	for _, entry := range entries {
		name := strings.TrimSuffix(entry.Name(), ".yml")
		if strings.HasPrefix(name, id) {
			matches = append(matches, name)
		}
	}

	switch len(matches) {
	case 0:
		return "", errors.New(i18n.Tf("store.ticket.not.found", id))
	case 1:
		return matches[0], nil
	default:
		return "", errors.New(i18n.Tf("store.ticket.ambiguous", id, matches))
	}
}

// GenerateID はタイムスタンプ+ランダム文字列のIDを生成する
func GenerateID() string {
	now := time.Now()
	b := make([]byte, 3)
	rand.Read(b)
	return fmt.Sprintf("%s-%x", now.Format("20060102-150405"), b)
}
