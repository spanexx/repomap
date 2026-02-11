/**
 * Code Map: Chat – Composed chat panel using useChat hook and sub-components.
 * CID: 1.5.2-Chat
 */

import { useState } from 'react';
import { Bot, MessageSquare, X, RefreshCcw } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';
import { type FileNode } from '../types';
import { useChat } from '../hooks/useChat';
import { MessageList } from './chat/MessageList';
import { MessageInput } from './chat/MessageInput';

interface ChatProps {
    files?: FileNode[];
    onHighlight?: (paths: string[]) => void;
    selectedNode?: string;
    viewMode?: string;
    embedded?: boolean;
}

export const Chat: React.FC<ChatProps> = ({
    files = [],
    onHighlight,
    selectedNode = '',
    viewMode = '',
    embedded = false
}) => {
    const [isOpen, setIsOpen] = useState(false);

    const {
        input,
        setInput,
        messages,
        isTyping,
        scrollRef,
        handleSend,
        resetHighlight,
        newSession,
    } = useChat({
        files,
        onHighlight,
        uiContext: { selectedNode, viewMode }
    });

    const ChatContent = (
        <div className={`flex flex-col overflow-hidden bg-[#161b22] border-[#30363d] ${embedded ? 'h-full border-0 rounded-none' : 'w-full h-full sm:w-[400px] sm:h-[600px] border rounded-none sm:rounded-2xl shadow-2xl z-50'}`}>
            {/* Header */}
            <div className="p-4 border-b border-[#30363d] flex items-center justify-between bg-[#1f242c]">
                <div className="flex items-center gap-2">
                    <div className="w-8 h-8 rounded-full bg-[#238636]/20 flex items-center justify-center text-[#238636]">
                        <Bot size={18} />
                    </div>
                    <span className="font-semibold text-sm">Planning Agent</span>
                </div>
                <div className="flex items-center gap-2">
                    <button
                        onClick={resetHighlight}
                        className="text-[#8b949e] hover:text-[#58a6ff] transition-colors p-1"
                        title="Reset Highlight"
                    >
                        <RefreshCcw size={16} />
                    </button>
                    <button
                        onClick={newSession}
                        className="text-[#8b949e] hover:text-white transition-colors"
                        title="New Session"
                    >
                        <MessageSquare size={18} />
                    </button>
                    {!embedded && (
                        <button onClick={() => setIsOpen(false)} className="text-[#8b949e] hover:text-white transition-colors">
                            <X size={20} />
                        </button>
                    )}
                </div>
            </div>

            <MessageList
                messages={messages}
                isTyping={isTyping}
                scrollRef={scrollRef}
            />

            <MessageInput
                input={input}
                isTyping={isTyping}
                onInputChange={setInput}
                onSend={handleSend}
            />
        </div>
    );

    if (embedded) {
        return ChatContent;
    }

    return (
        <>
            {/* Floating Toggle Button */}
            {!isOpen && (
                <motion.button
                    initial={{ scale: 0 }}
                    animate={{ scale: 1 }}
                    onClick={() => setIsOpen(true)}
                    className="fixed bottom-6 right-6 w-14 h-14 bg-[#238636] hover:bg-[#2ea043] rounded-full flex items-center justify-center shadow-lg text-white z-50 transition-colors"
                >
                    <MessageSquare size={24} />
                </motion.button>
            )}

            {/* Chat Window */}
            <AnimatePresence>
                {isOpen && (
                    <motion.div
                        initial={{ opacity: 0, y: 50, scale: 0.9 }}
                        animate={{ opacity: 1, y: 0, scale: 1 }}
                        exit={{ opacity: 0, y: 50, scale: 0.9 }}
                        className="fixed inset-0 sm:inset-auto sm:bottom-6 sm:right-6 z-50"
                    >
                        {ChatContent}
                    </motion.div>
                )}
            </AnimatePresence>
        </>
    );
};
