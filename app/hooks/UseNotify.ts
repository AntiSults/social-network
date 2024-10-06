import { useEffect } from 'react';

interface User {
    ID: number;
    firstName: string;
    lastName: string;
}

interface Event {
    type: string;
    payload: User;
    token: string;
}

export const useNotificationWS = (setNotifications: (user: User) => void) => {

    useEffect(() => {

        const ws = new WebSocket("ws://localhost:8080/notify");
        ws.onopen = () => {
            console.log("Connected to notify WebSocket");
        };
        ws.onmessage = (event) => {
            const data: Event = JSON.parse(event.data);

            if (data.type === "pending_follow_request") {
                setNotifications(data.payload);
                console.log("WS received user", data.payload)
            }
        };
        ws.onclose = () => {
            console.log("Disconnected from WebSocket");
        };

        return () => {
            ws.close();
        };


    }, [setNotifications]);
}


