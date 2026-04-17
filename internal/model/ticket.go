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
	Author string    `yaml:"author"`
	Date   time.Time `yaml:"date"`
	Body   string    `yaml:"body"`
}

type Ticket struct {
	ID          string    `yaml:"id"`
	Title       string    `yaml:"title"`
	Status      Status    `yaml:"status"`
	Priority    Priority  `yaml:"priority"`
	Assignee    string    `yaml:"assignee,omitempty"`
	Labels      []string  `yaml:"labels,omitempty"`
	Created     time.Time `yaml:"created"`
	Updated     time.Time `yaml:"updated"`
	Author      string    `yaml:"author"`
	Description string    `yaml:"description"`
	Comments    []Comment `yaml:"comments,omitempty"`
}
