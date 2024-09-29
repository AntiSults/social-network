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
    Message: Message[];
    User: User[];
  };
  token: string;
}

const useChat = (initialGroupId?: string) => {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [message, setMessage] = useState("");
  const [messages, setMessages] = useState<Message[]>([]);
  const [users, setUsers] = useState<Record<number, User>>({});
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const { user: currentUser } = useUser();

  // Add groupId state here
  const [groupId, setGroupId] = useState<string | null>(initialGroupId || null);

  const [recipients, setRecipients] = useState<Recipient[]>([]);
  const [selectedRecipient, setSelectedRecipient] = useState<number | null>(null);

  useEffect(() => {
    setIsLoggedIn(checkLoginStatus());
    const clientToken = clientCookieToken();
    const socketInstance = new WebSocket("ws://localhost:8080/ws");

    socketInstance.onopen = () => {
      console.log("Connected to WebSocket server");
      if (clientToken) {
        const uploadRequest = {
          type: groupId ? "initial_group_upload" : "initial_upload",
          payload: { groupId: groupId || null },
          sessionToken: clientToken,
        };
        socketInstance.send(JSON.stringify(uploadRequest));
      }
    };

    socketInstance.onmessage = (event) => {
      console.log("Received message:", event.data);
      const incomingEvent: Event = JSON.parse(event.data);

      if (incomingEvent.type === "initial_upload_response" || incomingEvent.type === "initial_group_upload_response") {
        if (Array.isArray(incomingEvent.payload.Message)) {
          setMessages(incomingEvent.payload.Message);
        }

        if (Array.isArray(incomingEvent.payload.User)) {
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
      }

      if (incomingEvent.type === "chat_message" || incomingEvent.type === "group_chat_message") {
        if (incomingEvent.payload.Message) {
          setMessages((prevMessages) => [...prevMessages, ...incomingEvent.payload.Message]);
        }
      }
    };

    socketInstance.onerror = (error) => {
      console.error("Socket error", error);
    };

    socketInstance.onclose = () => {
      console.log("Disconnected from WebSocket server");
    };

    setSocket(socketInstance);

    return () => {
      socketInstance.close();
    };
  }, [groupId]);

  const sendChatMessage = (e: React.FormEvent) => {
    e.preventDefault();

    const clientToken = clientCookieToken();
    if (!clientToken) {
      console.error("No session token found.");
      return;
    }

    if (socket && message.trim() !== "" && currentUser && selectedRecipient !== null) {
      const payload: Message = {
        content: message,
        fromUserID: currentUser.ID,
        toUserID: selectedRecipient,
        created: new Date().toISOString(),
      };

      const event: Event = {
        type: groupId ? "group_chat_message" : "chat_message",
        payload: { Message: [payload], User: [] },
        token: clientToken,
      };

      socket.send(JSON.stringify(event));
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
    setGroupId, // Return setGroupId to be able to update it
  };
};

export default useChat;
