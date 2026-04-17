package i18n

var messagesEN = map[string]string{
	// root
	"root.short": "A distributed ticket management tool for Git repositories",
	"root.long":  "gicket manages tickets as text (YAML) files inside a Git repository.\nNo web server needed — share tickets with your team via git push/pull.",

	// init
	"init.short":   "Initialize gicket in the current directory",
	"init.success": "Initialized gicket (.gicket/)",

	// new
	"new.short":         "Create a new ticket",
	"new.title.required": "Title is required (use -t flag)",
	"new.success":       "Created ticket: %s - %s",
	"new.flag.title":    "Ticket title (required)",
	"new.flag.priority": "Priority (low/medium/high)",
	"new.flag.label":    "Labels (multiple allowed)",
	"new.flag.assignee": "Assignee",

	// list
	"list.short":      "List tickets",
	"list.no.tickets": "No tickets found",
	"list.flag.all":   "Show tickets of all statuses",

	// show
	"show.short": "Show ticket details",

	// edit
	"edit.short":            "Edit a ticket",
	"edit.success":          "Updated ticket: %s",
	"edit.flag.title":       "Title",
	"edit.flag.priority":    "Priority (low/medium/high)",
	"edit.flag.status":      "Status (open/in-progress/closed)",
	"edit.flag.assignee":    "Assignee",
	"edit.flag.label":       "Labels",
	"edit.flag.description": "Description",

	// comment
	"comment.short":         "Add a comment to a ticket",
	"comment.body.required": "Comment body is required (use -m flag)",
	"comment.success":       "Added comment: %s",
	"comment.flag.message":  "Comment body (required)",

	// close
	"close.short":          "Close a ticket",
	"close.already.closed": "Ticket %s is already closed",
	"close.success":        "Closed ticket: %s - %s",

	// serve
	"serve.short": "Start Web UI",
	"serve.long":  "Start a Web UI server for ticket management. Browse, create, and edit tickets from your browser.",
	"serve.open":  "Opening http://localhost:%d in your browser...",
	"serve.flag.port": "Server port number",

	// hook
	"hook.short":              "Manage Git hooks",
	"hook.long":               "Install/uninstall commit-msg hook and custom merge driver.",
	"hook.install.short":      "Install Git hooks and merge driver",
	"hook.install.long":       "Installs the following:\n  - commit-msg hook: validates ticket ID references in commit messages\n  - Custom merge driver: smart 3-way merge for .gicket/issues/*.yml\n  - .gitattributes: merge driver rules",
	"hook.install.success":    "Installed Git hooks:",
	"hook.install.commitmsg":  "  ✓ commit-msg hook",
	"hook.install.mergedriver":"  ✓ Custom merge driver (merge.gicket)",
	"hook.install.gitattr":    "  ✓ .gitattributes",
	"hook.install.require.id": "To require ticket ID references:",
	"hook.uninstall.short":    "Uninstall Git hooks and merge driver",
	"hook.uninstall.success":  "Uninstalled Git hooks",

	// log
	"log.short":      "Show Git commit history related to a ticket",
	"log.long":       "Search and display commits that mention the ticket ID in their message.",
	"log.header":     "Ticket: %s - %s",
	"log.no.commits": "No related commits found",
	"log.count":      "\n%d commit(s) found",
	"log.flag.count": "Maximum number of commits to display",

	// merge-driver
	"merge_driver.short": "Custom merge driver (called internally by git)",

	// store errors
	"store.gicket.not.found":   ".gicket directory not found. Run 'gicket init' to initialize.",
	"store.dir.create.failed":  "Failed to create directory: %w",
	"store.config.marshal":     "Failed to marshal config: %w",
	"store.config.write":       "Failed to write config: %w",
	"store.ticket.marshal":     "Failed to marshal ticket: %w",
	"store.ticket.save":        "Failed to save ticket: %w",
	"store.ticket.read":        "Failed to read ticket: %w",
	"store.ticket.parse":       "Failed to parse ticket: %w",
	"store.ticket.list":        "Failed to list tickets: %w",
	"store.ticket.search":      "Failed to search tickets: %w",
	"store.ticket.not.found":   "Ticket '%s' not found",
	"store.ticket.ambiguous":   "ID '%s' matches multiple tickets: %v",

	// git errors
	"git.repo.not.found":   "Git repository not found",
	"git.not.installed":    "git command not found. Please install Git.",
	"git.hook.exists":      "commit-msg hook already exists. Please merge manually: %s",
	"git.hook.install.fail":"Failed to install commit-msg hook: %w",
	"git.merge.driver.fail":"Failed to configure merge driver: %w",
	"git.gitattr.fail":     "Failed to configure .gitattributes: %w",
	"git.hook.remove.fail": "Failed to remove commit-msg hook: %w",
	"git.gitattr.update":   "Failed to update .gitattributes: %w",

	// merge errors
	"merge.ancestor.read": "Failed to read ancestor: %w",
	"merge.ours.read":     "Failed to read ours: %w",
	"merge.theirs.read":   "Failed to read theirs: %w",
	"merge.marshal":       "Failed to marshal merge result: %w",
	"merge.write":         "Failed to write merge result: %w",
	"merge.conflict":      "CONFLICT: Ticket %s has fields changed to different values in both branches",
}
