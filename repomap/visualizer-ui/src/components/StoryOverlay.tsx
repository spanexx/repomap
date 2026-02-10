import { motion } from 'framer-motion';
import { STORY_STEPS } from '../types';

interface StoryOverlayProps {
    stepIdx: number;
    progress: number;
}

export const StoryOverlay: React.FC<StoryOverlayProps> = ({ stepIdx, progress }) => {
    return (
        <motion.div
            initial={{ opacity: 0, y: 100, x: '-50%' }}
            animate={{ opacity: 1, y: 0, x: '-50%' }}
            exit={{ opacity: 0, y: 100, x: '-50%' }}
            className="absolute bottom-12 left-1/2 z-40 w-[500px] bg-[#161b22]/90 backdrop-blur-md border border-[#a371f7]/50 rounded-xl p-6 shadow-2xl text-center"
        >
            <h2 className="text-xl font-bold text-[#a371f7] mb-2">{STORY_STEPS[stepIdx].title}</h2>
            <p className="text-sm text-[#c9d1d9] leading-relaxed mb-4">{STORY_STEPS[stepIdx].text}</p>
            <div className="h-1 bg-white/10 rounded-full overflow-hidden">
                <motion.div
                    className="h-full bg-[#a371f7]"
                    initial={{ width: 0 }}
                    animate={{ width: `${progress}%` }}
                    transition={{ duration: 0.1 }}
                />
            </div>
        </motion.div>
    );
};
