"use client";

import React, { useEffect, useState } from "react";
import FieldInput from "../components/FieldInput";
import Button from "../components/Button";
import NavBar from "../components/NavBar";
import checkLoginStatus from "../utils/checkLoginStatus";
import { clientCookieToken } from "../utils/auth";
import { useUser } from "../context/UserContext";

interface Message {
    id: number;
    content: string;
    fromUserID: number;
    toUserID: number;
    created: string;
}

interface User {
    ID: number;
    firstName: string;
    lastName: string;
}

interface Payload {
    Message: Message[];
    User: User[];
}

interface Event {
    type: string;
    payload: Payload;
    token: string;
}

const ChatMessage = () => {
    const [isLoggedIn, setIsLoggedIn] = useState(false);
    const [message, setMessage] = useState("");
    const [messages, setMessages] = useState<Message[]>([]);
    const [users, setUsers] = useState<Record<number, User>>({});
    const [socket, setSocket] = useState<WebSocket | null>(null);
    const { user: currentUser } = useUser();

    useEffect(() => {
        setIsLoggedIn(checkLoginStatus());
        const clientToken = clientCookieToken();
        const socketInstance = new WebSocket("ws://localhost:8080/ws");

        socketInstance.onopen = () => {
            console.log("Connected to WebSocket server");
            if (clientToken) {
                const uploadRequest = {
                    type: "initial_upload",
                    payload: {},
                    sessionToken: clientToken,
                };
                socketInstance.send(JSON.stringify(uploadRequest));
            }
        };

        socketInstance.onmessage = (event) => {
            console.log("Received message:", event.data);

            const incomingEvent: Event = JSON.parse(event.data);

            if (
                incomingEvent.type === "initial_upload_response" &&
                incomingEvent.payload &&
                Array.isArray(incomingEvent.payload.Message) &&
                Array.isArray(incomingEvent.payload.User)
            ) {
                // Set all messages received from the server
                setMessages(incomingEvent.payload.Message);

                // Create a mapping of user ID to user details for quick lookup
                const usersById = incomingEvent.payload.User.reduce((acc, user) => {
                    acc[user.ID] = user;
                    return acc;
                }, {} as Record<number, User>);

                setUsers(usersById);
            } else if (incomingEvent.type === "chat_message" && incomingEvent.payload.Message) {
                // Handle individual incoming chat messages
                setMessages((prevMessages) => [...prevMessages, ...incomingEvent.payload.Message]);
            }
        };

        socketInstance.onerror = (error) => {
            console.error("Socket error", error);
        };

        socketInstance.onclose = () => {
            console.log("Disconnected from WebSocket server");
        };

        setSocket(socketInstance);

        return () => {
            socketInstance.close();
        };
    }, []);

    const sendChatMessage = (e: React.FormEvent) => {
        e.preventDefault();

        const clientToken = clientCookieToken();

        if (!clientToken) {
            console.error("No session token found.");
            return;
        }

        if (socket && message.trim() !== "" && currentUser) {
            const messageId: number = 1;
            const messageFromID: number = currentUser.ID
            const messageToID: number = 3;
            const payload: Message = {
                id: messageId,
                content: message,
                fromUserID: messageFromID,
                toUserID: messageToID,
                created: new Date().toISOString(),
            };

            const event: Event = {
                type: "chat_message",
                payload: { Message: [payload], User: [] },
                token: clientToken,
            };

            socket.send(JSON.stringify(event));
            setMessage(""); // Clear the input field after sending
            setMessages((prevMessages) => [...prevMessages, payload]);
        }
    };

    return (
        <>
            <NavBar logged={isLoggedIn}></NavBar>
            <div className="max-w-xl mx-auto p-4 bg-white shadow-md rounded-md">
                <h1 className="text-xl font-bold mb-4">Chat Component</h1>
                <div className="chat-messages mb-4 max-h-96 overflow-y-auto border border-gray-300 rounded-md p-2">
                    {/* Display all messages */}
                    {messages.map((msg, index) => {
                        const sender = users[msg.fromUserID];
                        const senderName = sender ? `${sender.firstName} ${sender.lastName}` : "Unknown User";
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
                <form onSubmit={sendChatMessage} className="flex items-center space-x-2">
                    <FieldInput
                        name="Text:"
                        type="text"
                        placeholder="Push your imagination"
                        required={true}
                        value={message}
                        onChange={(e) => setMessage(e.target.value)}
                        className="flex-grow p-2 border border-gray-300 rounded-md"
                    />
                    <Button
                        type="submit"
                        name="Submit Message"
                        className="p-2 bg-blue-500 text-white rounded-md"
                    />
                </form>
            </div>
        </>
    );
};

export default ChatMessage;
