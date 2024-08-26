//this will not be here forever it is a component of out future main page

"use client";
import React, { useEffect, useState } from "react";
import FieldInput from "../components/FieldInput";
import Button from "../components/Button";

const ChatMessage = () => {
    const [message, setMessage] = useState("");
    const [socket, setSocket] = useState<WebSocket | null>(null);

    class Payload {
        id: number;
        content: string;
        from: number;
        to: number;
        created: string;
        constructor(id: number, content: string, from: number, to: number) {
            this.id = id;
            this.content = content;
            this.from = from;
            this.to = to;
            this.created = new Date().toISOString();
        }
    }

    class Event {
        type: string;
        payload: Payload;

        constructor(type: string, payload: Payload) {
            this.type = type;
            this.payload = payload;
        }
    }

    useEffect(() => {
        const socketInstance = new WebSocket('ws://localhost:8080/ws');

        socketInstance.onopen = () => {
            console.log('Connected to WebSocket server');
        };

        socketInstance.onmessage = (event) => {
            console.log('Received message:', event.data);
            // Handle incoming messages
        };

        socketInstance.onerror = (error) => {
            console.error('Socket error', error);
        };

        socketInstance.onclose = () => {
            console.log('Disconnected from WebSocket server');
        };

        setSocket(socketInstance);

        // Cleanup on component unmount
        return () => {
            socketInstance.close();
        };
    }, []);

    const sendChatMessage = (e: React.FormEvent) => {
        const messageId: number = 1;
        const messageFrom: number = 2;
        const messageTo: number = 5;

        e.preventDefault();
        if (socket && message.trim() !== "") {
            const payload = new Payload(messageId, message, messageFrom, messageTo);
            const event = new Event('chat_message', payload);
            socket.send(JSON.stringify(event));
            setMessage(""); // Clear the input field after sending
        }
    };

    return (
        <div>
            <h1>Chat Component</h1>
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
}

export default ChatMessage;