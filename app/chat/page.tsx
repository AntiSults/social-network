"use client";
import React, { useEffect, useState } from "react";
import FieldInput from "../components/FieldInput";
import Button from "../components/Button";

// Define types for message handling
interface Payload {
    id: number;
    content: string;
    fromUserID: number;
    toUserID: number;
    created: string;
}

interface Event {
    type: string;
    payload: Payload;
}

const ChatMessage = () => {
    const [message, setMessage] = useState("");
    const [messages, setMessages] = useState<Payload[]>([]); // State to hold messages
    const [socket, setSocket] = useState<WebSocket | null>(null);

    useEffect(() => {
        const socketInstance = new WebSocket("ws://localhost:8080/ws");

        socketInstance.onopen = () => {
            console.log("Connected to WebSocket server");
        };

        socketInstance.onmessage = (event) => {
            console.log("Received message:", event.data);

            // Parse the incoming message as an Event object
            const incomingEvent: Event = JSON.parse(event.data);

            // Add the new message to the messages state
            setMessages((prevMessages) => [...prevMessages, incomingEvent.payload]);
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
        const messageId: number = 1; // Example message ID
        const messageFromID: number = 2; // Example sender ID
        const messageToID: number = 5; // Example receiver ID

        e.preventDefault();
        if (socket && message.trim() !== "") {
            const payload: Payload = {
                id: messageId,
                content: message,
                fromUserID: messageFromID,
                toUserID: messageToID,
                created: new Date().toISOString(),
            };

            const event: Event = { type: "chat_message", payload };
            socket.send(JSON.stringify(event));
            setMessage(""); // Clear the input field after sending

            // Add outgoing message to the messages state
            setMessages((prevMessages) => [...prevMessages, payload]);
        }
    };

    return (
        <div>
            <h1>Chat Component</h1>
            <div className="chat-messages">
                {/* Display all messages */}
                {messages.map((msg, index) => (
                    <div key={index} className={msg.fromUserID === 2 ? "my-message" : "other-message"}>
                        <p>{msg.content}</p>
                        <small>{new Date(msg.created).toLocaleTimeString()}</small>
                    </div>
                ))}
            </div>
            <form onSubmit={sendChatMessage}>
                <FieldInput
                    name="Text:"
                    type="text"
                    placeholder="Push your imagination"
                    required={true}
                    value={message}
                    onChange={(e) => setMessage(e.target.value)}
                />
                <Button type="submit" name="Submit Message" />
            </form>
        </div>
    );
};

export default ChatMessage;
