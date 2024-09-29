"use client";

import React, { useState } from "react";
import NavBar from "@/app/components/NavBar";
import ChatMessages from "@/app/components/ChatMessages";
import ChatInput from "@/app/components/ChatInput";
import useChat from "@/app/hooks/UseChat";
import { useUser } from "@/app/context/UserContext";
import GroupList from "@/app/components/GroupList"; // Add this import

const ChatMessage = () => {
    const {
        isLoggedIn,
        messages,
        sendChatMessage,
        message,
        setMessage,
        users,
        recipients,
        selectedRecipient,
        setSelectedRecipient,
    } = useChat();

    const { user: currentUser } = useUser();

    const [selectedGroup, setSelectedGroup] = useState<number | null>(null); // Add state for selected group

    const handleGroupSelect = (groupId: number) => {
        setSelectedGroup(groupId);
        setSelectedRecipient(null); // Reset the selected recipient if switching to group chat
    };

    return (
        <>
            <NavBar logged={isLoggedIn} />
            <div className="max-w-xl mx-auto p-4 bg-white shadow-md rounded-md">
                <h1 className="text-xl font-bold mb-4">Chat Component</h1>
                {isLoggedIn ? (
                    <>
                        {/* Add GroupList for selecting a group */}
                        <GroupList onSelectGroup={handleGroupSelect} actionType="chat" />

                        {/* Display Chat Messages */}
                        <ChatMessages messages={messages} users={users} currentUser={currentUser} />

                        {/* Display input for group or individual chat */}
                        {selectedGroup !== null ? (
                            <p className="text-sm">Group Chat: {selectedGroup}</p>
                        ) : (
                            <ChatInput
                                message={message}
                                setMessage={setMessage}
                                onSubmit={sendChatMessage}
                                recipients={recipients}
                                selectedRecipient={selectedRecipient}
                                setSelectedRecipient={setSelectedRecipient}
                            />
                        )}
                    </>
                ) : (
                    <p className="text-center text-gray-600">Please login for chatting!</p>
                )}
            </div>
        </>
    );
};

export default ChatMessage;
