"use client";

import React, { createContext, useContext, useEffect, useState } from "react";
import { User } from "@/app/utils/types";

interface Props {
    user: User | null;
    selectedUser: User | null;
    setUser: React.Dispatch<React.SetStateAction<User | null>>;
    setSelectedUser: React.Dispatch<React.SetStateAction<User | null>>;
    updateUser: (updatedUser: Partial<User>) => void;  // Function to update user
}
// Create the context
const UserContext = createContext<Props | undefined>(undefined);

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
    const [selectedUser, setSelectedUser] = useState<User | null>(null);

    useEffect(() => {
        const getUserData = async () => {
            try {
                const response = await fetch("http://localhost:8080/getUserData", {
                    method: "GET",
                    credentials: "include",
                });

                if (response.status === 401) {
                    console.log("User is not logged in.");
                    return;
                }

                if (response.ok) {
                    const userData: User = await response.json();
                    setUser(userData);
                } else {
                    console.log("Failed to retrieve user data. Status:", response.status);
                }
            } catch (error) {
                console.error("Error retrieving user data:", error);
            }
        };
        getUserData();
    }, []);

    // Function to update user details
    const updateUser = (updatedUser: Partial<User>) => {
        if (user) {
            setUser({ ...user, ...updatedUser }); // Merge the updated fields with the current user
        }
    };
    return (
        <UserContext.Provider value={{ user, setUser, selectedUser, setSelectedUser, updateUser }}>
            {children}
        </UserContext.Provider>
    );
};
