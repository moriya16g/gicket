package model

import "time"

type Status string

const (
	StatusOpen       Status = "open"
	StatusInProgress Status = "in-progress"
	StatusClosed     Status = "closed"
)

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

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
