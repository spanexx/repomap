
// --- Types ---
export interface FileNode {
    path: string;
    language: string;
    importance: 'high' | 'medium' | 'low';
    rank: number;
    definitions?: string[];
    imports?: string[];
    token_count?: number;

    // --- Planning Features ---
    status?: 'existing' | 'planned' | 'modified';
    intent?: string;
    issues?: { type: string; description: string; severity: 'high' | 'medium' | 'low' }[];
    comments?: { user: string; text: string; }[];
}

export interface RepoMapData {
    repomap?: {
        files: FileNode[];
    };
    files?: FileNode[]; // Fallback
}

export interface StoryStep {
    title: string;
    text: string;
    filter: (f: FileNode) => boolean;
}

// --- Constants ---
export const COLORS = {
    high: '#ff7b72',
    medium: '#d29922',
    low: '#3fb950',
    pkg: '#8b949e',
    bg: '#161b22',
    edge: '#30363d',
    accent: '#58a6ff'
};

export const STORY_STEPS: StoryStep[] = [
    {
        title: "The Foundation",
        text: "Every great app starts with a solid core. These high-rank modules define the data structures and utilities that power everything else.",
        filter: (f: FileNode) => f.importance === 'high' && f.rank > 0.5
    },
    {
        title: "The Logic",
        text: "Building upon the core, these modules implement the business logic and algorithms. The complexity grows as the graph expands.",
        filter: (f: FileNode) => f.importance === 'medium' || (f.importance === 'high' && f.rank <= 0.5)
    },
    {
        title: "The Application",
        text: "Finally, the interface layers and entry points connect the logic to the user. The architecture is complete.",
        filter: (f: FileNode) => f.importance === 'low'
    }
];
