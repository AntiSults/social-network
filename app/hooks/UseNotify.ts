import { useEffect, useRef } from 'react';
import checkLoginStatus from "@/app/utils/checkLoginStatus";

interface User {
    ID: number;
    firstName: string;
    lastName: string;
}
interface Data {
    User: User;
    GroupName: string | null;
}
interface Event {
    type: string;
    payload: Data;
    token: string;
}

export const useNotificationWS = (setNotifications: (user: User, type: string, group: string | null) => void) => {
    const socketRef = useRef<WebSocket | null>(null);
    const reconnectTimeout = useRef<NodeJS.Timeout | null>(null); // For reconnect logic

    useEffect(() => {
        const loggedIn = checkLoginStatus();

        if (loggedIn && !socketRef.current) {
            const connectSocket = () => {
                if (reconnectTimeout.current) {
                    clearTimeout(reconnectTimeout.current);
                    reconnectTimeout.current = null;
                }
                const socket = new WebSocket("ws://localhost:8080/notify");
                socketRef.current = socket;

                socket.onopen = () => {
                    console.log("Connected to notify WebSocket");
                };
                socket.onmessage = (event) => {
                    const data: Event = JSON.parse(event.data);
                    if (data.type === "Pending-follow-request") {
                        const { ID, firstName, lastName } = data.payload.User;
                        const filteredUser: User = { ID, firstName, lastName };
                        setNotifications(filteredUser, data.type, null);
                    } else if (data.type === "Group-Invite-Notification") {
                        const { ID, firstName, lastName } = data.payload.User;
                        const filteredUser: User = { ID, firstName, lastName };
                        setNotifications(filteredUser, data.type, data.payload.GroupName);
                    }
                };
                socket.onclose = (event) => {
                    console.log(`Disconnected: ${event.reason} (Code: ${event.code})`);
                    if (event.code !== 1000 && event.code !== 1001) {
                        console.log("Attempting to reconnect...");
                        reconnectTimeout.current = setTimeout(connectSocket, 5000); // Reconnect after delay
                    }
                };
                socket.onerror = (error) => {
                    console.error("WebSocket error:", error);
                };
            };
            connectSocket();

            return () => {
                if (socketRef.current?.readyState === WebSocket.OPEN) {
                    socketRef.current.close(1000, "Component unmounted");
                }
                if (reconnectTimeout.current) {
                    clearTimeout(reconnectTimeout.current);
                }
            };
        }
        return () => {
            if (socketRef.current?.readyState === WebSocket.OPEN) {
                socketRef.current.close(1000, "Component unmounted");
            }
            if (reconnectTimeout.current) {
                clearTimeout(reconnectTimeout.current);
            }
        };
    }, [setNotifications]); // Depend only on setNotifications
};
