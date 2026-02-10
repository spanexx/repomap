import { useState } from 'react';
import { motion } from 'framer-motion';
import { Folder, Info, AlertTriangle, MessageSquare, BookOpen, Code, Activity } from 'lucide-react';
import { type FileNode, COLORS } from '../types';

interface SidebarProps {
    file: FileNode;
    onClose: () => void;
}

type TabType = 'overview' | 'code' | 'plan';

export const Sidebar: React.FC<SidebarProps> = ({ file, onClose }) => {
    const [activeTab, setActiveTab] = useState<TabType>('overview');

    const getColor = (f: FileNode) => {
        if (f.importance === 'high') return COLORS.high;
        if (f.importance === 'medium') return COLORS.medium;
        return COLORS.low;
    };

    const StatusBadge = ({ status }: { status?: string }) => {
        if (!status || status === 'existing') return <span className="bg-[#30363d] text-xs px-2 py-0.5 rounded text-[#8b949e] border border-[#30363d]">Existing</span>;
        if (status === 'planned') return <span className="bg-[#a371f7]/20 text-[#a371f7] text-xs px-2 py-0.5 rounded border border-[#a371f7]/30">Planned</span>;
        return <span className="bg-[#d29922]/20 text-[#d29922] text-xs px-2 py-0.5 rounded border border-[#d29922]/30">Modified</span>;
    };

    return (
        <motion.aside
            initial={{ x: 400 }}
            animate={{ x: 0 }}
            exit={{ x: 400 }}
            transition={{ type: 'spring', damping: 20, stiffness: 100 }}
            className="absolute top-4 right-4 bottom-4 w-[400px] bg-[#161b22]/95 backdrop-blur-xl border border-[#30363d] rounded-2xl shadow-2xl flex flex-col z-40 overflow-hidden"
        >
            {/* Header */}
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
                    <button onClick={onClose} className="text-[#8b949e] hover:text-white transition-colors p-1">
                        <div className="sr-only">Close</div>
                        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
                    </button>
                </div>
            </div>

            {/* Tabs */}
            <div className="flex border-b border-[#30363d] bg-[#0d1117]/30">
                <button
                    onClick={() => setActiveTab('overview')}
                    className={`flex-1 py-3 text-xs font-medium uppercase tracking-wider transition-colors border-b-2 ${activeTab === 'overview' ? 'border-[#58a6ff] text-[#58a6ff] bg-[#58a6ff]/5' : 'border-transparent text-[#8b949e] hover:text-[#c9d1d9]'}`}
                >
                    Overview
                </button>
                <button
                    onClick={() => setActiveTab('code')}
                    className={`flex-1 py-3 text-xs font-medium uppercase tracking-wider transition-colors border-b-2 ${activeTab === 'code' ? 'border-[#58a6ff] text-[#58a6ff] bg-[#58a6ff]/5' : 'border-transparent text-[#8b949e] hover:text-[#c9d1d9]'}`}
                >
                    Structure
                </button>
                <button
                    onClick={() => setActiveTab('plan')}
                    className={`flex-1 py-3 text-xs font-medium uppercase tracking-wider transition-colors border-b-2 ${activeTab === 'plan' ? 'border-[#58a6ff] text-[#58a6ff] bg-[#58a6ff]/5' : 'border-transparent text-[#8b949e] hover:text-[#c9d1d9]'}`}
                >
                    Plan & Issues
                </button>
            </div>

            {/* Content */}
            <div className="flex-1 overflow-y-auto custom-scrollbar bg-[#0d1117]/20 p-4">

                {activeTab === 'overview' && (
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
                )}

                {activeTab === 'code' && (
                    <div className="space-y-6">
                        <div>
                            <h3 className="text-xs font-bold text-[#8b949e] uppercase tracking-wider mb-3 flex items-center gap-2">
                                <Folder size={14} /> Definitions
                            </h3>
                            <div className="space-y-2">
                                {file.definitions?.map((def, i) => (
                                    <pre key={i} className="bg-[#0d1117] border border-[#30363d] rounded-md p-2 text-xs font-mono text-[#d2a8ff] overflow-x-auto whitespace-pre-wrap border-l-2 border-l-[#58a6ff]">
                                        {def}
                                    </pre>
                                ))}
                                {!file.definitions?.length && <span className="text-sm text-[#8b949e] italic">No public definitions</span>}
                            </div>
                        </div>

                        <div>
                            <h3 className="text-xs font-bold text-[#8b949e] uppercase tracking-wider mb-3 flex items-center gap-2">
                                <Info size={14} /> Dependencies
                            </h3>
                            <div className="flex flex-wrap gap-2">
                                {file.imports?.map((imp, i) => (
                                    <span key={i} className="bg-[#30363d]/50 text-[#c9d1d9] px-2 py-1 rounded text-xs border border-[#30363d] flex items-center gap-1 hover:bg-[#30363d] transition-colors cursor-default">
                                        <Code size={10} className="text-[#8b949e]" /> {imp}
                                    </span>
                                ))}
                                {!file.imports?.length && <span className="text-sm text-[#8b949e] italic">No dependencies</span>}
                            </div>
                        </div>
                    </div>
                )}

                {activeTab === 'plan' && (
                    <div className="space-y-6">
                        {/* Issues */}
                        <div>
                            <h3 className="text-xs font-bold text-[#8b949e] uppercase tracking-wider mb-3 flex items-center gap-2">
                                <AlertTriangle size={14} /> Detected Issues
                            </h3>
                            <div className="space-y-2">
                                {file.issues?.map((issue, i) => (
                                    <div key={i} className={`p-3 rounded-lg border text-sm ${issue.severity === 'high' ? 'bg-red-900/10 border-red-900/30 text-red-400' :
                                            issue.severity === 'medium' ? 'bg-[#d29922]/10 border-[#d29922]/30 text-[#d29922]' :
                                                'bg-[#30363d]/30 border-[#30363d] text-[#8b949e]'
                                        }`}>
                                        <div className="flex items-center gap-2 font-semibold mb-1">
                                            <span className="capitalize">{issue.type.replace('_', ' ')}</span>
                                        </div>
                                        <p className="opacity-90">{issue.description}</p>
                                    </div>
                                ))}
                                {!file.issues?.length && (
                                    <div className="text-center py-6 border border-dashed border-[#30363d] rounded-lg text-[#8b949e]">
                                        <div className="text-sm">No issues detected</div>
                                        <div className="text-xs opacity-60 mt-1">Ready for planning</div>
                                    </div>
                                )}
                            </div>
                        </div>

                        {/* Comments */}
                        <div>
                            <h3 className="text-xs font-bold text-[#8b949e] uppercase tracking-wider mb-3 flex items-center gap-2">
                                <MessageSquare size={14} /> Discussion
                            </h3>
                            <div className="space-y-3">
                                {file.comments?.map((comment, i) => (
                                    <div key={i} className="bg-[#161b22] border border-[#30363d] rounded-lg p-3">
                                        <div className="flex justify-between items-center mb-2">
                                            <span className="text-xs font-bold text-[#58a6ff]">{comment.user}</span>
                                            <span className="text-[10px] text-[#8b949e]">Just now</span>
                                        </div>
                                        <p className="text-sm text-[#c9d1d9]">{comment.text}</p>
                                    </div>
                                ))}
                                <div className="bg-[#0d1117] border border-[#30363d] rounded-lg p-2 text-[#8b949e] text-sm hover:border-[#58a6ff] cursor-text transition-colors flex items-center gap-2">
                                    <MessageSquare size={14} /> Add a comment...
                                </div>
                            </div>
                        </div>
                    </div>
                )}

            </div>
        </motion.aside>
    );
};
