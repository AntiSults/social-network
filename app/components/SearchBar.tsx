// app/components/SearchBar.tsx
"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import { useUser } from "../context/UserContext"; // Import useUser

interface SearchResult {
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

const SearchBar: React.FC = () => {
    const [searchQuery, setSearchQuery] = useState("");
    const [searchResults, setSearchResults] = useState<SearchResult[]>([]);
    const router = useRouter();
    const { setSelectedUser } = useUser(); // Access setSelectedUser

    const handleSearch = async () => {
        try {
            const response = await fetch(`http://localhost:8080/search?query=${encodeURIComponent(searchQuery)}`);
            if (response.ok) {
                const users = await response.json();
                setSearchResults(users);
            } else {
                console.log("Search failed");
            }
        } catch (error) {
            console.error("Error during search:", error);
        }
    };

    // Navigate to the selected user's profile and set the selected user in context
    const goToUserProfile = (user: SearchResult) => {
        setSelectedUser(user); // Set the selected user
        router.push(`/users/${user.ID}`); // Navigate to user profile
    };

    return (
        <div>
            <input
                type="text"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder="Search for other users..."
            />
            <button onClick={handleSearch}>Search</button>

            {searchResults.length > 0 && (
                <ul>
                    {searchResults.map((result: SearchResult) => (
                        <li key={result.ID} onClick={() => goToUserProfile(result)}>
                            {result.firstName} {result.lastName} ({result.email})
                        </li>
                    ))}
                </ul>
            )}
        </div>
    );
};

export default SearchBar;
