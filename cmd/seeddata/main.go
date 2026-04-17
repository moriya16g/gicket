package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gicket/gicket/internal/model"
	"github.com/gicket/gicket/internal/store"
)

var titles = []struct {
	title    string
	labels   []string
	priority model.Priority
}{
	{"Fix login page crash on mobile Safari", []string{"bug", "frontend", "mobile"}, model.PriorityHigh},
	{"Add password strength indicator", []string{"feature", "frontend", "security"}, model.PriorityMedium},
	{"Database connection pool exhaustion under load", []string{"bug", "backend", "database"}, model.PriorityHigh},
	{"Implement user avatar upload", []string{"feature", "frontend", "backend"}, model.PriorityMedium},
	{"API rate limiting not working for /search endpoint", []string{"bug", "backend", "api"}, model.PriorityHigh},
	{"Add dark mode toggle", []string{"feature", "frontend", "ui"}, model.PriorityLow},
	{"Optimize image compression pipeline", []string{"improvement", "backend", "performance"}, model.PriorityMedium},
	{"Email notification delay exceeds 5 minutes", []string{"bug", "backend", "notification"}, model.PriorityHigh},
	{"Create onboarding tutorial for new users", []string{"feature", "frontend", "ux"}, model.PriorityMedium},
	{"Fix broken pagination on search results", []string{"bug", "frontend", "api"}, model.PriorityMedium},
	{"Add CSV export for analytics dashboard", []string{"feature", "backend", "analytics"}, model.PriorityLow},
	{"Memory leak in WebSocket handler", []string{"bug", "backend", "websocket"}, model.PriorityHigh},
	{"Implement two-factor authentication", []string{"feature", "security", "backend"}, model.PriorityHigh},
	{"Update third-party dependencies", []string{"chore", "security"}, model.PriorityMedium},
	{"Fix timezone handling in scheduling module", []string{"bug", "backend", "i18n"}, model.PriorityMedium},
	{"Add keyboard shortcuts for power users", []string{"feature", "frontend", "ux"}, model.PriorityLow},
	{"Redesign settings page layout", []string{"improvement", "frontend", "ui"}, model.PriorityLow},
	{"Fix file upload size limit error message", []string{"bug", "frontend", "ux"}, model.PriorityLow},
	{"Implement audit log for admin actions", []string{"feature", "backend", "security"}, model.PriorityMedium},
	{"Add Slack integration for notifications", []string{"feature", "integration"}, model.PriorityMedium},
	{"Performance regression in dashboard loading", []string{"bug", "frontend", "performance"}, model.PriorityHigh},
	{"Create API documentation with OpenAPI spec", []string{"docs", "api"}, model.PriorityMedium},
	{"Fix CORS error on subdomain requests", []string{"bug", "backend", "api"}, model.PriorityHigh},
	{"Add bulk delete functionality for admin", []string{"feature", "frontend", "backend"}, model.PriorityMedium},
	{"Implement search suggestions / autocomplete", []string{"feature", "frontend", "ux"}, model.PriorityLow},
	{"Fix data race in concurrent user updates", []string{"bug", "backend", "concurrency"}, model.PriorityHigh},
	{"Add multi-language support (i18n)", []string{"feature", "frontend", "i18n"}, model.PriorityMedium},
	{"Improve error messages for form validation", []string{"improvement", "frontend", "ux"}, model.PriorityLow},
	{"Set up CI/CD pipeline with GitHub Actions", []string{"chore", "devops"}, model.PriorityMedium},
	{"Fix SSO login redirect loop", []string{"bug", "backend", "security"}, model.PriorityHigh},
}

var assignees = []string{
	"tanaka@example.com",
	"suzuki@example.com",
	"sato@example.com",
	"yamada@example.com",
	"watanabe@example.com",
	"",
}

var authors = []string{
	"Tanaka Taro <tanaka@example.com>",
	"Suzuki Hanako <suzuki@example.com>",
	"Sato Kenji <sato@example.com>",
	"Yamada Yuki <yamada@example.com>",
	"Watanabe Rina <watanabe@example.com>",
}

var commentBodies = []string{
	"I can reproduce this issue. Working on a fix.",
	"This is related to #12. We should fix them together.",
	"Added a unit test to cover this case.",
	"Needs design review before implementation.",
	"Fixed in the latest commit. Please verify.",
	"Low priority for now, moving to next sprint.",
	"Discussed in today's standup. Will start tomorrow.",
	"Blocked by the database migration task.",
	"Deployed to staging for testing.",
	"Confirmed fixed. Closing this ticket.",
}

func main() {
	targetDir := os.Args[1]

	s, _ := store.NewStore(targetDir)
	s.Init()

	statuses := []model.Status{model.StatusOpen, model.StatusOpen, model.StatusOpen, model.StatusInProgress, model.StatusInProgress, model.StatusClosed}

	for i, item := range titles {
		baseTime := time.Now().Add(-time.Duration(len(titles)-i) * 24 * time.Hour)

		ticket := &model.Ticket{
			ID:          store.GenerateID(),
			Title:       item.title,
			Status:      statuses[rand.Intn(len(statuses))],
			Priority:    item.priority,
			Assignee:    assignees[rand.Intn(len(assignees))],
			Labels:      item.labels,
			Created:     baseTime,
			Updated:     baseTime.Add(time.Duration(rand.Intn(48)) * time.Hour),
			Author:      authors[rand.Intn(len(authors))],
			Description: fmt.Sprintf("Detailed description for: %s\n\nThis needs to be addressed as part of our ongoing improvements.", item.title),
		}

		// Add 0-3 comments
		numComments := rand.Intn(4)
		for j := 0; j < numComments; j++ {
			ticket.Comments = append(ticket.Comments, model.Comment{
				Author: authors[rand.Intn(len(authors))],
				Date:   baseTime.Add(time.Duration(j+1) * 6 * time.Hour),
				Body:   commentBodies[rand.Intn(len(commentBodies))],
			})
		}

		s.Save(ticket)
		fmt.Printf("[%2d] %s  %-12s %-6s %s\n", i+1, ticket.ID, ticket.Status, ticket.Priority, ticket.Title)
	}

	fmt.Printf("\n%d tickets created in %s\n", len(titles), targetDir)
}
