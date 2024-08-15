//this will not be here forever it is a component of out future main page
"use client";
import React, { useEffect, useState } from "react";
import FieldInput from "../components/FieldInput";
import Button from "../components/Button";



const ChatMessage = () => {
    const [message, setMessageText] = useState("");
    //const [error, setError] = useState("");
    class Payload {
        content: string;
        created: string;
        constructor(content: string) {
            this.content = content;
            this.created = new Date().toISOString();;
        }
    }

    class Event {
        type: string;
        payload: object;

        constructor(type: string, payload: object) {
            this.type = type;
            this.payload = payload;
        }
    }
    useEffect(() => {
        const socket = new WebSocket('ws://localhost:8080/ws');

        socket.onopen = () => {
            console.log('Connected to WebSocket server');
        };

        socket.onmessage = (event) => {
            console.log('Received message:', event.data);
            // Handle incoming messages
        };

        socket.onerror = (error) => {
            console.log('Socket error', error);
        }

        socket.onclose = () => {
            console.log('Disconnected from WebSocket server');
        };

        // Cleanup on component unmount
        return () => {
            socket.close();
        };
    }, []);

    const sendChatMessage = () => {
        let event = new Event('chat message', new Payload('test message'))
    }

    return (
        <div>
            <h1>Chat Component</h1>
            <form onSubmit={sendChatMessage}>
                <FieldInput
                    name="Text:"
                    type="text"
                    placeholder="push your imagination"
                    required={true}
                    value={message}
                    onChange={(e) => setMessageText(e.target.value)}
                />
                <Button type="submit" name="Submit Message" />

            </form>
        </div>
    )
}


export default ChatMessage