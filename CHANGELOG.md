# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/), and this project adheres to [Semantic Versioning](https://semver.org/).

## [1.0.0] - 2026-04-17

### Added

- **CLI commands**: `init`, `new`, `list`, `show`, `edit`, `comment`, `close`, `reopen`, `search`, `stats`, `log`
- **Web UI**: Built-in single-page application with Kanban board, dashboard, filters, and light/dark theme (`gicket serve`)
- **REST API**: Full CRUD API for tickets and comments (`/api/tickets`)
- **Git integration**: commit-msg hook for ticket ID validation, custom 3-way merge driver for `.gicket/issues/*.yml`
- **VS Code extension**: Sidebar tree view, ticket detail webview, quick commands, auto-refresh
- **Multilingual support**: English (default) and Japanese — switchable via `GICKET_LANG` environment variable
- **JSON output**: `--json` flag on `list`, `show`, `search`, `stats` commands
- **Search**: Full-text search across title, description, comments, labels, and assignee
- **Reopen**: Reopen closed tickets with `gicket reopen`
- **Statistics**: Ticket summary by status and priority with `gicket stats`
- **Input validation**: Status and priority values are validated on `new` and `edit`
- **Configuration**: Project-level (`.gicket/config.yml`) and user-level (`~/.config/gicket/config.yml`) config files
- **Short ID lookup**: Reference tickets by unique prefix instead of full ID
- **GitHub Actions**: Automated multi-platform binary releases on tag push (Windows/macOS/Linux, amd64/arm64)
