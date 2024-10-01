"use client";
import React, { useEffect, useState } from "react";
import { User } from "@/app/utils/types";

interface FollowListProps {
    user: User | null;
}

const FollowList: React.FC<FollowListProps> = ({ user }) => {
    const [following, setFollowing] = useState<User[]>([]);
    const [followers, setFollowers] = useState<User[]>([]);

    useEffect(() => {
        if (user) {
            fetch(`http://localhost:8080/followers/followersList?userId=${user.ID}`)
                .then((res) => res.json())
                .then((data: { followers: User[], following: User[] }) => {
                    setFollowers(data.followers);
                    setFollowing(data.following);
                })
                .catch((err) => console.error("Error fetching followers and following:", err));
        }
    }, [user]);

    return (
        <div className="flex flex-col items-center mt-6">
            <div className="w-full max-w-md mb-10">
                <h2 className="text-xl font-bold mb-4">Following</h2>
                {following.length > 0 ? (
                    following.map((follower) => (
                        <div key={follower.ID} className="bg-white p-4 rounded-lg shadow-md mb-2">
                            {follower.firstName} {follower.lastName}
                        </div>
                    ))
                ) : (
                    <p className="text-gray-500 text-center">You are not following anyone yet.</p>
                )}
            </div>
            <div className="w-full max-w-md">
                <h2 className="text-xl font-bold mb-4">Followers</h2>
                {followers.length > 0 ? (
                    followers.map((follower) => (
                        <div key={follower.ID} className="bg-white p-4 rounded-lg shadow-md mb-2">
                            {follower.firstName} {follower.lastName}
                        </div>
                    ))
                ) : (
                    <p className="text-gray-500 text-center">You have no followers yet.</p>
                )}
            </div>
        </div>
    );
};

export default FollowList;
