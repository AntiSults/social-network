"use client";

import React from "react";
import { useUser } from "../../context/UserContext";
import NavBar from "../../components/NavBar";
import SearchBar from "../../components/SearchBar";

const ProfilePage = () => {
    const { user, selectedUser } = useUser();

    const profileUser = selectedUser || user;

    if (!profileUser) {
        return <p>Loading...</p>;
    }

    return (
        <div>
            <NavBar logged={true} />
            <h1>{`${profileUser.firstName} ${profileUser.lastName}'s Profile`}</h1>
            <div>
                <img
                    src={profileUser.avatarPath || "/default_avatar.jpg"}
                    alt={`${profileUser.firstName}'s Avatar`}
                />
                <p>About Me: {profileUser.aboutMe || "No details provided"}</p>
            </div>
            <SearchBar />
        </div>
    );
};

export default ProfilePage;




