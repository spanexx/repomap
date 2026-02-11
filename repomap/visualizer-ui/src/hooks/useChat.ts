/**
 * Code Map: useChat – Custom hook encapsulating chat session management,
 * streaming API calls, history loading, and file-mention highlighting.
 * CID: 1.5.2-useChat
 */

import { useState, useEffect, useRef } from 'react';
import { type FileNode } from '../types';

export interface Message {
    role: 'user' | 'assistant';
    content: string;
}

interface UseChatOptions {
    files?: FileNode[];
    onHighlight?: (paths: string[]) => void;
    uiContext?: {
        selectedNode: string;
        viewMode: string;
    };
}

export function useChat({ files = [], onHighlight, uiContext }: UseChatOptions) {
    const [input, setInput] = useState('');
    const [messages, setMessages] = useState<Message[]>([]);
    const [isTyping, setIsTyping] = useState(false);
    const scrollRef = useRef<HTMLDivElement>(null);

    const [sessionId, setSessionId] = useState(() => {
        const stored = localStorage.getItem('repomap_session_id');
        if (stored) return stored;
        const newId = crypto.randomUUID();
        localStorage.setItem('repomap_session_id', newId);
        return newId;
    });

    // Auto-scroll on new messages
    useEffect(() => {
        if (scrollRef.current) {
            scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
        }
    }, [messages]);

    // Context-based file highlighting
    useEffect(() => {
        if (!onHighlight || messages.length === 0) return;

        const lastMsg = messages[messages.length - 1];
        if (!lastMsg || !lastMsg.content) return;

        if (lastMsg.content === '_RESET_') {
            onHighlight([]);
            return;
        }

        const content = lastMsg.content.toLowerCase();
        const matches: string[] = [];

        files.forEach(f => {
            const basename = f.path.split('/').pop()?.toLowerCase();
            if (content.includes(f.path.toLowerCase()) || (basename && content.includes(basename))) {
                matches.push(f.path);
            }
        });

        if (matches.length > 0) {
            onHighlight(matches);
        } else if (lastMsg.role === 'user') {
            onHighlight([]);
        }
    }, [messages, files, onHighlight]);

    // Load history on mount
    useEffect(() => {
        if (!sessionId) return;

        const fetchHistory = async () => {
            try {
                const res = await fetch(`/api/chat?sessionId=${sessionId}`);
                if (res.ok) {
                    const history = await res.json();
                    if (Array.isArray(history)) {
                        setMessages(history);
                    }
                }
            } catch (err) {
                console.error("Failed to load history", err);
            }
        };
        fetchHistory();
    }, [sessionId]);

    /** Send a message and stream the assistant response. */
    const handleSend = async () => {
        if (!input.trim() || isTyping) return;

        const userMsg: Message = { role: 'user', content: input };
        setMessages(prev => [...prev, userMsg]);
        setInput('');
        setIsTyping(true);

        try {
            const response = await fetch('/api/chat', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    message: input,
                    sessionId,
                    context: uiContext
                }),
            });

            if (!response.body) throw new Error('No body');

            const reader = response.body.getReader();
            const decoder = new TextDecoder();
            let assistantMsg = '';

            setMessages(prev => [...prev, { role: 'assistant', content: '' }]);

            while (true) {
                const { value, done } = await reader.read();
                if (done) break;

                const chunk = decoder.decode(value);
                const lines = chunk.split('\n');

                for (const line of lines) {
                    if (line.startsWith('data: ')) {
                        try {
                            const jsonStr = line.replace('data: ', '');
                            const data = JSON.parse(jsonStr);
                            if (data.token) {
                                assistantMsg += data.token;
                                setMessages(prev => {
                                    const last = prev[prev.length - 1];
                                    if (last.role === 'assistant') {
                                        return [...prev.slice(0, -1), { ...last, content: assistantMsg }];
                                    }
                                    return prev;
                                });
                            }
                        } catch (_e) {
                            // Ignore parse errors (empty lines, keep-alives)
                        }
                    }
                }
            }
        } catch (err) {
            console.error('Chat error:', err);
            setMessages(prev => [...prev, { role: 'assistant', content: 'Sorry, I encountered an error.' }]);
        } finally {
            setIsTyping(false);
        }
    };

    /** Reset highlight via the onHighlight callback. */
    const resetHighlight = () => onHighlight?.([]);

    /** Start a new chat session. */
    const newSession = () => {
        const newId = crypto.randomUUID();
        localStorage.setItem('repomap_session_id', newId);
        setSessionId(newId);
        setMessages([]);
        resetHighlight();
    };

    return {
        input,
        setInput,
        messages,
        isTyping,
        scrollRef,
        handleSend,
        resetHighlight,
        newSession,
    };
}
