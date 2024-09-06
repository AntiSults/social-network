"use client";
import React, { useEffect, useState } from "react";
import FieldInput from "../components/FieldInput";
import Button from "../components/Button";
import { clientCookieToken } from "../utils/auth"; // Utility function to get the token from the cookie

interface Payload {
    id: number;
    content: string;
    fromUserID: number;
    toUserID: number;
    created: string;
}

interface Event {
    type: string;
    payload: Payload | Payload[]; // Can be a single message or an array of messages
    token: string;
}

const ChatMessage = () => {
    const [message, setMessage] = useState("");
    const [messages, setMessages] = useState<Payload[]>([]); // State to hold messages
    const [socket, setSocket] = useState<WebSocket | null>(null);

    useEffect(() => {
        // Get the token directly from client-side cookie
        const clientToken = clientCookieToken();

        const socketInstance = new WebSocket("ws://localhost:8080/ws");

        socketInstance.onopen = () => {
            console.log("Connected to WebSocket server");
            if (clientToken) {
                // Prepare and send the initial upload request with token
                const uploadRequest = {
                    type: 'initial_upload',
                    payload: {}, // No specific payload needed, just token for identification
                    sessionToken: clientToken,
                };
                socketInstance.send(JSON.stringify(uploadRequest));

            }
        };

        socketInstance.onmessage = (event) => {
            console.log("Received message:", event.data);

            const incomingEvent: Event = JSON.parse(event.data);

            // Check if the incoming event is the initial messages load
            if (incomingEvent.type === 'initial_upload_response' && Array.isArray(incomingEvent.payload)) {
                // Set all messages received from the server
                setMessages(incomingEvent.payload);
            } else if (incomingEvent.type === 'chat_message') {
                // Handle individual incoming chat messages
                setMessages((prevMessages) => [...prevMessages, incomingEvent.payload as Payload]);
            }
        };

        socketInstance.onerror = (error) => {
            console.error("Socket error", error);
        };

        socketInstance.onclose = () => {
            console.log("Disconnected from WebSocket server");
        };

        setSocket(socketInstance);

        // Cleanup on component unmount
        return () => {
            socketInstance.close();
        };
    }, []);

    const sendChatMessage = (e: React.FormEvent) => {
        e.preventDefault();

        const clientToken = clientCookieToken(); // Read the token directly from the cookie

        if (!clientToken) {
            console.error("No session token found.");
            return;
        }

        const messageId: number = 1; // Example message ID
        const messageFromID: number = 2; // Example sender ID
        const messageToID: number = 5; // Example receiver ID

        if (socket && message.trim() !== "") {
            const payload: Payload = {
                id: messageId,
                content: message,
                fromUserID: messageFromID,
                toUserID: messageToID,
                created: new Date().toISOString(),
            };

            const event: Event = {
                type: "chat_message",
                payload,
                token: clientToken,
            };

            socket.send(JSON.stringify(event));
            setMessage(""); // Clear the input field after sending

            setMessages((prevMessages) => [...prevMessages, payload]);
        }
    };

    return (
        <div className="max-w-xl mx-auto p-4 bg-white shadow-md rounded-md">
            <h1 className="text-xl font-bold mb-4">Chat Component</h1>
            <div className="chat-messages mb-4 max-h-96 overflow-y-auto border border-gray-300 rounded-md p-2">
                {/* Display all messages */}
                {messages.map((msg, index) => (
                    <div
                        key={index}
                        className={`p-2 my-1 rounded-md ${msg.fromUserID === 2
                            ? "bg-blue-100 text-right self-end"
                            : "bg-gray-100 text-left"
                            }`}
                    >
                        <p className="text-sm">{msg.content}</p>
                        <small className="text-xs text-gray-500">
                            {new Date(msg.created).toLocaleTimeString()}
                        </small>
                    </div>
                ))}
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
    );
};

export default ChatMessage;
