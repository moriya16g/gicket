package model

import (
	"fmt"
	"time"
)

type Status string

const (
	StatusOpen       Status = "open"
	StatusInProgress Status = "in-progress"
	StatusClosed     Status = "closed"
)

// ValidStatuses は有効なステータスの一覧
var ValidStatuses = []Status{StatusOpen, StatusInProgress, StatusClosed}

// IsValidStatus はステータスが有効かどうかを返す
func IsValidStatus(s string) bool {
	for _, v := range ValidStatuses {
		if string(v) == s {
			return true
		}
	}
	return false
}

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

// ValidPriorities は有効な優先度の一覧
var ValidPriorities = []Priority{PriorityLow, PriorityMedium, PriorityHigh}

// IsValidPriority は優先度が有効かどうかを返す
func IsValidPriority(p string) bool {
	for _, v := range ValidPriorities {
		if string(v) == p {
			return true
		}
	}
	return false
}

// ValidateStatus はステータスを検証し、無効な場合エラーを返す
func ValidateStatus(s string) error {
	if !IsValidStatus(s) {
		return fmt.Errorf("invalid status: %q (must be one of: open, in-progress, closed)", s)
	}
	return nil
}

// ValidatePriority は優先度を検証し、無効な場合エラーを返す
func ValidatePriority(p string) error {
	if !IsValidPriority(p) {
		return fmt.Errorf("invalid priority: %q (must be one of: low, medium, high)", p)
	}
	return nil
}

type Comment struct {
	Author string    `yaml:"author" json:"author"`
	Date   time.Time `yaml:"date" json:"date"`
	Body   string    `yaml:"body" json:"body"`
}

type Ticket struct {
	ID          string    `yaml:"id" json:"id"`
	Title       string    `yaml:"title" json:"title"`
	Status      Status    `yaml:"status" json:"status"`
	Priority    Priority  `yaml:"priority" json:"priority"`
	Assignee    string    `yaml:"assignee,omitempty" json:"assignee"`
	Labels      []string  `yaml:"labels,omitempty" json:"labels"`
	Created     time.Time `yaml:"created" json:"created"`
	Updated     time.Time `yaml:"updated" json:"updated"`
	Author      string    `yaml:"author" json:"author"`
	Description string    `yaml:"description" json:"description"`
	Comments    []Comment `yaml:"comments,omitempty" json:"comments"`
}
