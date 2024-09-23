"use client";
import Image from 'next/image';
import React, { useState, useEffect } from "react";
import { useUser } from "../../context/UserContext";
import NavBar from "../../components/NavBar";
import SearchBar from "../../components/SearchBar";

const ProfilePage = () => {
    const { user, selectedUser } = useUser();
    const profileUser = selectedUser || user;

    const [isFollowing, setIsFollowing] = useState(false);
    const [isPending, setIsPending] = useState(false);

    useEffect(() => {
        if (profileUser && user) {
            // Check if already following
            fetch(`http://localhost:8080/followers/status?userId=${profileUser.ID}&followerId=${user.ID}`)
                .then((res) => res.json())
                .then((data) => {
                    setIsFollowing(data.isFollowing);
                    setIsPending(data.isPending);
                });
        }
    }, [profileUser, user]);

    const handleFollow = async () => {
        const response = await fetch("http://localhost:8080/followers", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ userId: profileUser?.ID, followerId: user?.ID }),
        });

        if (response.ok) {
            if (profileUser?.profileVisibility === "private") {
                setIsPending(true);
            } else {
                setIsFollowing(true);
            }
        }
    };

    const handleUnfollow = async () => {
        const response = await fetch(`http://localhost:8080/followers?userId=${profileUser?.ID}&followerId=${user?.ID}`, {
            method: "DELETE",
        });

        if (response.ok) {
            setIsFollowing(false);
            setIsPending(false);
        }
    };

    if (!profileUser) {
        return <p>Loading...</p>;
    }

    return (
        <div>
            <NavBar logged={true} />
            <h1>{`${profileUser.firstName} ${profileUser.lastName}'s Profile`}</h1>
            <div>
                <Image
                    src={profileUser.avatarPath || "/default_avatar.jpg"}
                    alt={`${profileUser.firstName}'s Avatar`}
                    width={250}
                    height={250}
                />
                <p>About Me: {profileUser.aboutMe || "No details provided"}</p>
            </div>

            {/* Follow / Unfollow Button */}
            {profileUser?.ID !== user?.ID && (
                <div>
                    {isFollowing ? (
                        <button onClick={handleUnfollow}>Unfollow</button>
                    ) : isPending ? (
                        <p>Follow request pending...</p>
                    ) : (
                        <button onClick={handleFollow}>Follow</button>
                    )}
                </div>
            )}

            <SearchBar />
        </div>
    );
};

export default ProfilePage;





