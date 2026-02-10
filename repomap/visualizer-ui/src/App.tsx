
import { useState, useEffect } from 'react';
import { Layers } from 'lucide-react';
import { AnimatePresence } from 'framer-motion';

import { Header } from './components/Header';
import { Sidebar } from './components/Sidebar';
import { StoryOverlay } from './components/StoryOverlay';
import { Graph } from './components/Graph';
import { Chat } from './components/Chat';
import { type RepoMapData, type FileNode } from './types';

function App() {
  const [data, setData] = useState<RepoMapData | null>(null);
  const [mode, setMode] = useState<'cluster' | 'flow' | 'rank'>('cluster');
  const [selectedFile, setSelectedFile] = useState<FileNode | null>(null);
  const [isStoryPlaying, setIsStoryPlaying] = useState(false);
  const [storyStep, setStoryStep] = useState(0);
  const [storyProgress, setStoryProgress] = useState(0);

  // --- Effects ---
  useEffect(() => {
    // Fetch initial plan from backend
    const fetchPlan = async () => {
      try {
        const response = await fetch('/api/plan');
        if (response.ok) {
          const json = await response.json();
          setData(json);
        }
      } catch (err) {
        console.warn('Backend not available or plan missing');
      }
    };
    fetchPlan();
  }, []);

  // --- Handlers ---
  const handleUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    const reader = new FileReader();
    reader.onload = (ev) => {
      try {
        const json = JSON.parse(ev.target?.result as string);
        setData(json);
      } catch (err) {
        alert("Invalid JSON");
      }
    };
    reader.readAsText(file);
  };

  const getFiles = () => (data?.repomap?.files || data?.files || []);

  const handleNodeSelect = (file: FileNode | null) => {
    setSelectedFile(file);
  };

  const handleStoryUpdate = (step: number, progress: number) => {
    setStoryStep(step);
    setStoryProgress(progress);
  };

  const handleStoryComplete = () => {
    setIsStoryPlaying(false);
  };

  return (
    <div className="w-full h-full relative bg-[#0d1117] text-[#c9d1d9] overflow-hidden font-sans">

      <Header
        mode={mode}
        isStoryPlaying={isStoryPlaying}
        onUpload={handleUpload}
        onSetMode={setMode}
        onToggleStory={() => setIsStoryPlaying(!isStoryPlaying)}
      />

      <Graph
        files={getFiles()}
        mode={mode}
        isStoryPlaying={isStoryPlaying}
        onNodeSelect={handleNodeSelect}
        onStoryUpdate={handleStoryUpdate}
        onStoryComplete={handleStoryComplete}
      />

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
          <StoryOverlay stepIdx={storyStep} progress={storyProgress} />
        )}
      </AnimatePresence>

      {/* Sidebar */}
      <AnimatePresence>
        {selectedFile && (
          <Sidebar file={selectedFile} onClose={() => setSelectedFile(null)} />
        )}
      </AnimatePresence>

      <Chat />
    </div>
  );
}

export default App;
