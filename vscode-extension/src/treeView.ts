import * as vscode from 'vscode';
import { Ticket, PRIORITY_ICONS, STATUS_LABELS } from './types';
import * as store from './store';

export class TicketTreeProvider implements vscode.TreeDataProvider<TicketTreeItem> {
    private _onDidChangeTreeData = new vscode.EventEmitter<TicketTreeItem | undefined>();
    readonly onDidChangeTreeData = this._onDidChangeTreeData.event;

    private root: string | undefined;

    constructor() {
        this.detectRoot();
    }

    refresh(): void {
        this.detectRoot();
        this._onDidChangeTreeData.fire(undefined);
    }

    private detectRoot(): void {
        const folders = vscode.workspace.workspaceFolders;
        if (folders && folders.length > 0) {
            this.root = store.findGicketRoot(folders[0].uri.fsPath);
        }
    }

    getTreeItem(element: TicketTreeItem): vscode.TreeItem {
        return element;
    }

    getChildren(element?: TicketTreeItem): TicketTreeItem[] {
        if (element) {
            return [];
        }

        if (!this.root) {
            return [];
        }

        const tickets = store.listTickets(this.root);

        // Group by status
        const groups: Record<string, Ticket[]> = {
            'open': [],
            'in-progress': [],
            'closed': [],
        };

        for (const ticket of tickets) {
            const status = ticket.status || 'open';
            if (groups[status]) {
                groups[status].push(ticket);
            }
        }

        const items: TicketTreeItem[] = [];

        for (const status of ['open', 'in-progress', 'closed']) {
            const group = groups[status];
            if (group.length === 0) {
                continue;
            }

            const groupItem = new TicketTreeItem(
                `${STATUS_LABELS[status]} (${group.length})`,
                vscode.TreeItemCollapsibleState.Expanded,
            );
            groupItem.contextValue = 'group';
            items.push(groupItem);

            for (const ticket of group) {
                items.push(TicketTreeItem.fromTicket(ticket));
            }
        }

        return items;
    }
}

export class TicketTreeItem extends vscode.TreeItem {
    ticket?: Ticket;

    static fromTicket(ticket: Ticket): TicketTreeItem {
        const icon = PRIORITY_ICONS[ticket.priority] || '⚪';
        const label = `${icon} ${ticket.title}`;
        const item = new TicketTreeItem(label, vscode.TreeItemCollapsibleState.None);
        item.ticket = ticket;
        item.contextValue = 'ticket';
        item.tooltip = new vscode.MarkdownString(
            `**${ticket.title}**\n\n` +
            `Status: ${ticket.status}  \n` +
            `Priority: ${ticket.priority}  \n` +
            (ticket.assignee ? `Assignee: ${ticket.assignee}  \n` : '') +
            (ticket.labels && ticket.labels.length > 0 ? `Labels: ${ticket.labels.join(', ')}  \n` : '') +
            `ID: \`${ticket.id}\``
        );
        item.command = {
            command: 'gicket.showTicket',
            title: 'Show Ticket',
            arguments: [item],
        };
        return item;
    }
}
