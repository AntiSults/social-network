"use client";

import React from "react";
import NavBar from "../components/NavBar";
import ChatMessages from "../components/ChatMessages";
import ChatInput from "../components/ChatInput";
import useChat from "../hooks/UseChat";
import { useUser } from "../context/UserContext";

const ChatMessage = () => {
    const { isLoggedIn, messages, sendChatMessage, message, setMessage, users } = useChat();
    const { user: currentUser } = useUser();

    return (
        <>
            <NavBar logged={isLoggedIn} />
            <div className="max-w-xl mx-auto p-4 bg-white shadow-md rounded-md">
                <h1 className="text-xl font-bold mb-4">Chat Component</h1>
                {isLoggedIn ? (
                    <>
                        <ChatMessages messages={messages} users={users} currentUser={currentUser} />
                        <ChatInput message={message} setMessage={setMessage} onSubmit={sendChatMessage} />
                    </>
                ) : (
                    <p className="text-center text-gray-600">Please login for chatting!</p>
                )}
            </div>
        </>

    );
};

export default ChatMessage;
