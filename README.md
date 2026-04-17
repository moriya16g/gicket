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

## Roadmap

- [ ] **Phase 1**: CLI core (current)
- [ ] **Phase 2**: REST API + Web UI (`gicket serve`)
- [ ] **Phase 3**: Kanban board, real-time filters, dashboard
- [ ] **Phase 4**: VS Code extension

## Similar Projects

| Project | Language | Approach |
|---------|----------|----------|
| [git-bug](https://github.com/git-bug/git-bug) | Go | Stores issues as Git objects (not files) |
| [git-issue](https://github.com/dspinellis/git-issue) | Shell | Text files in `.issues/` directory |
| [SIT](https://github.com/sit-fyi/sit) | Rust | Serverless Information Tracker |
| [Bugs Everywhere](http://www.bugseverywhere.org/) | Python | Multi-VCS support |

**gicket** differentiates itself by combining human-readable YAML files with a rich Web UI (planned), all delivered as a single binary.

## License

MIT License. See [LICENSE](LICENSE) for details.
