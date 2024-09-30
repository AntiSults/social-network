"use client";

import React, { useState } from "react";
import NavBar from "@/app/components/NavBar";
import ChatMessages from "@/app/components/ChatMessages";
import ChatInput from "@/app/components/ChatInput";
import useChat from "@/app/hooks/UseChat";
import { useUser } from "@/app/context/UserContext";
import GroupList from "@/app/components/GroupList";

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
        groupId,
        setGroupId,
    } = useChat();

    const { user: currentUser } = useUser();

    const handleGroupSelect = (groupId: number) => {
        setGroupId(groupId.toString()); // Update the groupId in chat state
    };

    return (
        <>
            <NavBar logged={isLoggedIn} />
            <div className="max-w-xl mx-auto p-4 bg-white shadow-md rounded-md">
                <h1 className="text-xl font-bold mb-4">Chat Component</h1>

                <GroupList onSelectGroup={handleGroupSelect} actionType="chat" />

                {isLoggedIn ? (
                    <>
                        <ChatMessages
                            messages={messages}
                            users={users}
                            currentUser={currentUser}
                            groupId={groupId} // Pass groupId prop
                            setGroupId={setGroupId}
                        />
                        <ChatInput
                            message={message}
                            setMessage={setMessage}
                            onSubmit={sendChatMessage} // Ensure this is connected
                            recipients={recipients}
                            selectedRecipient={selectedRecipient}
                            setSelectedRecipient={setSelectedRecipient}
                            groupId={groupId} // Pass groupId prop to ChatInput
                        />
                    </>
                ) : (
                    <p className="text-center text-gray-600">Please login for chatting!</p>
                )}
            </div>
        </>
    );
};

export default ChatMessage;
