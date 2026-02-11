/**
 * Code Map: DependencyList – Structure tab showing definitions and import dependencies.
 * CID: 1.5.1-DependencyList
 */

import { Folder, Info, Code } from 'lucide-react';
import { type FileNode } from '../../types';

interface DependencyListProps {
    file: FileNode;
}

export const DependencyList: React.FC<DependencyListProps> = ({ file }) => {
    return (
        <div className="space-y-6">
            {/* Definitions */}
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

            {/* Dependencies */}
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
    );
};
