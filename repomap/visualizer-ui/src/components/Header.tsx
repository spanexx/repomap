import { useState } from 'react';
import { Share2, Upload, Layout, Play, Pause, Menu, X, Table, Palette } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';

interface HeaderProps {
    mode: 'cluster' | 'flow' | 'simple' | 'table';
    isStoryPlaying: boolean;
    onUpload: (e: React.ChangeEvent<HTMLInputElement>) => void;
    onSetMode: (mode: 'cluster' | 'flow' | 'simple' | 'table') => void;
    onToggleStory: () => void;
    colorMode?: 'importance' | 'intent';
    onSetColorMode?: (mode: 'importance' | 'intent') => void;
}

export const Header: React.FC<HeaderProps> = ({
    mode,
    isStoryPlaying,
    onUpload,
    onSetMode,
    onToggleStory,
    colorMode = 'importance',
    onSetColorMode
}) => {
    const [isMenuOpen, setIsMenuOpen] = useState(false);

    const menuItems = (
        <>
            <label className="flex items-center justify-center p-2 rounded-lg hover:bg-white/5 cursor-pointer transition-colors text-[#8b949e] hover:text-[#c9d1d9]" title="Load JSON">
                <Upload size={20} />
                <input type="file" className="hidden" accept=".json" onChange={onUpload} />
            </label>

            <div className="w-px h-6 bg-[#30363d] hidden md:block" />

            <div className="flex flex-col md:flex-row gap-1">
                <button
                    onClick={() => { onSetMode('table'); setIsMenuOpen(false); }}
                    className={`p-2 rounded-lg transition-colors ${mode === 'table' ? 'bg-[#58a6ff]/20 text-[#58a6ff]' : 'text-[#8b949e] hover:bg-white/5 hover:text-[#c9d1d9]'}`}
                    title="Table View"
                >
                    <Table size={20} />
                </button>
                <div className="w-px h-4 bg-[#30363d] hidden md:block self-center mx-1" />
                <button
                    onClick={() => { onSetMode('cluster'); setIsMenuOpen(false); }}
                    className={`p-2 rounded-lg transition-colors ${mode === 'cluster' ? 'bg-[#58a6ff]/20 text-[#58a6ff]' : 'text-[#8b949e] hover:bg-white/5 hover:text-[#c9d1d9]'}`}
                    title="Cluster View"
                >
                    <Layout size={20} />
                </button>
                <button
                    onClick={() => { onSetMode('simple'); setIsMenuOpen(false); }}
                    className={`p-2 rounded-lg transition-colors ${mode === 'simple' ? 'bg-[#58a6ff]/20 text-[#58a6ff]' : 'text-[#8b949e] hover:bg-white/5 hover:text-[#c9d1d9]'}`}
                    title="Simple View"
                >
                    <Layout size={20} className="rotate-45" />
                </button>
            </div>

            {mode !== 'table' && onSetColorMode && (
                <>
                    <div className="w-px h-6 bg-[#30363d] hidden md:block" />
                    <button
                        onClick={() => {
                            onSetColorMode(colorMode === 'importance' ? 'intent' : 'importance');
                            setIsMenuOpen(false);
                        }}
                        className={`p-2 rounded-lg transition-colors hover:bg-white/5 ${colorMode === 'intent' ? 'text-[#a371f7]' : 'text-[#8b949e]'}`}
                        title={`Color by: ${colorMode === 'intent' ? 'Intent' : 'Importance'}`}
                    >
                        <Palette size={20} />
                    </button>
                </>
            )}

            <div className="w-px h-6 bg-[#30363d] hidden md:block" />

            <button
                onClick={() => { onToggleStory(); setIsMenuOpen(false); }}
                className={`p-2 rounded-lg transition-colors ${isStoryPlaying ? 'bg-red-500/20 text-red-500 hover:bg-red-500/30' : 'text-[#a371f7] hover:bg-[#a371f7]/10'
                    }`}
                title={isStoryPlaying ? 'Stop Story' : 'Play Story'}
            >
                {isStoryPlaying ? <Pause size={20} /> : <Play size={20} />}
            </button>
        </>
    );

    return (
        <>
            <header className="absolute top-4 left-4 right-4 md:left-1/2 md:right-auto md:-translate-x-1/2 z-50 bg-[#161b22]/90 backdrop-blur-md border border-[#30363d] rounded-2xl px-4 py-3 flex items-center justify-between md:justify-start md:gap-6 shadow-2xl">
                <h1 className="text-lg font-semibold flex items-center gap-2 text-white shrink-0">
                    <Share2 className="text-[#58a6ff]" size={20} />
                    <span className="hidden sm:inline">Repomap</span>
                </h1>

                {/* Desktop Menu */}
                <div className="hidden md:flex items-center gap-6">
                    {menuItems}
                </div>

                {/* Mobile Menu Toggle */}
                <button
                    onClick={() => setIsMenuOpen(!isMenuOpen)}
                    className="md:hidden p-2 text-[#8b949e] hover:text-white transition-colors"
                >
                    {isMenuOpen ? <X size={24} /> : <Menu size={24} />}
                </button>
            </header>

            {/* Mobile Menu Dropdown */}
            <AnimatePresence>
                {isMenuOpen && (
                    <motion.div
                        initial={{ opacity: 0, y: -20 }}
                        animate={{ opacity: 1, y: 0 }}
                        exit={{ opacity: 0, y: -20 }}
                        className="fixed top-20 left-4 right-4 z-40 bg-[#161b22] border border-[#30363d] rounded-2xl shadow-2xl p-4 flex flex-col gap-4 md:hidden"
                    >
                        {menuItems}
                    </motion.div>
                )}
            </AnimatePresence>
        </>
    );
};
