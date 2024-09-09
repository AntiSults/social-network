import { useState, useEffect } from "react";
import { clientCookieToken } from "../utils/auth";
import checkLoginStatus from "../utils/checkLoginStatus";
import { useUser } from "../context/UserContext";

interface Message {
  id: number;
  content: string;
  fromUserID: number;
  toUserID: number;
  created: string;
}

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

const useChat = () => {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [message, setMessage] = useState("");
  const [messages, setMessages] = useState<Message[]>([]);
  const [users, setUsers] = useState<Record<number, User>>({});
  const [socket, setSocket] = useState<WebSocket | null>(null);

  const { user: currentUser } = useUser();

  useEffect(() => {
    setIsLoggedIn(checkLoginStatus());
    const clientToken = clientCookieToken();
    const socketInstance = new WebSocket("ws://localhost:8080/ws");

    socketInstance.onopen = () => {
      console.log("Connected to WebSocket server");
      if (clientToken) {
        const uploadRequest = {
          type: "initial_upload",
          payload: {},
          sessionToken: clientToken,
        };
        socketInstance.send(JSON.stringify(uploadRequest));
      }
    };

    socketInstance.onmessage = (event) => {
      console.log("Received message:", event.data);
      const incomingEvent: Event = JSON.parse(event.data);

      if (
        incomingEvent.type === "initial_upload_response" &&
        incomingEvent.payload &&
        Array.isArray(incomingEvent.payload.Message) &&
        Array.isArray(incomingEvent.payload.User)
      ) {
        setMessages(incomingEvent.payload.Message);

        const usersById = incomingEvent.payload.User.reduce((acc, user) => {
          acc[user.ID] = user;
          return acc;
        }, {} as Record<number, User>);

        setUsers(usersById);
      } else if (incomingEvent.type === "chat_message" && incomingEvent.payload.Message) {
        setMessages((prevMessages) => [...prevMessages, ...incomingEvent.payload.Message]);
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
  }, []);

  const sendChatMessage = (e: React.FormEvent) => {
    e.preventDefault();

    const clientToken = clientCookieToken();
    if (!clientToken) {
      console.error("No session token found.");
      return;
    }

    if (socket && message.trim() !== "" && currentUser) {
      const payload: Message = {
        id: 1,
        content: message,
        fromUserID: currentUser.ID,
        toUserID: 3, // Example toUserID
        created: new Date().toISOString(),
      };

      const event: Event = {
        type: "chat_message",
        payload: { Message: [payload], User: [] },
        token: clientToken,
      };

      socket.send(JSON.stringify(event));
      setMessage("");
      setMessages((prevMessages) => [...prevMessages, payload]);
    }
  };

  return { isLoggedIn, messages, sendChatMessage, message, setMessage, users };
};

export default useChat;
