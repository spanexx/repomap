import { Share2, Upload, Layout, Play, Pause } from 'lucide-react';

interface HeaderProps {
    mode: 'cluster' | 'flow' | 'rank';
    isStoryPlaying: boolean;
    onUpload: (e: React.ChangeEvent<HTMLInputElement>) => void;
    onSetMode: (mode: 'cluster' | 'flow' | 'rank') => void;
    onToggleStory: () => void;
}

export const Header: React.FC<HeaderProps> = ({ mode, isStoryPlaying, onUpload, onSetMode, onToggleStory }) => (
    <header className="absolute top-4 left-1/2 -translate-x-1/2 z-50 bg-[#161b22]/90 backdrop-blur-md border border-[#30363d] rounded-2xl px-6 py-3 flex items-center gap-6 shadow-2xl">
        <h1 className="text-lg font-semibold flex items-center gap-2 text-white">
            <Share2 className="text-[#58a6ff]" size={20} />
            Repomap
        </h1>

        <div className="w-px h-6 bg-[#30363d]" />

        <label className="flex items-center gap-2 px-3 py-2 rounded-lg hover:bg-white/5 cursor-pointer transition-colors text-sm font-medium">
            <Upload size={16} /> Load JSON
            <input type="file" className="hidden" accept=".json" onChange={onUpload} />
        </label>

        <div className="w-px h-6 bg-[#30363d]" />

        <div className="flex gap-1">
            <button
                onClick={() => onSetMode('cluster')}
                className={`flex items-center gap-2 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${mode === 'cluster' ? 'bg-[#58a6ff]/20 text-[#58a6ff]' : 'hover:bg-white/5'}`}
            >
                <Layout size={16} /> Cluster
            </button>
            <button
                onClick={() => onSetMode('flow')}
                className={`flex items-center gap-2 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${mode === 'flow' ? 'bg-[#58a6ff]/20 text-[#58a6ff]' : 'hover:bg-white/5'}`}
            >
                <Share2 size={16} /> Flow
            </button>
        </div>

        <div className="w-px h-6 bg-[#30363d]" />

        <button
            onClick={onToggleStory}
            className={`flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-colors ${isStoryPlaying ? 'bg-red-500/20 text-red-500 hover:bg-red-500/30' : 'text-[#a371f7] hover:bg-[#a371f7]/10'
                }`}
        >
            {isStoryPlaying ? <Pause size={16} /> : <Play size={16} />}
            {isStoryPlaying ? 'Stop' : 'Play Story'}
        </button>
    </header>
);
