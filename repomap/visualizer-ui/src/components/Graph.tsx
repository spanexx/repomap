import { useEffect, useRef } from 'react';
import { Network, type Options } from 'vis-network';
import { DataSet } from 'vis-data';
import { type FileNode, COLORS, STORY_STEPS } from '../types';

interface GraphProps {
    files: FileNode[];
    mode: 'cluster' | 'flow' | 'rank';
    isStoryPlaying: boolean;
    onNodeSelect: (file: FileNode | null) => void;
    onStoryUpdate: (step: number, progress: number) => void;
    onStoryComplete: () => void;
}

interface GraphState {
    nodes: DataSet<any>;
    edges: DataSet<any>;
}

export const Graph: React.FC<GraphProps> = ({
    files,
    mode,
    isStoryPlaying,
    onNodeSelect,
    onStoryUpdate,
    onStoryComplete
}) => {
    const containerRef = useRef<HTMLDivElement>(null);
    const networkRef = useRef<Network | null>(null);
    const graphState = useRef<GraphState>({
        nodes: new DataSet([]),
        edges: new DataSet([])
    });

    // --- Helpers ---
    const clearGraph = () => {
        graphState.current.nodes.clear();
        graphState.current.edges.clear();
    };

    const addFileNode = (f: FileNode, withEdges = true) => {
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
                color: { background: COLORS.bg, border: COLORS.edge },
                shape: 'hexagon',
                font: { color: COLORS.pkg, size: 14 }
            });
        }

        let color = COLORS.low;
        if (f.importance === 'medium') color = COLORS.medium;
        if (f.importance === 'high') color = COLORS.high;

        graphState.current.nodes.add({
            id: f.path,
            label: f.path.split('/').pop(),
            group: 'file',
            value: f.rank ? (f.rank * 10) + 5 : 5,
            color: { background: color, border: '#fff' }
        });

        if (withEdges) {
            // Structural Edge
            graphState.current.edges.add({
                from: f.path,
                to: pkgId,
                color: { color: color, opacity: 0.15 },
                width: 1
            });

            // Flow Edges
            f.imports?.forEach(imp => {
                if (imp.startsWith('.')) return;
                // Heuristic mapping
                const pkgNodes = graphState.current.nodes.get({ filter: (n: any) => n.group === 'package' });
                const target = pkgNodes.find((n: any) => imp.endsWith(n.label));

                if (target && target.id !== pkgId) {
                    graphState.current.edges.add({
                        from: f.path,
                        to: target.id,
                        arrows: 'to',
                        color: { color: COLORS.accent, opacity: 0.3 },
                        dashes: true
                    });
                }
            });
        }
    };

    // --- Init ---
    useEffect(() => {
        if (!containerRef.current) return;

        const options: Options = {
            nodes: {
                shape: 'dot',
                font: { color: '#f0f6fc', face: 'Inter, system-ui' }
            },
            physics: {
                forceAtlas2Based: {
                    gravitationalConstant: -26,
                    centralGravity: 0.005,
                    springLength: 230,
                    springConstant: 0.18
                },
                maxVelocity: 146,
                solver: 'forceAtlas2Based',
                timestep: 0.35,
                stabilization: { enabled: true, iterations: 200 }
            },
            interaction: { hover: true, tooltipDelay: 200, hideEdgesOnDrag: true },
            layout: {
                hierarchical: {
                    enabled: mode === 'flow' || mode === 'rank',
                    direction: mode === 'rank' ? 'UD' : 'LR', // UD for Rank, LR for Flow
                    sortMethod: 'directed',
                    levelSeparation: 150,
                    nodeSpacing: 100
                }
            }
        };

        const net = new Network(containerRef.current, {
            nodes: graphState.current.nodes,
            edges: graphState.current.edges
        }, options);

        net.on('click', (params) => {
            if (params.nodes.length) {
                const nodeId = params.nodes[0];
                const file = files.find(f => f.path === nodeId);
                onNodeSelect(file || null);
            } else {
                onNodeSelect(null);
            }
        });

        networkRef.current = net;
        return () => net.destroy();
    }, [files, onNodeSelect, mode]); // Re-init on mode change to apply layout

    // --- Render Logic ---
    useEffect(() => {
        if (isStoryPlaying) return;

        clearGraph();

        if (mode === 'cluster' || mode === 'flow' || mode === 'rank') {
            files.forEach(f => addFileNode(f, true));
        }

    }, [files, mode, isStoryPlaying]);

    // --- Story Logic ---
    useEffect(() => {
        if (!isStoryPlaying) return;

        let isActive = true;

        const playStep = async (stepIdx: number) => {
            if (!isActive) return;
            if (stepIdx >= STORY_STEPS.length) {
                onStoryComplete();
                return;
            }

            onStoryUpdate(stepIdx, 0);

            const step = STORY_STEPS[stepIdx];
            const batch = files.filter(step.filter);

            for (let i = 0; i < batch.length; i++) {
                if (!isActive) break;

                addFileNode(batch[i], true);
                onStoryUpdate(stepIdx, ((i + 1) / batch.length) * 100);

                await new Promise(r => setTimeout(r, 50));
            }

            if (isActive) {
                await new Promise(r => setTimeout(r, 2000));
                playStep(stepIdx + 1);
            }
        };

        clearGraph();
        playStep(0);

        return () => { isActive = false; };
    }, [isStoryPlaying, files]); // Removed onStoryUpdate/Complete from deps loop

    return <div ref={containerRef} className="w-full h-full cursor-grab active:cursor-grabbing" />;
};
