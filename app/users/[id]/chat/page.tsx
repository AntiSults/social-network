"use client";

import React from "react";
import NavBar from "@/app/components/NavBar";
import ChatMessages from "@/app/components/ChatMessages";
import ChatInput from "@/app/components/ChatInput";
import useChat from "@/app/hooks/UseChat";
import { useUser } from "@/app/context/UserContext";

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

    return (
        <>
            <NavBar logged={isLoggedIn} />
            <div className="max-w-xl mx-auto p-4 bg-white shadow-md rounded-md">
                <h1 className="text-xl font-bold mb-4">Chat Component</h1>
                {isLoggedIn ? (
                    <>
                        <ChatMessages messages={messages} users={users} currentUser={currentUser} />
                        <ChatInput
                            message={message}
                            setMessage={setMessage}
                            onSubmit={sendChatMessage}
                            recipients={recipients}
                            selectedRecipient={selectedRecipient}
                            setSelectedRecipient={setSelectedRecipient}
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
