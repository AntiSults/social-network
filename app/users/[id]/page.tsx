"use client";

import React, { useState } from "react";
import { useUser } from "../../context/UserContext";
import NavBar from "../../components/NavBar";

const ProfilePage = () => {
    const { user } = useUser();

    const [searchQuery, setSearchQuery] = useState("");
    const [searchResults, setSearchResults] = useState([]);

    const handleSearch = async () => {
        const response = await fetch(`http://localhost:8080/search?query=${searchQuery}`);
        if (response.ok) {
            const users = await response.json();
            setSearchResults(users);
        } else {
            console.log("Search failed");
        }
    };

    if (!user) {
        return <p>Loading...</p>;
    }

    return (
        <div>
            <NavBar logged={true} /> {/* Add the NavBar */}
            <h1>{`${user.firstName} ${user.lastName}'s Profile`}</h1>
            <div>
                <img
                    src={user.avatarPath || "/default_avatar.jpg"}
                    alt={`${user.firstName}'s Avatar`}
                />
                <p>About Me: {user.aboutMe || "No details provided"}</p>
            </div>
            {/* Search Bar */}
            <input
                type="text"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder="Search for other users..."
            />
            <button onClick={handleSearch}>Search</button>

            {/* Search Results */}
            {searchResults.length > 0 && (
                <ul>
                    {searchResults.map((result) => (
                        <li key={result.id}>
                            {result.firstName} {result.lastName} ({result.email})
                        </li>
                    ))}
                </ul>
            )}
        </div>
    );
};

export default ProfilePage;