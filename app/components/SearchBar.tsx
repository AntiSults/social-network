"use client";
import Image from 'next/image';
import React, { useState } from "react";
import { useRouter } from "next/navigation";
import { useUser } from "../context/UserContext";

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
    const { setSelectedUser } = useUser();

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

    const goToUserProfile = (user: SearchResult) => {
        setSelectedUser(user);
        router.push(`/users/${user.ID}`);
    };

    return (
        <div className="relative w-full max-w-md mx-auto mt-8">
            <div className="flex items-center bg-white shadow-md rounded-full p-2">
                <input
                    type="text"
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    placeholder="Search for users..."
                    className="w-full px-4 py-2 text-sm text-gray-700 rounded-full focus:outline-none"
                />
                <button
                    onClick={handleSearch}
                    className="bg-gray-600 text-white px-4 py-2 rounded-full hover:bg-gray-700"
                >
                    Search
                </button>
            </div>

            {searchResults.length > 0 && (
                <ul className="absolute w-full mt-2 bg-white shadow-lg rounded-lg max-h-60 overflow-y-auto z-10">
                    {searchResults.map((result: SearchResult) => (
                        <li
                            key={result.ID}
                            onClick={() => goToUserProfile(result)}
                            className="flex items-center px-4 py-2 cursor-pointer hover:bg-gray-100 transition"
                        >
                            <Image
                                src={result.avatarPath || "/default_avatar.jpg"}
                                alt={`${result.firstName}'s avatar`}
                                className="w-10 h-10 rounded-full mr-3"
                                width={250}
                                height={250}
                            />
                            <div>
                                <p className="font-semibold">
                                    {result.firstName} {result.lastName}
                                </p>
                                <p className="text-sm text-gray-500">{result.email}</p>
                            </div>
                        </li>
                    ))}
                </ul>
            )}
        </div>
    );
};

export default SearchBar;