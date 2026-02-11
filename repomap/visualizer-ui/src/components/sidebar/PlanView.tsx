/**
 * Code Map: PlanView – Plan & Issues tab showing detected issues and discussion comments.
 * CID: 1.5.1-PlanView
 */

import { AlertTriangle, MessageSquare } from 'lucide-react';
import { type FileNode } from '../../types';
import { Chat } from '../Chat';

interface PlanViewProps {
    file: FileNode;
}

export const PlanView: React.FC<PlanViewProps> = ({ file }) => {
    return (
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

            {/* Discussion (Chat) */}
            <div className="pt-4 border-t border-[#30363d]">
                <h3 className="text-xs font-bold text-[#8b949e] uppercase tracking-wider mb-3 flex items-center gap-2">
                    <MessageSquare size={14} /> Discussion
                </h3>
                <div className="h-[500px] border border-[#30363d] rounded-lg overflow-hidden relative">
                    <Chat
                        files={[file]}
                        selectedNode={file.path}
                        viewMode="plan"
                        embedded={true}
                    />
                </div>
            </div>
        </div>
    );
};
