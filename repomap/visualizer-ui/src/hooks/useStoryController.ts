/**
 * Code Map: useStoryController – Custom hook for managing cinematic story playback state.
 * CID: 1.5.3-useStoryController
 */

import { useState, useCallback } from 'react';

export function useStoryController() {
    const [isStoryPlaying, setIsStoryPlaying] = useState(false);
    const [storyStep, setStoryStep] = useState(0);
    const [storyProgress, setStoryProgress] = useState(0);
    const [isStoryPaused, setIsStoryPaused] = useState(false);

    const toggleStory = useCallback(() => setIsStoryPlaying(prev => !prev), []);

    const onStoryUpdate = useCallback((step: number, progress: number) => {
        setStoryStep(step);
        setStoryProgress(progress);
    }, []);

    const onStoryComplete = useCallback(() => {
        setIsStoryPlaying(false);
    }, []);

    const togglePause = useCallback(() => setIsStoryPaused(prev => !prev), []);

    const nextStep = useCallback(() => {
        // Placeholder for future skip logic if needed
    }, []);

    return {
        isStoryPlaying,
        storyStep,
        storyProgress,
        isStoryPaused,
        toggleStory,
        onStoryUpdate,
        onStoryComplete,
        togglePause,
        nextStep,
        stopStory: () => setIsStoryPlaying(false)
    };
}
