/**
 * Code Map: MessageInput – Chat input field with send button.
 * CID: 1.5.2-MessageInput
 */

import { Send } from 'lucide-react';

interface MessageInputProps {
    input: string;
    isTyping: boolean;
    onInputChange: (value: string) => void;
    onSend: () => void;
}

export const MessageInput: React.FC<MessageInputProps> = ({ input, isTyping, onInputChange, onSend }) => {
    return (
        <div className="p-4 border-t border-[#30363d] bg-[#1f242c]">
            <div className="relative">
                <input
                    type="text"
                    value={input}
                    onChange={(e) => onInputChange(e.target.value)}
                    onKeyPress={(e) => e.key === 'Enter' && onSend()}
                    placeholder="Type a message..."
                    className="w-full bg-[#0d1117] border border-[#30363d] rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-[#238636] transition-colors pr-10"
                />
                <button
                    onClick={onSend}
                    disabled={!input.trim() || isTyping}
                    className="absolute right-2 top-1/2 -translate-y-1/2 text-[#8b949e] hover:text-[#238636] disabled:opacity-30 p-1 transition-colors"
                >
                    <Send size={18} />
                </button>
            </div>
        </div>
    );
};
