/**
 * Code Map: useSelection – Custom hook for managing file selection and node highlighting.
 * CID: 1.5.3-useSelection
 */

import { useState, useCallback } from 'react';
import { type FileNode } from '../types';

export function useSelection() {
    const [selectedFile, setSelectedFile] = useState<FileNode | null>(null);
    const [highlightedNodes, setHighlightedNodes] = useState<string[]>([]);

    const toggleHighlight = useCallback((path: string) => {
        setHighlightedNodes(prev =>
            prev.includes(path)
                ? prev.filter(p => p !== path)
                : [...prev, path]
        );
    }, []);

    return {
        selectedFile,
        setSelectedFile,
        highlightedNodes,
        setHighlightedNodes,
        toggleHighlight
    };
}
