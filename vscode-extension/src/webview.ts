import * as vscode from 'vscode';
import { Ticket } from './types';

export class TicketWebview {
    private static panels: Map<string, vscode.WebviewPanel> = new Map();

    static show(ticket: Ticket, extensionUri: vscode.Uri, onUpdate: () => void): void {
        const existing = this.panels.get(ticket.id);
        if (existing) {
            existing.reveal();
            existing.webview.html = this.getHtml(ticket);
            return;
        }

        const panel = vscode.window.createWebviewPanel(
            'gicketTicket',
            `Ticket: ${ticket.title}`,
            vscode.ViewColumn.One,
            { enableScripts: false },
        );

        panel.webview.html = this.getHtml(ticket);
        this.panels.set(ticket.id, panel);

        panel.onDidDispose(() => {
            this.panels.delete(ticket.id);
        });
    }

    static update(ticket: Ticket): void {
        const panel = this.panels.get(ticket.id);
        if (panel) {
            panel.webview.html = this.getHtml(ticket);
            panel.title = `Ticket: ${ticket.title}`;
        }
    }

    private static getHtml(ticket: Ticket): string {
        const priorityColor = ticket.priority === 'high' ? '#e74c3c' :
                              ticket.priority === 'medium' ? '#f39c12' : '#27ae60';
        const statusColor = ticket.status === 'open' ? '#3498db' :
                            ticket.status === 'in-progress' ? '#f39c12' : '#95a5a6';

        const labelsHtml = ticket.labels && ticket.labels.length > 0
            ? ticket.labels.map(l => `<span class="label">${escapeHtml(l)}</span>`).join(' ')
            : '<span class="muted">None</span>';

        const commentsHtml = ticket.comments && ticket.comments.length > 0
            ? ticket.comments.map(c => `
                <div class="comment">
                    <div class="comment-header">
                        <strong>${escapeHtml(c.author)}</strong>
                        <span class="muted">${formatDate(c.date)}</span>
                    </div>
                    <div class="comment-body">${escapeHtml(c.body)}</div>
                </div>
            `).join('')
            : '<p class="muted">No comments</p>';

        return `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<style>
    body {
        font-family: var(--vscode-font-family);
        color: var(--vscode-foreground);
        background: var(--vscode-editor-background);
        padding: 16px;
        line-height: 1.6;
    }
    h1 { margin: 0 0 8px 0; font-size: 1.4em; }
    .meta { display: flex; gap: 12px; flex-wrap: wrap; margin-bottom: 16px; }
    .badge {
        display: inline-block;
        padding: 2px 10px;
        border-radius: 12px;
        font-size: 0.85em;
        font-weight: 600;
        color: #fff;
    }
    .field { margin-bottom: 12px; }
    .field-label {
        font-size: 0.85em;
        color: var(--vscode-descriptionForeground);
        margin-bottom: 2px;
    }
    .label {
        display: inline-block;
        background: var(--vscode-badge-background);
        color: var(--vscode-badge-foreground);
        padding: 1px 8px;
        border-radius: 10px;
        font-size: 0.8em;
        margin-right: 4px;
    }
    .muted { color: var(--vscode-descriptionForeground); }
    .description {
        background: var(--vscode-textBlockQuote-background);
        border-left: 3px solid var(--vscode-textBlockQuote-border);
        padding: 8px 12px;
        margin: 8px 0;
        white-space: pre-wrap;
    }
    hr {
        border: none;
        border-top: 1px solid var(--vscode-widget-border);
        margin: 16px 0;
    }
    .comment {
        background: var(--vscode-editor-inactiveSelectionBackground);
        border-radius: 6px;
        padding: 8px 12px;
        margin-bottom: 8px;
    }
    .comment-header {
        display: flex;
        justify-content: space-between;
        margin-bottom: 4px;
        font-size: 0.9em;
    }
    .comment-body { white-space: pre-wrap; }
    .id-text {
        font-family: var(--vscode-editor-font-family);
        font-size: 0.85em;
        color: var(--vscode-descriptionForeground);
    }
</style>
</head>
<body>
    <h1>${escapeHtml(ticket.title)}</h1>
    <div class="id-text">${escapeHtml(ticket.id)}</div>

    <div class="meta" style="margin-top: 12px">
        <span class="badge" style="background:${statusColor}">${escapeHtml(ticket.status)}</span>
        <span class="badge" style="background:${priorityColor}">${escapeHtml(ticket.priority)}</span>
    </div>

    <div class="field">
        <div class="field-label">Assignee</div>
        <div>${ticket.assignee ? escapeHtml(ticket.assignee) : '<span class="muted">Unassigned</span>'}</div>
    </div>

    <div class="field">
        <div class="field-label">Labels</div>
        <div>${labelsHtml}</div>
    </div>

    <div class="field">
        <div class="field-label">Author</div>
        <div>${escapeHtml(ticket.author)}</div>
    </div>

    <div class="field">
        <div class="field-label">Created</div>
        <div>${formatDate(ticket.created)}</div>
    </div>

    <div class="field">
        <div class="field-label">Updated</div>
        <div>${formatDate(ticket.updated)}</div>
    </div>

    ${ticket.description ? `
    <hr>
    <div class="field">
        <div class="field-label">Description</div>
        <div class="description">${escapeHtml(ticket.description)}</div>
    </div>
    ` : ''}

    <hr>
    <div class="field">
        <div class="field-label">Comments (${ticket.comments?.length || 0})</div>
        ${commentsHtml}
    </div>
</body>
</html>`;
    }
}

function escapeHtml(text: string): string {
    return text
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;');
}

function formatDate(dateStr: string): string {
    try {
        const d = new Date(dateStr);
        return d.toLocaleString();
    } catch {
        return dateStr;
    }
}
