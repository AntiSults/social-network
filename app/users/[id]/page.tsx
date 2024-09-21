"use client";

import React from "react";
import { useUser } from "../../context/UserContext";
import NavBar from "../../components/NavBar";
import SearchBar from "../../components/SearchBar";

const ProfilePage = () => {
    const { user } = useUser();
    if (!user) {
        return <p>Loading...</p>;
    }
    return (
        <div>
            <NavBar logged={true} />
            <h1>{`${user.firstName} ${user.lastName}'s Profile`}</h1>
            <div>
                <img
                    src={user.avatarPath || "/default_avatar.jpg"}
                    alt={`${user.firstName}'s Avatar`}
                />
                <p>About Me: {user.aboutMe || "No details provided"}</p>
            </div>
            <SearchBar />
        </div>
    );
};

export default ProfilePage;
