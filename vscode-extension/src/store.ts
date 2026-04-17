import * as fs from 'fs';
import * as path from 'path';
import * as nodeCrypto from 'crypto';
import * as YAML from 'yaml';
import { Ticket, Comment } from './types';

const GICKET_DIR = '.gicket';
const ISSUES_DIR = 'issues';

/**
 * .gicket ディレクトリを探してリポジトリルートを返す
 */
export function findGicketRoot(startDir: string): string | undefined {
    let dir = startDir;
    while (true) {
        const gicketPath = path.join(dir, GICKET_DIR);
        if (fs.existsSync(gicketPath) && fs.statSync(gicketPath).isDirectory()) {
            return dir;
        }
        const parent = path.dirname(dir);
        if (parent === dir) {
            return undefined;
        }
        dir = parent;
    }
}

/**
 * issues ディレクトリのパスを返す
 */
function issuesDir(root: string): string {
    return path.join(root, GICKET_DIR, ISSUES_DIR);
}

/**
 * 全チケットを読み込む
 */
export function listTickets(root: string): Ticket[] {
    const dir = issuesDir(root);
    if (!fs.existsSync(dir)) {
        return [];
    }

    const entries = fs.readdirSync(dir);
    const tickets: Ticket[] = [];

    for (const entry of entries) {
        if (!entry.endsWith('.yml')) {
            continue;
        }
        try {
            const ticket = loadTicket(root, entry.replace('.yml', ''));
            if (ticket) {
                tickets.push(ticket);
            }
        } catch {
            // skip invalid files
        }
    }

    tickets.sort((a, b) => new Date(b.created).getTime() - new Date(a.created).getTime());
    return tickets;
}

/**
 * 指定IDのチケットを読み込む
 */
export function loadTicket(root: string, id: string): Ticket | undefined {
    const filePath = ticketPath(root, id);
    if (!fs.existsSync(filePath)) {
        return undefined;
    }
    const content = fs.readFileSync(filePath, 'utf-8');
    return YAML.parse(content) as Ticket;
}

/**
 * チケットを保存する
 */
export function saveTicket(root: string, ticket: Ticket): void {
    ticket.updated = new Date().toISOString();
    const filePath = ticketPath(root, ticket.id);
    const dir = issuesDir(root);
    if (!fs.existsSync(dir)) {
        fs.mkdirSync(dir, { recursive: true });
    }
    const content = YAML.stringify(ticket);
    fs.writeFileSync(filePath, content, 'utf-8');
}

/**
 * チケットを削除する
 */
export function deleteTicket(root: string, id: string): void {
    const filePath = ticketPath(root, id);
    if (fs.existsSync(filePath)) {
        fs.unlinkSync(filePath);
    }
}

/**
 * 新しいIDを生成する
 */
export function generateId(): string {
    const now = new Date();
    const pad2 = (n: number) => n.toString().padStart(2, '0');
    const timestamp = `${now.getFullYear()}${pad2(now.getMonth() + 1)}${pad2(now.getDate())}-${pad2(now.getHours())}${pad2(now.getMinutes())}${pad2(now.getSeconds())}`;
    const hex = nodeCrypto.randomBytes(3).toString('hex');
    return `${timestamp}-${hex}`;
}

/**
 * .gicket を初期化する
 */
export function initGicket(root: string): void {
    const dir = issuesDir(root);
    fs.mkdirSync(dir, { recursive: true });

    const configPath = path.join(root, GICKET_DIR, 'config.yml');
    if (!fs.existsSync(configPath)) {
        fs.writeFileSync(configPath, YAML.stringify({ version: '1' }), 'utf-8');
    }
}

/**
 * コメントを追加する
 */
export function addComment(root: string, ticketId: string, body: string, author: string): Ticket | undefined {
    const ticket = loadTicket(root, ticketId);
    if (!ticket) {
        return undefined;
    }

    const comment: Comment = {
        author,
        date: new Date().toISOString(),
        body,
    };

    if (!ticket.comments) {
        ticket.comments = [];
    }
    ticket.comments.push(comment);
    saveTicket(root, ticket);
    return ticket;
}

function ticketPath(root: string, id: string): string {
    return path.join(root, GICKET_DIR, ISSUES_DIR, `${id}.yml`);
}

/**
 * チケットYAMLファイルのURIを返す
 */
export function ticketUri(root: string, id: string): string {
    return ticketPath(root, id);
}
