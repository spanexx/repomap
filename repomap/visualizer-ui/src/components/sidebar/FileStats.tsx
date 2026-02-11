/**
 * Code Map: FileStats – Overview tab content showing intent and metrics grid.
 * CID: 1.5.1-FileStats
 */

import { BookOpen, Activity } from 'lucide-react';
import { type FileNode, COLORS } from '../../types';

interface FileStatsProps {
    file: FileNode;
}

const getColor = (f: FileNode) => {
    if (f.importance === 'high') return COLORS.high;
    if (f.importance === 'medium') return COLORS.medium;
    return COLORS.low;
};

export const FileStats: React.FC<FileStatsProps> = ({ file }) => {
    return (
        <div className="space-y-6">
            {/* Intent Section */}
            <div className="bg-[#161b22] border border-[#30363d] rounded-xl p-4">
                <h3 className="text-xs font-bold text-[#8b949e] uppercase tracking-wider mb-2 flex items-center gap-2">
                    <BookOpen size={14} /> Intent
                </h3>
                <p className="text-sm text-[#c9d1d9] leading-relaxed">
                    {file.intent || "No intent description provided for this module."}
                </p>
            </div>

            {/* Metrics Grid */}
            <div className="grid grid-cols-2 gap-2">
                <div className="bg-[#161b22] p-3 rounded-lg border border-[#30363d]">
                    <span className="text-[10px] uppercase tracking-wider text-[#8b949e] block mb-1">Importance</span>
                    <span className="font-bold text-lg" style={{ color: getColor(file) }}>
                        {file.importance.toUpperCase()}
                    </span>
                </div>
                <div className="bg-[#161b22] p-3 rounded-lg border border-[#30363d]">
                    <span className="text-[10px] uppercase tracking-wider text-[#8b949e] block mb-1">Rank</span>
                    <span className="font-mono font-bold text-lg">
                        {file.rank ? file.rank.toFixed(3) : '0.00'}
                    </span>
                </div>
                <div className="bg-[#161b22] p-3 rounded-lg border border-[#30363d]">
                    <span className="text-[10px] uppercase tracking-wider text-[#8b949e] block mb-1">Tokens</span>
                    <span className="font-bold text-lg">{file.token_count || 0}</span>
                </div>
                <div className="bg-[#161b22] p-3 rounded-lg border border-[#30363d]">
                    <span className="text-[10px] uppercase tracking-wider text-[#8b949e] block mb-1">Activity</span>
                    <span className="font-bold text-lg text-[#3fb950] flex items-center gap-1">
                        <Activity size={14} /> Low
                    </span>
                </div>
            </div>
        </div>
    );
};
