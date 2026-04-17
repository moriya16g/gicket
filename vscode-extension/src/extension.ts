import * as vscode from 'vscode';
import { TicketTreeProvider, TicketTreeItem } from './treeView';
import { TicketWebview } from './webview';
import * as store from './store';
import { Ticket } from './types';

let treeProvider: TicketTreeProvider;

export function activate(context: vscode.ExtensionContext) {
    treeProvider = new TicketTreeProvider();
    vscode.window.registerTreeDataProvider('gicketTickets', treeProvider);

    // Watch for file changes in .gicket directory
    const watcher = vscode.workspace.createFileSystemWatcher('**/.gicket/issues/*.yml');
    watcher.onDidChange(() => treeProvider.refresh());
    watcher.onDidCreate(() => treeProvider.refresh());
    watcher.onDidDelete(() => treeProvider.refresh());
    context.subscriptions.push(watcher);

    context.subscriptions.push(
        vscode.commands.registerCommand('gicket.refresh', () => {
            treeProvider.refresh();
        }),

        vscode.commands.registerCommand('gicket.init', async () => {
            const folders = vscode.workspace.workspaceFolders;
            if (!folders || folders.length === 0) {
                vscode.window.showErrorMessage('ワークスペースフォルダが開かれていません');
                return;
            }
            const root = folders[0].uri.fsPath;
            store.initGicket(root);
            treeProvider.refresh();
            vscode.window.showInformationMessage('gicket を初期化しました (.gicket/)');
        }),

        vscode.commands.registerCommand('gicket.newTicket', async () => {
            const root = getRoot();
            if (!root) { return; }

            const title = await vscode.window.showInputBox({
                prompt: 'チケットのタイトル',
                placeHolder: 'Fix login validation',
            });
            if (!title) { return; }

            const priority = await vscode.window.showQuickPick(
                ['medium', 'high', 'low'],
                { placeHolder: '優先度を選択' },
            );
            if (!priority) { return; }

            const assignee = await vscode.window.showInputBox({
                prompt: 'アサイン先（省略可）',
                placeHolder: 'user@example.com',
            });

            const labelsStr = await vscode.window.showInputBox({
                prompt: 'ラベル（カンマ区切り、省略可）',
                placeHolder: 'bug, frontend',
            });

            const description = await vscode.window.showInputBox({
                prompt: '説明（省略可）',
                placeHolder: 'Detailed description...',
            });

            const now = new Date().toISOString();
            const ticket: Ticket = {
                id: store.generateId(),
                title,
                status: 'open',
                priority: priority as Ticket['priority'],
                assignee: assignee || undefined,
                labels: labelsStr ? labelsStr.split(',').map(l => l.trim()).filter(l => l) : undefined,
                created: now,
                updated: now,
                author: getGitUser(),
                description: description || undefined,
            };

            store.saveTicket(root, ticket);
            treeProvider.refresh();
            vscode.window.showInformationMessage(`チケットを作成しました: ${ticket.id}`);
        }),

        vscode.commands.registerCommand('gicket.showTicket', (item: TicketTreeItem) => {
            if (!item.ticket) { return; }
            TicketWebview.show(item.ticket, context.extensionUri, () => treeProvider.refresh());
        }),

        vscode.commands.registerCommand('gicket.editTicket', async (item: TicketTreeItem) => {
            if (!item.ticket) { return; }
            const root = getRoot();
            if (!root) { return; }

            const ticket = store.loadTicket(root, item.ticket.id);
            if (!ticket) { return; }

            const field = await vscode.window.showQuickPick(
                [
                    { label: 'Title', description: ticket.title },
                    { label: 'Status', description: ticket.status },
                    { label: 'Priority', description: ticket.priority },
                    { label: 'Assignee', description: ticket.assignee || '(none)' },
                    { label: 'Labels', description: ticket.labels?.join(', ') || '(none)' },
                    { label: 'Description', description: ticket.description ? '(set)' : '(none)' },
                ],
                { placeHolder: '編集するフィールドを選択' },
            );
            if (!field) { return; }

            switch (field.label) {
                case 'Title': {
                    const v = await vscode.window.showInputBox({ value: ticket.title, prompt: 'タイトル' });
                    if (v !== undefined) { ticket.title = v; }
                    break;
                }
                case 'Status': {
                    const v = await vscode.window.showQuickPick(['open', 'in-progress', 'closed'], { placeHolder: 'ステータス' });
                    if (v) { ticket.status = v as Ticket['status']; }
                    break;
                }
                case 'Priority': {
                    const v = await vscode.window.showQuickPick(['low', 'medium', 'high'], { placeHolder: '優先度' });
                    if (v) { ticket.priority = v as Ticket['priority']; }
                    break;
                }
                case 'Assignee': {
                    const v = await vscode.window.showInputBox({ value: ticket.assignee || '', prompt: 'アサイン先' });
                    if (v !== undefined) { ticket.assignee = v || undefined; }
                    break;
                }
                case 'Labels': {
                    const v = await vscode.window.showInputBox({ value: ticket.labels?.join(', ') || '', prompt: 'ラベル（カンマ区切り）' });
                    if (v !== undefined) { ticket.labels = v ? v.split(',').map(l => l.trim()).filter(l => l) : undefined; }
                    break;
                }
                case 'Description': {
                    const v = await vscode.window.showInputBox({ value: ticket.description || '', prompt: '説明' });
                    if (v !== undefined) { ticket.description = v || undefined; }
                    break;
                }
            }

            store.saveTicket(root, ticket);
            treeProvider.refresh();
            TicketWebview.update(ticket);
            vscode.window.showInformationMessage(`チケットを更新しました: ${ticket.id}`);
        }),

        vscode.commands.registerCommand('gicket.closeTicket', async (item: TicketTreeItem) => {
            if (!item.ticket) { return; }
            const root = getRoot();
            if (!root) { return; }

            const ticket = store.loadTicket(root, item.ticket.id);
            if (!ticket) { return; }

            ticket.status = 'closed';
            store.saveTicket(root, ticket);
            treeProvider.refresh();
            TicketWebview.update(ticket);
            vscode.window.showInformationMessage(`チケットをクローズしました: ${ticket.id}`);
        }),

        vscode.commands.registerCommand('gicket.reopenTicket', async (item: TicketTreeItem) => {
            if (!item.ticket) { return; }
            const root = getRoot();
            if (!root) { return; }

            const ticket = store.loadTicket(root, item.ticket.id);
            if (!ticket) { return; }

            ticket.status = 'open';
            store.saveTicket(root, ticket);
            treeProvider.refresh();
            TicketWebview.update(ticket);
            vscode.window.showInformationMessage(`チケットを再オープンしました: ${ticket.id}`);
        }),

        vscode.commands.registerCommand('gicket.addComment', async (item: TicketTreeItem) => {
            if (!item.ticket) { return; }
            const root = getRoot();
            if (!root) { return; }

            const body = await vscode.window.showInputBox({
                prompt: 'コメント',
                placeHolder: 'Write a comment...',
            });
            if (!body) { return; }

            const ticket = store.addComment(root, item.ticket.id, body, getGitUser());
            if (ticket) {
                treeProvider.refresh();
                TicketWebview.update(ticket);
                vscode.window.showInformationMessage('コメントを追加しました');
            }
        }),

        vscode.commands.registerCommand('gicket.openYaml', (item: TicketTreeItem) => {
            if (!item.ticket) { return; }
            const root = getRoot();
            if (!root) { return; }

            const uri = vscode.Uri.file(store.ticketUri(root, item.ticket.id));
            vscode.window.showTextDocument(uri);
        }),
    );
}

function getRoot(): string | undefined {
    const folders = vscode.workspace.workspaceFolders;
    if (!folders || folders.length === 0) {
        vscode.window.showErrorMessage('ワークスペースフォルダが開かれていません');
        return undefined;
    }
    const root = store.findGicketRoot(folders[0].uri.fsPath);
    if (!root) {
        vscode.window.showErrorMessage('.gicket が見つかりません。gicket init を実行してください');
        return undefined;
    }
    return root;
}

function getGitUser(): string {
    try {
        const cp = require('child_process');
        const name = cp.execSync('git config user.name', { encoding: 'utf-8' }).trim();
        const email = cp.execSync('git config user.email', { encoding: 'utf-8' }).trim();
        if (name && email) {
            return `${name} <${email}>`;
        }
        return name || 'vscode-user';
    } catch {
        return 'vscode-user';
    }
}

export function deactivate() {}
