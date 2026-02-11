import React, { useState } from 'react';
import type { FileNode } from '../types';
import { ArrowUpDown, FileText, AlertTriangle, MessageSquare } from 'lucide-react';

interface TableViewProps {
    files: FileNode[];
    onFileSelect: (file: FileNode) => void;
}

type SortField = 'path' | 'intent' | 'importance' | 'tokens';
type SortOrder = 'asc' | 'desc';

export const TableView: React.FC<TableViewProps> = ({ files, onFileSelect }) => {
    const [sortField, setSortField] = useState<SortField>('importance');
    const [sortOrder, setSortOrder] = useState<SortOrder>('desc');
    const [filter, setFilter] = useState('');

    const handleSort = (field: SortField) => {
        if (sortField === field) {
            setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
        } else {
            setSortField(field);
            setSortOrder('asc');
        }
    };

    const getSortValue = (file: FileNode, field: SortField): string | number => {
        if (field === 'tokens') return file.token_count || 0;
        if (field === 'importance') {
            const rank = { high: 3, medium: 2, low: 1 };
            return rank[file.importance as keyof typeof rank] || 0;
        }
        return (file[field as keyof FileNode] as string) || '';
    };

    const sortedFiles = [...files]
        .filter(f => f.path.toLowerCase().includes(filter.toLowerCase()) || f.intent?.toLowerCase().includes(filter.toLowerCase()))
        .sort((a, b) => {
            const valA = getSortValue(a, sortField);
            const valB = getSortValue(b, sortField);

            if (valA < valB) return sortOrder === 'asc' ? -1 : 1;
            if (valA > valB) return sortOrder === 'asc' ? 1 : -1;
            return 0;
        });

    return (
        <div className="w-full h-full overflow-hidden flex flex-col pt-24 px-4 sm:px-8 pb-4">
            <div className="flex justify-between items-center mb-6">
                <h2 className="text-2xl font-semibold text-white">Repository Overview</h2>
                <input
                    type="text"
                    placeholder="Filter files..."
                    value={filter}
                    onChange={(e) => setFilter(e.target.value)}
                    className="bg-[#161b22] border border-[#30363d] rounded-lg px-4 py-2 text-sm text-white focus:outline-none focus:border-[#58a6ff] transition-colors w-64"
                />
            </div>

            <div className="flex-1 overflow-auto bg-[#161b22] border border-[#30363d] rounded-xl shadow-xl custom-scrollbar">
                <table className="w-full text-left border-collapse">
                    <thead className="bg-[#0d1117] sticky top-0 z-10">
                        <tr>
                            <th className="p-4 border-b border-[#30363d] font-semibold text-[#8b949e] cursor-pointer hover:text-white transition-colors" onClick={() => handleSort('path')}>
                                <div className="flex items-center gap-2">Path <ArrowUpDown size={14} /></div>
                            </th>
                            <th className="p-4 border-b border-[#30363d] font-semibold text-[#8b949e] cursor-pointer hover:text-white transition-colors" onClick={() => handleSort('intent')}>
                                <div className="flex items-center gap-2">Intent <ArrowUpDown size={14} /></div>
                            </th>
                            <th className="p-4 border-b border-[#30363d] font-semibold text-[#8b949e] cursor-pointer hover:text-white transition-colors" onClick={() => handleSort('importance')}>
                                <div className="flex items-center gap-2">Importance <ArrowUpDown size={14} /></div>
                            </th>
                            <th className="p-4 border-b border-[#30363d] font-semibold text-[#8b949e] cursor-pointer hover:text-white transition-colors" onClick={() => handleSort('tokens')}>
                                <div className="flex items-center gap-2">Tokens <ArrowUpDown size={14} /></div>
                            </th>
                            <th className="p-4 border-b border-[#30363d] font-semibold text-[#8b949e]">Issues</th>
                        </tr>
                    </thead>
                    <tbody>
                        {sortedFiles.map(file => (
                            <tr
                                key={file.path}
                                onClick={() => onFileSelect(file)}
                                className="border-b border-[#30363d] hover:bg-[#1f242c] cursor-pointer transition-colors group"
                            >
                                <td className="p-4 text-sm font-medium text-[#c9d1d9] font-mono group-hover:text-[#58a6ff]">
                                    <div className="flex items-center gap-2">
                                        <FileText size={14} className="text-[#8b949e]" />
                                        {file.path}
                                    </div>
                                </td>
                                <td className="p-4 text-sm text-[#8b949e]">
                                    {file.intent ? (
                                        <span className="px-2 py-1 rounded-full bg-[#30363d]/50 border border-[#30363d] text-xs">
                                            {file.intent}
                                        </span>
                                    ) : (
                                        <span className="text-[#484f58] italic">Unassigned</span>
                                    )}
                                </td>
                                <td className="p-4 text-sm">
                                    <span className={`px-2 py-1 rounded-full text-xs font-semibold
                                        ${file.importance === 'high' ? 'bg-red-900/30 text-red-400 border border-red-900/50' :
                                            file.importance === 'medium' ? 'bg-yellow-900/30 text-yellow-400 border border-yellow-900/50' :
                                                'bg-[#30363d] text-[#8b949e] border border-[#30363d]'}`}>
                                        {file.importance.toUpperCase()}
                                    </span>
                                </td>
                                <td className="p-4 text-sm text-[#8b949e] font-mono">
                                    {file.token_count?.toLocaleString()}
                                </td>
                                <td className="p-4 text-sm">
                                    <div className="flex gap-2">
                                        {file.issues && file.issues.length > 0 && (
                                            <span className="flex items-center gap-1 text-yellow-500 bg-yellow-900/20 px-2 py-1 rounded-full text-xs">
                                                <AlertTriangle size={12} /> {file.issues.length}
                                            </span>
                                        )}
                                        {file.comments && file.comments.length > 0 && (
                                            <span className="flex items-center gap-1 text-[#58a6ff] bg-[#58a6ff]/20 px-2 py-1 rounded-full text-xs">
                                                <MessageSquare size={12} /> {file.comments.length}
                                            </span>
                                        )}
                                    </div>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
};
