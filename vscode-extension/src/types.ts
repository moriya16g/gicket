export interface Ticket {
    id: string;
    title: string;
    status: 'open' | 'in-progress' | 'closed';
    priority: 'low' | 'medium' | 'high';
    assignee?: string;
    labels?: string[];
    created: string;
    updated: string;
    author: string;
    description?: string;
    comments?: Comment[];
}

export interface Comment {
    author: string;
    date: string;
    body: string;
}

export const STATUS_LABELS: Record<string, string> = {
    'open': 'Open',
    'in-progress': 'In Progress',
    'closed': 'Closed',
};

export const PRIORITY_ICONS: Record<string, string> = {
    'high': '🔴',
    'medium': '🟡',
    'low': '🟢',
};
