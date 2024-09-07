"use client";

import React, { createContext, useContext, useState, useEffect, ReactNode } from "react";

interface User {
    ID: number;
    firstName: string;
    lastName: string;
    email: string;
    dob: string;
    nickname?: string;
    avatarPath?: string;
    profileVisibility?: string;
    aboutMe?: string;
}

// Define the context value type
interface UserContextType {
    user: User | null;
    setUser: React.Dispatch<React.SetStateAction<User | null>>;
}

// Create the UserContext
const UserContext = createContext<UserContextType | undefined>(undefined);

// Custom hook to use the UserContext
export const useUser = () => {
    const context = useContext(UserContext);
    if (!context) {
        throw new Error("useUser must be used within a UserProvider");
    }
    return context;
};

// UserProvider component to provide the context
export const UserProvider = ({ children }: { children: ReactNode }) => {
    const [user, setUser] = useState<User | null>(null);

    // Example: Fetch user data on component mount
    useEffect(() => {
        const fetchUser = async () => {
            try {
                const response = await fetch("http://localhost:8080/getUserData", {
                    method: "GET",
                    credentials: "include",
                });
                if (response.ok) {
                    const userData = await response.json();
                    setUser(userData); // Update the user state
                }
            } catch (error) {
                console.error("Failed to fetch user data", error);
            }
        };

        fetchUser();
    }, []);

    return (
        <UserContext.Provider value={{ user, setUser }}>
            {children}
        </UserContext.Provider>
    );
};