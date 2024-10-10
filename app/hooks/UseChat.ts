import { useState, useEffect } from "react";
import { clientCookieToken } from "@/app/utils/auth";
import checkLoginStatus from "@/app/utils/checkLoginStatus";
import { useUser } from "@/app/context/UserContext";
import { Message, Recipient } from "@/app/utils/types";

interface User {
    ID: number;
    firstName: string;
    lastName: string;
}

interface Event {
    type: string;
    payload: {
        Message: Message[] | [];
        User: User[] | [];
        GroupName?: string;
    };
    token: string;
}

export const useChatAndNotify = (setNotifications?: (type: string, userNotify: User | null, group: string | null) => void, initialGroupId?: string) => {

    const [isLoggedIn, setIsLoggedIn] = useState(false);
    const [message, setMessage] = useState("");
    const [messages, setMessages] = useState<Message[]>([]);
    const [users, setUsers] = useState<Record<number, User>>({});
    const [socketChat, setSocketChat] = useState<WebSocket | null>(null);
    const [socketNotify, setSocketNotify] = useState<WebSocket | null>(null);
    const { user: currentUser } = useUser();

    const [groupId, setGroupId] = useState<string | null>(initialGroupId || null);
    const [recipients, setRecipients] = useState<Recipient[]>([]);
    const [selectedRecipient, setSelectedRecipient] = useState<number | null>(null);

    // Reset recipients and messages when groupId changes
    useEffect(() => {
        setRecipients([]);
        setMessages([]);
    }, [groupId]);

    useEffect(() => {
        setIsLoggedIn(checkLoginStatus());
        const clientToken = clientCookieToken();
        const newChatSocket = new WebSocket("ws://localhost:8080/ws");

        newChatSocket.onopen = () => {
            console.log("Connected to WS  server");

            if (clientToken) {
                const uploadRequest = groupId
                    ? {
                        type: "initial_group_upload",
                        payload: { groupId },
                        sessionToken: clientToken,
                    }
                    : {
                        type: "initial_upload",
                        sessionToken: clientToken,
                    };
                newChatSocket?.send(JSON.stringify(uploadRequest));
            }
        };

        newChatSocket.onmessage = (event) => {
            const incomingEvent: Event = JSON.parse(event.data);

            // Handle chat-related events
            if (incomingEvent.type === "initial_upload_response" || incomingEvent.type === "initial_group_upload_response") {
                if (Array.isArray(incomingEvent.payload?.Message)) {
                    setMessages(incomingEvent.payload.Message);
                }
                if (Array.isArray(incomingEvent.payload?.User)) {
                    const usersById = incomingEvent.payload.User.reduce((acc, user) => {
                        acc[user.ID] = user;
                        return acc;
                    }, {} as Record<number, User>);
                    setUsers(usersById);

                    const formattedRecipients: Recipient[] = incomingEvent.payload.User.map(user => ({
                        id: user.ID,
                        name: `${user.firstName} ${user.lastName}`,
                        type: groupId ? "group" : "user",
                    }));

                    setRecipients(formattedRecipients);
                }
            } else if (incomingEvent.type === "chat_message" || incomingEvent.type === "group_chat_message") {
                if (incomingEvent.payload?.Message) {
                    console.log("incoming message", incomingEvent.payload.Message)
                    setMessages(prevMessages => [...prevMessages, ...incomingEvent.payload.Message]);
                }
            }
        };

        newChatSocket.onerror = (error) => {
            console.error("Socket error", error);
        };

        newChatSocket.onclose = (event) => {
            console.log(`Disconnected: ${event.reason} (Code: ${event.code})`);
        };

        setSocketChat(newChatSocket);

        const newNotifySocket = new WebSocket("ws://localhost:8080/notify");

        newNotifySocket.onopen = () => {
            console.log("Connected to WS Notify server");
        };

        newNotifySocket.onmessage = (event) => {
            const incomingEvent: Event = JSON.parse(event.data);

            // Handle notification-related events
            if (incomingEvent.type === "Pending-follow-request" && incomingEvent.payload?.User) {
                const { ID, firstName, lastName } = incomingEvent.payload.User[0];
                const filteredUser: User = { ID, firstName, lastName };
                if (setNotifications) {
                    setNotifications(incomingEvent.type, filteredUser, null);
                }
            } else if (
                (incomingEvent.type === "Group-Invite-Notification" || incomingEvent.type === "Group-Join-Request") &&
                incomingEvent.payload?.User && incomingEvent.payload.GroupName
            ) {
                const { ID, firstName, lastName } = incomingEvent.payload.User[0];
                const filteredUser: User = { ID, firstName, lastName };
                if (setNotifications) {
                    setNotifications(incomingEvent.type, filteredUser, incomingEvent.payload.GroupName);
                }
            } else if (incomingEvent.type === "New-Group-Event") {
                if (setNotifications) {
                    setNotifications(incomingEvent.type, null, incomingEvent.payload.GroupName || null)
                }
            }
        };

        newNotifySocket.onerror = (error) => {
            console.error("Socket error", error);
        };

        newNotifySocket.onclose = (event) => {
            console.log(`Disconnected: ${event.reason} (Code: ${event.code})`);
        };

        setSocketNotify(newNotifySocket);
        // Cleanup on component unmount
        return () => {
            if (newChatSocket) newChatSocket.close();
            if (newNotifySocket) newNotifySocket.close();
        };
    }, [setNotifications, groupId]);

    const sendChatMessage = (e: React.FormEvent) => {
        e.preventDefault();

        const clientToken = clientCookieToken();
        if (!clientToken) {
            console.error("No session token found.");
            return;
        }

        if (socketChat && message.trim() !== "" && currentUser && selectedRecipient !== null) {
            const payload: Message = {
                content: message,
                fromUserID: currentUser.ID,
                toUserID: selectedRecipient,
                created: new Date().toISOString(),
                groupID: groupId !== null ? +groupId : 0,
            };

            const event: Event = {
                type: groupId ? "group_chat_message" : "chat_message",
                payload: { Message: [payload], User: [], GroupName: "" },
                token: clientToken,
            };

            socketChat.send(JSON.stringify(event));
            setMessage("");
            setMessages((prevMessages) => [...prevMessages, payload]);
        }
    };

    return {
        isLoggedIn,
        messages,
        sendChatMessage,
        message,
        setMessage,
        users,
        recipients,
        selectedRecipient,
        setSelectedRecipient,
        groupId,
        setGroupId,
    };
};

