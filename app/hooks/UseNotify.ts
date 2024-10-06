import { useEffect, useState } from 'react';

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
    const [socket, setSocket] = useState<WebSocket | null>(null);

    useEffect(() => {
        let socketInstance: WebSocket | null = null;

        const connectSocket = () => {
            socketInstance = new WebSocket("ws://localhost:8080/notify");

            socketInstance.onopen = () => {
                console.log("Connected to notify WebSocket");
            };
            socketInstance.onmessage = (event) => {
                const data: Event = JSON.parse(event.data);
                if (data.type === "pending_follow_request") {
                    const { ID, firstName, lastName } = data.payload;
                    const filteredUser: User = { ID, firstName, lastName };
                    setNotifications(filteredUser);
                }
            };
            socketInstance.onclose = (event) => {
                console.log(`Disconnected: ${event.reason} (Code: ${event.code})`);
                if (event.code !== 1000 && event.code !== 1001) {
                    console.log("Attempting to reconnect...");
                    setTimeout(connectSocket, 5000);
                }
            };
        };

        connectSocket();
        setSocket(socketInstance);

        return () => {
            if (socketInstance?.readyState === WebSocket.OPEN) {
                socketInstance.close(1000, "Component unmounted");
            }
        };
    }, [setNotifications]);
};
