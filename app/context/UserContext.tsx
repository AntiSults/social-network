// UserContext.tsx
"use client";

import React, { createContext, useContext, useEffect, useState } from "react";

interface User {
    ID: number;
    email: string;
    firstName: string;
    lastName: string;
    dob: string;
    nickname?: string;
    aboutMe?: string;
    avatarPath?: string | null;
    profileVisibility?: "public" | "private";
}

interface UserContextProps {
    user: User | null; // Current user
    selectedUser: User | null; // Selected user from search results
    setUser: React.Dispatch<React.SetStateAction<User | null>>;
    setSelectedUser: React.Dispatch<React.SetStateAction<User | null>>; // New setter for selected user
}


// Create the context
const UserContext = createContext<UserContextProps | undefined>(undefined);

// Hook to use the user context
export const useUser = () => {
    const context = useContext(UserContext);
    if (!context) {
        throw new Error("useUser must be used within a UserProvider");
    }
    return context;
};

// Provider component
export const UserProvider = ({ children }: { children: React.ReactNode }) => {
    const [user, setUser] = useState<User | null>(null);
    const [selectedUser, setSelectedUser] = useState<User | null>(null); // New state for selected user

    useEffect(() => {
        const getUserData = async () => {
            const response = await fetch("http://localhost:8080/getUserData", {
                method: "GET",
                credentials: "include",
            });

            if (response.ok) {
                const userData: User = await response.json();

                // Process avatar path with regex
                const regex = /\/uploads\/.*/;
                const paths = userData.avatarPath?.match(regex);
                const avatarUrl = paths ? paths[0] : undefined;
                userData.avatarPath = avatarUrl;

                setUser(userData);
            } else {
                console.log("Failed to retrieve user data");
            }
        };

        getUserData();
    }, []);

    return (
        <UserContext.Provider value={{ user, setUser, selectedUser, setSelectedUser }}>
            {children}
        </UserContext.Provider>
    );
};
