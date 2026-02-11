import { motion } from 'framer-motion';
import { Play, Pause, ChevronRight, X } from 'lucide-react';
import { STORY_STEPS } from '../types';

interface StoryOverlayProps {
    stepIdx: number;
    progress: number;
    isPaused: boolean;
    onTogglePause: () => void;
    onNext: () => void;
    onExit: () => void;
}

export const StoryOverlay: React.FC<StoryOverlayProps> = ({
    stepIdx,
    progress,
    isPaused,
    onTogglePause,
    onNext,
    onExit
}) => {
    return (
        <motion.div
            initial={{ opacity: 0, y: 100, x: '-50%' }}
            animate={{ opacity: 1, y: 0, x: '-50%' }}
            exit={{ opacity: 0, y: 100, x: '-50%' }}
            className="absolute bottom-12 left-1/2 z-40 w-[550px] bg-[#161b22]/95 backdrop-blur-xl border border-[#a371f7]/40 rounded-2xl p-6 shadow-2xl overflow-hidden"
        >
            <div className="flex justify-between items-start mb-4">
                <div className="text-left">
                    <h2 className="text-xl font-bold text-[#a371f7] mb-1">{STORY_STEPS[stepIdx].title}</h2>
                    <p className="text-sm text-[#8b949e]">Step {stepIdx + 1} of {STORY_STEPS.length}</p>
                </div>
                <button
                    onClick={onExit}
                    className="p-2 hover:bg-white/10 rounded-full text-[#8b949e] hover:text-white transition-colors"
                >
                    <X size={20} />
                </button>
            </div>

            <p className="text-sm text-[#c9d1d9] leading-relaxed mb-6 text-left">{STORY_STEPS[stepIdx].text}</p>

            <div className="flex items-center gap-4 mb-4">
                <div className="flex-1 h-1.5 bg-white/10 rounded-full overflow-hidden">
                    <motion.div
                        className="h-full bg-[#a371f7] shadow-[0_0_10px_#a371f7]"
                        initial={{ width: 0 }}
                        animate={{ width: `${progress}%` }}
                        transition={{ duration: 0.1 }}
                    />
                </div>
                <div className="flex items-center gap-2">
                    <button
                        onClick={onTogglePause}
                        className="w-10 h-10 flex items-center justify-center bg-white/5 hover:bg-white/10 rounded-full text-[#c9d1d9] border border-white/10 transition-colors"
                    >
                        {isPaused ? <Play size={18} fill="currentColor" /> : <Pause size={18} fill="currentColor" />}
                    </button>
                    <button
                        onClick={onNext}
                        disabled={stepIdx >= STORY_STEPS.length - 1}
                        className="w-10 h-10 flex items-center justify-center bg-[#a371f7]/20 hover:bg-[#a371f7]/30 rounded-full text-[#a371f7] border border-[#a371f7]/30 transition-colors disabled:opacity-30"
                    >
                        <ChevronRight size={20} />
                    </button>
                </div>
            </div>
        </motion.div>
    );
};
