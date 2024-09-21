"use client";

import React, { useState } from "react";

interface SearchResult {
    id: number;
    firstName: string;
    lastName: string;
    email: string;
    nickname?: string;
    aboutMe?: string;
    avatarPath?: string;
}

const SearchBar: React.FC = () => {
    const [searchQuery, setSearchQuery] = useState("");
    const [searchResults, setSearchResults] = useState<SearchResult[]>([]);

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
                        <li key={result.id}>
                            {result.firstName} {result.lastName} ({result.email})
                        </li>
                    ))}
                </ul>
            )}
        </div>
    );
};

export default SearchBar;
