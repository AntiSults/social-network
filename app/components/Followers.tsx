"use client";
import React, { useState, useEffect } from "react";
import { User } from "../utils/types"

interface FollowersProps {
    profileUser: User | null;
    user: User | null;
}

const Followers: React.FC<FollowersProps> = ({ profileUser, user }) => {
    const [isFollowing, setIsFollowing] = useState<boolean>(false);
    const [isPending, setIsPending] = useState<boolean>(false);

    useEffect(() => {
        if (profileUser && user) {
            // Fetch follow status from backend
            fetch(`http://localhost:8080/followers/status?userId=${profileUser.ID}&followerId=${user.ID}`)
                .then((res) => res.json())
                .then((data) => {
                    setIsFollowing(data.isFollowing);
                    setIsPending(data.isPending);
                });
        }
    }, [profileUser, user]);

    if (!profileUser || !user) {
        return null; // Don't render anything if profileUser or user is null
    }

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

    return (
        <div>
            {isFollowing ? (
                <button onClick={handleUnfollow}>Unfollow</button>
            ) : isPending ? (
                <p>Follow request pending...</p>
            ) : (
                <button onClick={handleFollow}>Follow</button>
            )}
        </div>
    );
};

export default Followers;
