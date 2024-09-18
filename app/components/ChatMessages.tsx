import React from "react";

interface Message {
    id?: number;
    content: string;
    fromUserID: number;
    toUserID: number[];
    created: string;
}

interface User {
    ID: number;
    firstName: string;
    lastName: string;
}

interface ChatMessagesProps {
    messages: Message[];
    users: Record<number, User>;
    currentUser: User | null;
}

const ChatMessages: React.FC<ChatMessagesProps> = ({ messages, users, currentUser }) => {
    return (
        <div className="chat-messages mb-4 max-h-96 overflow-y-auto border border-gray-300 rounded-md p-2">
            {messages.map((msg, index) => {
                const sender = users[msg.fromUserID];
                const senderName = currentUser && msg.fromUserID === currentUser.ID ? "Me" : `${sender?.firstName || 'Unknown'} ${sender?.lastName || 'User'}`;
                // Handle individual vs group message logic (for future use if needed)
                const isGroupMessage = msg.toUserID.length > 1;

                return (
                    <div
                        key={index}
                        className={`p-2 my-1 rounded-md ${currentUser && msg.fromUserID === currentUser.ID
                            ? "bg-blue-100 text-right self-end"
                            : "bg-gray-100 text-left"
                            }`}
                    >
                        <p className="text-sm font-bold">{senderName}</p>
                        <p className="text-sm">{msg.content}</p>
                        <small className="text-xs text-gray-500">
                            {new Date(msg.created).toLocaleString()}
                        </small>
                    </div>
                );
            })}
        </div>
    );
};

export default ChatMessages;
