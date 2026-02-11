import { useEffect, useRef } from 'react';
import { Network } from 'vis-network';
import { DataSet } from 'vis-data';
import { type FileNode, COLORS, STORY_STEPS } from '../types';

interface GraphProps {
    files: FileNode[];
    mode: 'cluster' | 'flow' | 'simple' | 'table';
    isStoryPlaying: boolean;
    onNodeSelect: (file: FileNode | null) => void;
    onStoryUpdate: (step: number, progress: number) => void;
    onStoryComplete: () => void;
    highlightedNodes?: string[];
    isPaused?: boolean;
    colorMode?: 'importance' | 'intent';
}

interface GraphState {
    nodes: DataSet<any>;
    edges: DataSet<any>;
}

// Intent Color Palette
const INTENT_COLORS: Record<string, string> = {
    'entry point': '#f778ba',   // Pink
    'core domain': '#a371f7',   // Purple
    'business logic': '#58a6ff',// Blue
    'infrastructure': '#8b949e',// Grey
    'data model': '#238636',    // Green
    'interface/api': '#f2cc60', // Yellow
    'utility': '#79c0ff',       // Light Blue
    'configuration': '#d2a8ff', // Light Purple
    'test': '#ff7b72',          // Red
    'unknown': '#30363d'        // Dark Grey
};

const getIntentColor = (intent: string) => {
    if (!intent) return COLORS.low;
    const key = intent.toLowerCase();
    // Prefix matching
    if (key.includes('test')) return INTENT_COLORS['test'];
    if (key.includes('config')) return INTENT_COLORS['configuration'];
    if (key.includes('util')) return INTENT_COLORS['utility'];
    if (key.includes('interface') || key.includes('api')) return INTENT_COLORS['interface/api'];
    if (key.includes('data') || key.includes('model') || key.includes('entity')) return INTENT_COLORS['data model'];
    if (key.includes('infra') || key.includes('adapter') || key.includes('db')) return INTENT_COLORS['infrastructure'];
    if (key.includes('usecase') || key.includes('service') || key.includes('logic')) return INTENT_COLORS['business logic'];
    if (key.includes('domain') || key.includes('core')) return INTENT_COLORS['core domain'];
    if (key.includes('main') || key.includes('cmd') || key.includes('entry')) return INTENT_COLORS['entry point'];

    return COLORS.low;
};

export const Graph: React.FC<GraphProps> = ({
    files,
    mode,
    isStoryPlaying,
    onNodeSelect,
    onStoryUpdate,
    onStoryComplete,
    highlightedNodes = [],
    isPaused = false,
    colorMode = 'importance'
}) => {
    const containerRef = useRef<HTMLDivElement>(null);
    const networkRef = useRef<Network | null>(null);
    const graphState = useRef<GraphState>({
        nodes: new DataSet([]),
        edges: new DataSet([])
    });
    const filesRef = useRef(files);
    const isStoryPlayingRef = useRef(isStoryPlaying);

    useEffect(() => {
        filesRef.current = files;
    }, [files]);

    useEffect(() => {
        isStoryPlayingRef.current = isStoryPlaying;
    }, [isStoryPlaying]);

    // --- Helpers ---
    const clearGraph = () => {
        graphState.current.nodes.clear();
        graphState.current.edges.clear();
    };

    const addFileNode = (f: FileNode, withEdges = true, isDimmed = false) => {
        if (graphState.current.nodes.get(f.path)) return;

        const dir = f.path.substring(0, f.path.lastIndexOf('/')) || ".";
        const pkgId = `pkg:${dir}`;

        // Add Package Node if needed
        if (!graphState.current.nodes.get(pkgId)) {
            graphState.current.nodes.add({
                id: pkgId,
                label: dir,
                group: 'package',
                value: 20,
                color: {
                    background: isDimmed ? 'rgba(30,30,30, 0.2)' : COLORS.bg,
                    border: isDimmed ? 'rgba(50,50,50, 0.2)' : COLORS.edge
                },
                shape: 'hexagon',
                font: {
                    color: isDimmed ? 'rgba(100,100,100, 0.2)' : COLORS.pkg,
                    size: 14
                }
            });
        }

        let color = COLORS.low;
        if (colorMode === 'intent') {
            color = getIntentColor(f.intent || '');
        } else {
            if (f.importance === 'medium') color = COLORS.medium;
            if (f.importance === 'high') color = COLORS.high;
        }

        graphState.current.nodes.add({
            id: f.path,
            label: f.path.split('/').pop(),
            group: 'file',
            value: f.rank ? (f.rank * 10) + 5 : 5,
            color: {
                background: isDimmed ? 'rgba(50,50,50, 0.1)' : color,
                border: isDimmed ? 'rgba(100,100,100, 0.1)' : '#fff'
            },
            font: {
                color: isDimmed ? 'rgba(150,150,150, 0.3)' : '#f0f6fc'
            }
        });

        if (withEdges) {
            // Structural Edge
            graphState.current.edges.add({
                from: f.path,
                to: pkgId,
                color: { color: color, opacity: isDimmed ? 0.05 : 0.15 },
                width: 1
            });

            // Flow Edges
            f.imports?.forEach(imp => {
                if (imp.startsWith('.')) return;
                const pkgNodes = graphState.current.nodes.get({ filter: (n: any) => n.group === 'package' });
                const target = pkgNodes.find((n: any) => imp.endsWith(n.label));

                if (target && target.id !== pkgId) {
                    graphState.current.edges.add({
                        from: f.path,
                        to: target.id,
                        arrows: 'to',
                        color: { color: COLORS.accent, opacity: isDimmed ? 0.05 : 0.3 },
                        dashes: true
                    });
                }
            });
        }
    };

    // --- Init (Once) ---
    useEffect(() => {
        if (!containerRef.current) return;

        const net = new Network(containerRef.current, {
            nodes: graphState.current.nodes,
            edges: graphState.current.edges
        }, {
            nodes: {
                shape: 'dot',
                font: { color: '#f0f6fc', face: 'Inter, system-ui' },
                borderWidth: 2,
                scaling: {
                    min: 10,
                    max: 30,
                    label: { enabled: true, min: 12, max: 20 }
                }
            },
            physics: {
                enabled: true,
                solver: 'barnesHut',
                barnesHut: {
                    gravitationalConstant: -2000,
                    centralGravity: 0.3,
                    springLength: 230,
                    springConstant: 0.04,
                    damping: 0.4,
                    avoidOverlap: 0.5
                },
                stabilization: {
                    enabled: true,
                    iterations: 200,
                    updateInterval: 25
                },
                maxVelocity: 30,
                minVelocity: 0.75,
                timestep: 0.35
            },
            interaction: {
                hover: true,
                tooltipDelay: 200,
                hideEdgesOnDrag: true,
                zoomView: true,
                dragView: true
            }
        });

        net.on('click', (params) => {
            if (params.nodes.length) {
                const nodeId = params.nodes[0];
                const file = filesRef.current.find(f => f.path === nodeId);
                onNodeSelect(file || null);
            } else {
                onNodeSelect(null);
            }
        });

        // Optimization: Stop physics after stabilization to save resources, 
        // BUT NOT during story mode reveal!
        net.on('stabilized', () => {
            if (!isStoryPlayingRef.current) {
                net.setOptions({ physics: { enabled: false } });
            }
        });

        networkRef.current = net;
        return () => {
            if (networkRef.current) {
                networkRef.current.destroy();
                networkRef.current = null;
            }
        };
    }, []); // Only once

    // --- Render Logic (Structure & Mode) ---
    useEffect(() => {
        if (isStoryPlaying) return;

        clearGraph();

        // Initial add of all nodes
        files.forEach(f => {
            addFileNode(f, true, false);
        });

        if (networkRef.current) {
            const isSimple = mode === 'simple';
            const usePhysics = mode === 'cluster';

            networkRef.current.setOptions({
                layout: {
                    hierarchical: {
                        enabled: mode === 'flow',
                        direction: 'LR',
                        sortMethod: 'directed',
                        levelSeparation: 250,
                        nodeSpacing: 200,
                        parentCentralization: true,
                        blockShifting: true,
                        edgeMinimization: true
                    }
                },
                physics: {
                    enabled: usePhysics,
                    solver: 'barnesHut',
                    barnesHut: {
                        gravitationalConstant: -3000,
                        centralGravity: 0.25,
                        springLength: 180,
                        springConstant: 0.04,
                        damping: 0.4,
                        avoidOverlap: 0.5
                    },
                    stabilization: {
                        enabled: true,
                        iterations: 150,
                        updateInterval: 25
                    },
                    maxVelocity: 30,
                    minVelocity: 0.75,
                    timestep: 0.35
                },
                edges: {
                    hidden: isSimple
                }
            });

            // --- Simple Mode: Static grid layout with box nodes ---
            if (isSimple) {
                const groups: Record<string, FileNode[]> = {};
                files.forEach(f => {
                    const dir = f.path.substring(0, f.path.lastIndexOf('/')) || "root";
                    if (!groups[dir]) groups[dir] = [];
                    groups[dir].push(f);
                });

                const groupKeys = Object.keys(groups).sort();
                const COLS = 4;
                const CELL_W = 220;
                const CELL_H = 60;
                const GROUP_GAP = 80;
                const HEADER_H = 50;

                let cursorY = 0;
                const nodeUpdates: any[] = [];
                const pkgUpdates: any[] = [];

                groupKeys.forEach((group) => {
                    const pkgId = `pkg:${group}`;
                    const children = groups[group];
                    const rows = Math.ceil(children.length / COLS);
                    const groupWidth = Math.min(children.length, COLS) * CELL_W;

                    // Directory header
                    pkgUpdates.push({
                        id: pkgId,
                        x: groupWidth / 2 - CELL_W / 2,
                        y: cursorY,
                        shape: 'box',
                        label: `📁 ${group === 'root' ? '/' : group}`,
                        color: {
                            background: 'rgba(88, 166, 255, 0.08)',
                            border: 'rgba(88, 166, 255, 0.25)'
                        },
                        font: { size: 13, color: '#8b949e', face: 'Inter, system-ui' },
                        borderWidth: 1,
                        physics: false,
                        fixed: true,
                        widthConstraint: { minimum: groupWidth + 40 },
                        margin: { top: 8, bottom: 8, left: 14, right: 14 }
                    });

                    cursorY += HEADER_H;

                    // File boxes in grid rows
                    children.forEach((f, cIdx) => {
                        const col = cIdx % COLS;
                        const row = Math.floor(cIdx / COLS);
                        let color = COLORS.low;
                        if (colorMode === 'intent') {
                            color = getIntentColor(f.intent || '');
                        } else {
                            if (f.importance === 'medium') color = COLORS.medium;
                            if (f.importance === 'high') color = COLORS.high;
                        }

                        nodeUpdates.push({
                            id: f.path,
                            x: col * CELL_W,
                            y: cursorY + row * CELL_H,
                            shape: 'box',
                            label: f.path.split('/').pop(),
                            color: {
                                background: color,
                                border: 'rgba(255,255,255,0.12)'
                            },
                            font: { color: '#f0f6fc', size: 12, face: 'Inter, system-ui' },
                            borderWidth: 1,
                            physics: false,
                            fixed: true,
                            widthConstraint: { minimum: 150, maximum: 200 },
                            margin: { top: 6, bottom: 6, left: 10, right: 10 }
                        });
                    });

                    cursorY += rows * CELL_H + GROUP_GAP;
                });

                graphState.current.nodes.update([...pkgUpdates, ...nodeUpdates]);
                networkRef.current.fit({ animation: { duration: 500, easingFunction: 'easeInOutQuad' } });
            }

            if (usePhysics) {
                networkRef.current.stabilize();
            }
        }

    }, [files, mode, isStoryPlaying, colorMode]); // Added colorMode dependence!

    // --- Highlighting Logic (Performance Optimized) ---
    useEffect(() => {
        if (isStoryPlaying || !graphState.current.nodes.length) return;

        const allNodeIds = graphState.current.nodes.getIds();
        const updates: any[] = [];

        const hasHighlight = highlightedNodes && highlightedNodes.length > 0;

        allNodeIds.forEach(id => {
            const nodeId = id as string;
            const isPackage = nodeId.startsWith('pkg:');

            if (isPackage) {
                updates.push({
                    id: nodeId,
                    color: {
                        background: hasHighlight ? 'rgba(30,30,30, 0.1)' : COLORS.bg,
                        border: hasHighlight ? 'rgba(50,50,50, 0.1)' : COLORS.edge
                    },
                    font: { color: hasHighlight ? 'rgba(100,100,100, 0.1)' : COLORS.pkg }
                });
                return;
            }

            const isHighlighted = highlightedNodes.includes(nodeId);
            const file = files.find(f => f.path === nodeId);
            if (!file) return;

            let baseColor = COLORS.low;
            if (colorMode === 'intent') {
                baseColor = getIntentColor(file.intent || '');
            } else {
                if (file.importance === 'medium') baseColor = COLORS.medium;
                if (file.importance === 'high') baseColor = COLORS.high;
            }

            if (hasHighlight) {
                if (isHighlighted) {
                    updates.push({
                        id: nodeId,
                        color: { background: '#ff4444', border: '#fff' },
                        font: { color: '#fff', size: 16 },
                        value: (file.rank ? (file.rank * 10) : 5) + 10 // Slightly larger
                    });
                } else {
                    updates.push({
                        id: nodeId,
                        color: { background: 'rgba(50,50,50, 0.05)', border: 'rgba(100,100,100, 0.05)' },
                        font: { color: 'rgba(150,150,150, 0.1)' },
                        value: file.rank ? (file.rank * 10) + 5 : 5
                    });
                }
            } else {
                // Reset to default
                updates.push({
                    id: nodeId,
                    color: { background: baseColor, border: '#fff' },
                    font: { color: '#f0f6fc', size: 14 },
                    value: file.rank ? (file.rank * 10) + 5 : 5
                });
            }
        });

        graphState.current.nodes.update(updates);

        // Also update edges
        const edgeUpdates: any[] = [];
        graphState.current.edges.getIds().forEach(id => {
            const edge = graphState.current.edges.get(id);
            const isRelatedToHighlight = hasHighlight && highlightedNodes.some(path => edge.from === path || edge.to === path);

            edgeUpdates.push({
                id,
                color: {
                    color: isRelatedToHighlight ? COLORS.accent : COLORS.edge,
                    opacity: hasHighlight ? (isRelatedToHighlight ? 0.6 : 0.02) : 0.15
                },
                width: isRelatedToHighlight ? 2 : 1
            });
        });
        graphState.current.edges.update(edgeUpdates);

    }, [highlightedNodes, isStoryPlaying, files, colorMode]); // Added colorMode dependence!

    // --- Story Logic ---
    useEffect(() => {
        if (!isStoryPlaying) return;

        let isActive = true;
        const state = {
            currentStep: 0,
            currentNodeIdx: 0,
            isWaitingBetweenSteps: false
        };

        const playCycle = async () => {
            if (!isActive) return;

            if (state.currentStep >= STORY_STEPS.length) {
                onStoryComplete();
                return;
            }

            if (isPaused) {
                setTimeout(playCycle, 100);
                return;
            }

            const step = STORY_STEPS[state.currentStep];
            const batch = files.filter(step.filter);

            if (state.currentNodeIdx === 0 && !state.isWaitingBetweenSteps) {
                // Start of a new step: Fit view to context or focus on previous
                if (networkRef.current) {
                    networkRef.current.fit({ animation: { duration: 1000, easingFunction: 'easeInOutQuad' } });
                    await new Promise(r => setTimeout(r, 1100));
                }
                onStoryUpdate(state.currentStep, 0);
            }

            if (state.currentNodeIdx < batch.length) {
                const f = batch[state.currentNodeIdx];

                // Ensure physics is enabled for the reveal
                if (networkRef.current) {
                    networkRef.current.setOptions({ physics: { enabled: true } });
                }

                addFileNode(f, true);

                // Focus camera on new node
                if (networkRef.current) {
                    networkRef.current.focus(f.path, {
                        scale: 1.2,
                        animation: { duration: 400, easingFunction: 'easeOutCubic' }
                    });
                }

                onStoryUpdate(state.currentStep, ((state.currentNodeIdx + 1) / batch.length) * 100);
                state.currentNodeIdx++;

                // Pulsing reveal timing
                await new Promise(r => setTimeout(r, 600));
                setTimeout(playCycle, 50);
            } else {
                // Step finished
                state.isWaitingBetweenSteps = true;
                await new Promise(r => setTimeout(r, 2500));
                if (!isActive) return;

                state.currentStep++;
                state.currentNodeIdx = 0;
                state.isWaitingBetweenSteps = false;
                setTimeout(playCycle, 100);
            }
        };

        clearGraph();
        playCycle();

        return () => { isActive = false; };
    }, [isStoryPlaying, files, isPaused]);

    return <div ref={containerRef} className="w-full h-full cursor-grab active:cursor-grabbing" />;
};
