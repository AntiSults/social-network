import { useEffect, useState } from 'react';
import { useUser } from '@/app/context/UserContext';

interface User {
    ID: number;
    firstName: string;
    lastName: string;
}
interface Data {
    User: User | null;
    GroupName: string | null;
}
interface Event {
    type: string;
    payload: Data;
    token: string;
}

export const useNotificationWS = (setNotifications: (type: string, user: User | null, group: string | null) => void) => {
    const { user } = useUser();
    const [_, setSocket] = useState<WebSocket | null>(null);
    useEffect(() => {
        let socketInstance: WebSocket | null = null;
        if (user) {
            const connectSocket = () => {

                const socketInstance = new WebSocket('ws://localhost:8080/notify');

                socketInstance.onopen = () => {
                    console.log('Connected to notify WebSocket');
                };
                socketInstance.onmessage = (event) => {
                    try {
                        const data: Event = JSON.parse(event.data);

                        if (data.type === 'Pending-follow-request' && data.payload.User) {
                            const { ID, firstName, lastName } = data.payload.User;
                            const filteredUser: User = { ID, firstName, lastName };
                            setNotifications(data.type, filteredUser, null);
                        } else if (
                            (data.type === 'Group-Invite-Notification' || data.type === 'Group-Join-Request') &&
                            data.payload.User && data.payload.GroupName
                        ) {
                            const { ID, firstName, lastName } = data.payload.User;
                            const filteredUser: User = { ID, firstName, lastName };
                            setNotifications(data.type, filteredUser, data.payload.GroupName);
                        } else if (data.type === 'New-Group-Event') {
                            setNotifications(data.type, null, data.payload.GroupName);
                        }
                    } catch (error) {
                        console.error('Failed to parse WebSocket message:', error);
                    }
                };
                socketInstance.onerror = (error) => {
                    console.error('WebSocket error:', error);
                };
                socketInstance.onclose = (event) => {
                    console.log(`WebSocket closed: ${event.reason} (Code: ${event.code})`);
                    if (event.code !== 1000 && event.code !== 1001) {
                        console.log('Attempting to reconnect...');
                        setTimeout(connectSocket, 5000); // Reconnect after 5 seconds
                    }
                };
            };

            connectSocket();
            setSocket(socketInstance);

        }


    }, [user, setNotifications]);
};
