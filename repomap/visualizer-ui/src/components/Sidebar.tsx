/**
 * Code Map: Sidebar – Composed file detail panel using modular sub-components.
 * CID: 1.5.1-Sidebar
 */

import { useState } from 'react';
import { motion } from 'framer-motion';
import { type FileNode } from '../types';
import { FileHeader } from './sidebar/FileHeader';
import { FileStats } from './sidebar/FileStats';
import { DependencyList } from './sidebar/DependencyList';
import { PlanView } from './sidebar/PlanView';

interface SidebarProps {
    file: FileNode;
    onClose: () => void;
    isHighlighted?: boolean;
    onToggleHighlight?: () => void;
}

type TabType = 'overview' | 'code' | 'plan';

export const Sidebar: React.FC<SidebarProps> = ({ file, onClose, isHighlighted, onToggleHighlight }) => {
    const [activeTab, setActiveTab] = useState<TabType>('overview');

    return (
        <motion.aside
            initial={{ x: '100%' }}
            animate={{ x: 0 }}
            exit={{ x: '100%' }}
            transition={{ type: 'spring', damping: 20, stiffness: 100 }}
            className="fixed right-0 bottom-0 top-[60px] sm:top-4 sm:right-4 sm:bottom-4 w-full sm:w-[400px] z-50 sm:z-40 bg-[#161b22]/95 backdrop-blur-xl border-l sm:border border-[#30363d] sm:rounded-2xl shadow-2xl flex flex-col overflow-hidden"
        >
            <FileHeader
                file={file}
                onClose={onClose}
                isHighlighted={isHighlighted}
                onToggleHighlight={onToggleHighlight}
            />

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
                {activeTab === 'overview' && <FileStats file={file} />}
                {activeTab === 'code' && <DependencyList file={file} />}
                {activeTab === 'plan' && <PlanView file={file} />}
            </div>
        </motion.aside>
    );
};
