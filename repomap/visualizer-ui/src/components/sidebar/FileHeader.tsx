/**
 * Code Map: FileHeader – Sidebar header displaying file path, status badge,
 * language, and highlight/close actions.
 * CID: 1.5.1-FileHeader
 */

import { type FileNode } from '../../types';

interface FileHeaderProps {
    file: FileNode;
    onClose: () => void;
    isHighlighted?: boolean;
    onToggleHighlight?: () => void;
}

/** Status badge with color-coded styling based on file status. */
const StatusBadge = ({ status }: { status?: string }) => {
    if (!status || status === 'existing')
        return <span className="bg-[#30363d] text-xs px-2 py-0.5 rounded text-[#8b949e] border border-[#30363d]">Existing</span>;
    if (status === 'planned')
        return <span className="bg-[#a371f7]/20 text-[#a371f7] text-xs px-2 py-0.5 rounded border border-[#a371f7]/30">Planned</span>;
    return <span className="bg-[#d29922]/20 text-[#d29922] text-xs px-2 py-0.5 rounded border border-[#d29922]/30">Modified</span>;
};

export const FileHeader: React.FC<FileHeaderProps> = ({ file, onClose, isHighlighted, onToggleHighlight }) => {
    return (
        <div className="p-4 border-b border-[#30363d] bg-[#0d1117]/50">
            <div className="flex justify-between items-start mb-2">
                <div className="flex flex-col gap-1 overflow-hidden">
                    <div className="font-mono text-sm font-semibold text-[#58a6ff] truncate" title={file.path}>
                        {file.path}
                    </div>
                    <div className="flex gap-2">
                        <StatusBadge status={file.status} />
                        <span className="text-xs text-[#8b949e] font-mono">{file.language}</span>
                    </div>
                </div>
                <div className="flex items-center gap-2">
                    {onToggleHighlight && (
                        <button
                            onClick={onToggleHighlight}
                            className={`text-[10px] uppercase tracking-widest px-3 py-1 rounded-full border transition-all ${isHighlighted
                                ? 'bg-[#ff4444]/20 border-[#ff4444]/50 text-[#ff4444] hover:bg-[#ff4444]/30'
                                : 'bg-[#58a6ff]/10 border-[#58a6ff]/30 text-[#58a6ff] hover:bg-[#58a6ff]/20'
                                }`}
                        >
                            {isHighlighted ? 'Clear Highlight' : 'Highlight Node'}
                        </button>
                    )}
                    <button onClick={onClose} className="text-[#8b949e] hover:text-white transition-colors p-1">
                        <div className="sr-only">Close</div>
                        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                            <line x1="18" y1="6" x2="6" y2="18"></line>
                            <line x1="6" y1="6" x2="18" y2="18"></line>
                        </svg>
                    </button>
                </div>
            </div>
        </div>
    );
};
