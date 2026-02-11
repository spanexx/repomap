/**
 * Code Map: App – Main application component, now refactored to compose custom hooks.
 * CID: 1.5.3-App
 */

import { useState } from 'react';
import { Layers } from 'lucide-react';
import { AnimatePresence } from 'framer-motion';

import { Header } from './components/Header';
import { Sidebar } from './components/Sidebar';
import { StoryOverlay } from './components/StoryOverlay';
import { Graph } from './components/Graph';
import { Chat } from './components/Chat';
import { TableView } from './components/TableView';

import { useRepoData } from './hooks/useRepoData';
import { useStoryController } from './hooks/useStoryController';
import { useSelection } from './hooks/useSelection';

function App() {
  const [mode, setMode] = useState<'cluster' | 'flow' | 'simple' | 'table'>('cluster');
  const [colorMode, setColorMode] = useState<'importance' | 'intent'>('importance');

  const { data, files, handleUpload } = useRepoData();

  const {
    isStoryPlaying,
    storyStep,
    storyProgress,
    isStoryPaused,
    toggleStory,
    onStoryUpdate,
    onStoryComplete,
    togglePause,
    nextStep,
    stopStory
  } = useStoryController();

  const {
    selectedFile,
    setSelectedFile,
    highlightedNodes,
    setHighlightedNodes,
    toggleHighlight
  } = useSelection();

  return (
    <div className="w-full h-full relative bg-[#0d1117] text-[#c9d1d9] overflow-hidden font-sans">

      <Header
        mode={mode}
        isStoryPlaying={isStoryPlaying}
        onUpload={handleUpload}
        onSetMode={(m) => setMode(m as any)}
        onToggleStory={toggleStory}
        colorMode={colorMode}
        onSetColorMode={setColorMode}
      />

      {mode === 'table' ? (
        <TableView
          files={files}
          onFileSelect={setSelectedFile}
        />
      ) : (
        <Graph
          files={files}
          mode={mode as any}
          isStoryPlaying={isStoryPlaying}
          isPaused={isStoryPaused}
          highlightedNodes={highlightedNodes}
          onNodeSelect={setSelectedFile}
          onStoryUpdate={onStoryUpdate}
          onStoryComplete={onStoryComplete}
          colorMode={colorMode}
        />
      )}

      {/* Loader / Empty State */}
      {!data && !isStoryPlaying && (
        <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
          <div className="text-[#8b949e] flex flex-col items-center gap-3">
            <Layers size={48} strokeWidth={1} />
            <p>Upload a <code>repomap.json</code> to verify</p>
          </div>
        </div>
      )}

      {/* Story Overlay */}
      <AnimatePresence>
        {isStoryPlaying && (
          <StoryOverlay
            stepIdx={storyStep}
            progress={storyProgress}
            isPaused={isStoryPaused}
            onTogglePause={togglePause}
            onNext={nextStep}
            onExit={stopStory}
          />
        )}
      </AnimatePresence>

      {/* Sidebar */}
      <AnimatePresence>
        {selectedFile && (
          <Sidebar
            file={selectedFile}
            onClose={() => setSelectedFile(null)}
            isHighlighted={highlightedNodes.includes(selectedFile.path)}
            onToggleHighlight={() => toggleHighlight(selectedFile.path)}
          />
        )}
      </AnimatePresence>

      <Chat
        files={files}
        onHighlight={setHighlightedNodes}
        selectedNode={selectedFile?.path || ''}
        viewMode={mode}
      />
    </div>
  );
}

export default App;
