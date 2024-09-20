"use client";

import React from "react";
import { useUser } from "../../context/UserContext";
import NavBar from "../../components/NavBar";

const ProfilePage = () => {
    const { user } = useUser();

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
        </div>
    );
};

export default ProfilePage;