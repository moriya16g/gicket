# gicket

**A distributed ticket management tool embedded in Git repositories.**

[日本語](README_ja.md)

gicket manages tickets as human-readable YAML text files inside a Git repository. No web server, no database, no vendor lock-in — just Git.

## Features

- **Text-based**: Tickets are stored as plain YAML files that anyone can read and edit
- **Distributed**: Share tickets using standard `git push` / `git pull`
- **Zero infrastructure**: No server, no database — works entirely on the filesystem
- **Vendor-neutral**: Your data lives in your repository, not on someone else's platform
- **Single binary**: One executable, no dependencies, cross-platform (Windows / macOS / Linux)
- **Short ID lookup**: Reference tickets by unique prefix instead of full ID
- **Built-in Web UI**: Rich browser-based interface with Kanban board, filters, and light/dark theme
- **REST API**: Full HTTP API for integration with external tools
- **VS Code extension**: Manage tickets directly from your editor
- **Multilingual**: English (default) and Japanese UI — switchable via environment variable

## Quick Start

### Install

```bash
go install github.com/gicket/gicket@latest
```

Or build from source:

```bash
git clone https://github.com/gicket/gicket.git
cd gicket
go build -o gicket .
```

### Usage

```bash
# Initialize gicket in your repository
gicket init

# Create a new ticket
gicket new -t "Fix login validation" -p high -l bug,frontend

# List open tickets
gicket list

# List all tickets (including closed)
gicket list --all

# Show ticket details (full or prefix ID)
gicket show 20260416-200633-709268
gicket show 20260416-20          # prefix match

# Edit a ticket
gicket edit <id> -s in-progress -a "dev@example.com"
gicket edit <id> -d "Detailed description here"

# Add a comment
gicket comment <id> -m "Working on this now"

# Close a ticket
gicket close <id>

# Start Web UI (default: http://localhost:8080)
gicket serve
gicket serve -p 3000   # custom port
```

### Share with your team

Since tickets are plain files tracked by Git:

```bash
git add .gicket/
git commit -m "Add new tickets"
git push
```

Other developers simply `git pull` to receive ticket updates.

## Data Format

Tickets are stored in `.gicket/issues/` as YAML files:

```yaml
id: 20260416-200633-709268
title: Fix login validation
status: open
priority: high
assignee: dev@example.com
labels:
    - bug
    - frontend
created: 2026-04-16T20:06:33+09:00
updated: 2026-04-16T20:07:15+09:00
author: John Doe <john@example.com>
description: |
    Email format validation is missing on the login form.
comments:
    - author: Jane Smith <jane@example.com>
      date: 2026-04-16T21:00:00+09:00
      body: I'll handle this.
```

## Directory Structure

```
your-project/
├── .gicket/
│   ├── config.yml        # Project configuration
│   └── issues/
│       ├── 20260416-200633-709268.yml
│       ├── 20260416-200633-f16bab.yml
│       └── ...
├── src/
└── ...
```

## Commands

| Command | Description |
|---------|-------------|
| `gicket init` | Initialize gicket in the current directory |
| `gicket new` | Create a new ticket |
| `gicket list` | List open tickets (`--all` for all statuses) |
| `gicket show <id>` | Show ticket details |
| `gicket edit <id>` | Edit ticket fields |
| `gicket comment <id>` | Add a comment to a ticket |
| `gicket close <id>` | Close a ticket |
| `gicket serve` | Start Web UI server (`-p` for port, default 8080) |
| `gicket hook install` | Install Git hooks and custom merge driver |
| `gicket hook uninstall` | Uninstall Git hooks and merge driver |
| `gicket log <id>` | Show Git commits related to a ticket |

## Web UI

Run `gicket serve` to launch a browser-based interface:

- **Dashboard**: Ticket count cards (Open / In Progress / Closed / Total)
- **List & Kanban views**: Toggle between table list and drag-free Kanban board
- **Filters**: Status tabs + full-text search
- **Ticket operations**: Create, edit, close, and comment — all from the browser
- **Light / Dark theme**: Toggle in the header, preference saved in localStorage

The Web UI is embedded in the binary via `go:embed`, so no extra files are needed.

## REST API

`gicket serve` also exposes a REST API:

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/tickets` | List all tickets |
| `POST` | `/api/tickets` | Create a new ticket |
| `GET` | `/api/tickets/{id}` | Get a ticket by ID |
| `PUT` | `/api/tickets/{id}` | Update a ticket |
| `DELETE` | `/api/tickets/{id}` | Delete a ticket |
| `POST` | `/api/tickets/{id}/comments` | Add a comment |

## Git Integration

### Hooks & Merge Driver

```bash
# Install Git hooks and merge driver
gicket hook install

# Uninstall
gicket hook uninstall
```

`gicket hook install` sets up:

- **commit-msg hook**: Validates ticket ID references in commit messages (pattern: `gicket:<ticket-id>`). Set `GICKET_HOOK_REQUIRE_ID=1` to make it mandatory.
- **Custom merge driver**: Automatically resolves merge conflicts in `.gicket/issues/*.yml` using smart 3-way merge:
  - Single-field changes: applies the change
  - Comments: combines all comments from both branches (deduplication + chronological sort)
  - Labels: set union with deletion tracking
  - Status conflicts: prefers `closed` > `in-progress` > `open`
- **.gitattributes**: Configures the merge driver for ticket files

### Commit Log

```bash
# Show commits related to a ticket
gicket log <id>
gicket log <id> -n 20   # limit to 20 results
```

Searches for commits that mention the ticket ID in their message or modify the ticket file.

## VS Code Extension

The `vscode-extension/` directory contains a VS Code extension that provides:

- **Sidebar tree view**: Tickets grouped by status (Open / In Progress / Closed) with priority icons
- **Ticket detail panel**: Rich webview showing all fields, description, and comments
- **Quick commands**: Create, edit, close, reopen tickets and add comments via input prompts
- **YAML file access**: Open the raw YAML file for any ticket
- **Auto-refresh**: File watcher detects changes in `.gicket/issues/` automatically
- **Zero dependencies at runtime**: Reads/writes YAML files directly — no gicket CLI needed

### Install from source

```bash
cd vscode-extension
npm install
npm run compile
# Then press F5 in VS Code to launch Extension Development Host
```

The extension activates automatically when a workspace contains a `.gicket` directory.

## Language / i18n

gicket supports English (default) and Japanese. Set the language with an environment variable:

```bash
# English (default)
gicket list

# Japanese
GICKET_LANG=ja gicket list

# Or set system-wide
export GICKET_LANG=ja
```

Detection priority: `GICKET_LANG` > `LANG` > English.

## Roadmap

- [x] **Phase 1**: CLI core
- [x] **Phase 2**: REST API + Web UI (`gicket serve`)
- [x] **Phase 3**: Git integration (hooks, merge conflict resolution)
- [x] **Phase 4**: VS Code extension

## Similar Projects

| Project | Language | Approach |
|---------|----------|----------|
| [git-bug](https://github.com/git-bug/git-bug) | Go | Stores issues as Git objects (not files) |
| [git-issue](https://github.com/dspinellis/git-issue) | Shell | Text files in `.issues/` directory |
| [SIT](https://github.com/sit-fyi/sit) | Rust | Serverless Information Tracker |
| [Bugs Everywhere](http://www.bugseverywhere.org/) | Python | Multi-VCS support |

**gicket** differentiates itself by combining human-readable YAML files with a built-in Web UI, all delivered as a single binary.

## License

MIT License. See [LICENSE](LICENSE) for details.
