import { useEffect, useState } from 'react';
import { useUser } from "@/app/context/UserContext";
import { clientCookieToken } from "@/app/utils/auth";
import checkLoginStatus from "@/app/utils/checkLoginStatus";


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
    const { user } = useUser();
    const [isLoggedIn, setIsLoggedIn] = useState(false);
    const [socket, setSocket] = useState<WebSocket | null>(null);

    useEffect(() => {
        setIsLoggedIn(checkLoginStatus());
        const clientToken = clientCookieToken();
        let socketInstance: WebSocket | null = null;

        const connectSocket = () => {
            if (clientToken) {
                socketInstance = new WebSocket("ws://localhost:8080/notify");

                socketInstance.onopen = () => {
                    console.log("Connected to notify WebSocket");
                };
                socketInstance.onmessage = (event) => {
                    const data: Event = JSON.parse(event.data);
                    if (data.type === "pending_follow_request") {
                        setNotifications(data.payload);
                        console.log("WS received user", data.payload)
                    }
                };
                socketInstance.onclose = (event) => {
                    console.log(`Disconnected: ${event.reason} (Code: ${event.code})`);
                    if (event.code !== 1000 && event.code !== 1001) {
                        console.log("Attempting to reconnect...");
                        setTimeout(connectSocket, 5000);
                    }
                };
            }
        }
        connectSocket();
        setSocket(socketInstance);
        return () => {
            if (socketInstance?.readyState === WebSocket.OPEN) {
                socketInstance.close(1000, "Component unmounted");
            }
        };
    }, [user, setNotifications]);
}


