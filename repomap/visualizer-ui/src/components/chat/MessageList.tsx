/**
 * Code Map: MessageList – Renders a scrollable list of chat messages with markdown support.
 * CID: 1.5.2-MessageList
 */

import { Bot, Loader2 } from 'lucide-react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { type Message } from '../../hooks/useChat';

interface MessageListProps {
    messages: Message[];
    isTyping: boolean;
    scrollRef: React.RefObject<HTMLDivElement | null>;
}

export const MessageList: React.FC<MessageListProps> = ({ messages, isTyping, scrollRef }) => {
    return (
        <div ref={scrollRef} className="flex-1 overflow-y-auto p-4 space-y-4 custom-scrollbar">
            {messages.length === 0 && (
                <div className="h-full flex flex-col items-center justify-center text-[#8b949e] opacity-50 text-center px-8">
                    <Bot size={48} strokeWidth={1} className="mb-4" />
                    <p className="text-sm">Hi! I can help you refine the plan. Just ask me to add tasks or update files.</p>
                </div>
            )}
            {messages.map((m, i) => (
                <div key={i} className={`flex ${m.role === 'user' ? 'justify-end' : 'justify-start'}`}>
                    <div className={`max-w-[85%] p-3 rounded-2xl text-sm ${m.role === 'user'
                        ? 'bg-[#238636] text-white rounded-tr-none'
                        : 'bg-[#21262d] text-[#c9d1d9] border border-[#30363d] rounded-tl-none'
                        }`}>
                        <ReactMarkdown
                            remarkPlugins={[remarkGfm]}
                            components={{
                                code({ node, inline, className, children, ...props }: any) {
                                    return !inline ? (
                                        <div className="bg-[#0d1117] p-2 rounded my-2 overflow-x-auto border border-[#30363d]">
                                            <code className={className} {...props}>
                                                {children}
                                            </code>
                                        </div>
                                    ) : (
                                        <code className="bg-[#0d1117] px-1 rounded border border-[#30363d]" {...props}>
                                            {children}
                                        </code>
                                    )
                                }
                            }}
                        >
                            {m.content || ''}
                        </ReactMarkdown>
                        {isTyping && i === messages.length - 1 && !m.content && <Loader2 className="animate-spin" size={16} />}
                    </div>
                </div>
            ))}
        </div>
    );
};
