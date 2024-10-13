import React from 'react';
import { Message } from '@/app/utils/types';

interface User {
    ID: number;
    firstName: string;
    lastName: string;
}

interface Props {
    messages: Message[];
    users: Record<number, User>;
    currentUser: User | null;
    groupId: string | null;
    setGroupId: React.Dispatch<React.SetStateAction<string | null>>;
}

const ChatMessages: React.FC<Props> = ({ messages, users, currentUser, groupId, setGroupId }) => {
    const handleBackToRegularChat = () => {
        setGroupId(null); // Reset the groupId to go back to regular chat
    };

    return (
        <div>
            {/* Show "Back to Regular Chat" button only when in group chat */}
            {groupId && (
                <button
                    onClick={handleBackToRegularChat}
                    className="mb-2 bg-red-500 hover:bg-red-600 text-white py-2 px-4 rounded-md"
                >
                    Back to Regular Chat
                </button>
            )}

            <div className="chat-messages mb-4 max-h-96 overflow-y-auto border border-gray-300 rounded-md p-2">
                {messages.map((msg, index) => {
                    const sender = users[msg.fromUserID];
                    const senderName = currentUser && msg.fromUserID === currentUser.ID
                        ? "Me"
                        : `${sender?.firstName || 'Unknown'} ${sender?.lastName || 'User'}`;

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
        </div>
    );
};

export default ChatMessages;

