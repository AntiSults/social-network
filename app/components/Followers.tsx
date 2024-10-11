import React, { useState, useEffect } from "react";
import { User } from "@/app/utils/types"
interface Props {
    profileUser: User | null;
    user: User | null;
}
const Followers: React.FC<Props> = ({ profileUser, user }) => {
    const [isFollowing, setIsFollowing] = useState<boolean>(false);
    const [isPending, setIsPending] = useState<boolean>(false);
    useEffect(() => {
        if (profileUser && user) {
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
        <div className="flex justify-center mt-6">
            {isFollowing ? (
                <button
                    onClick={handleUnfollow}
                    className="bg-red-500 hover:bg-red-600 text-white font-bold py-2 px-6 rounded-lg shadow-md transition duration-200 ease-in-out"

                >
                    Unfollow
                </button>
            ) : isPending ? (
                <p className="text-gray-500 text-center">Follow request pending...</p>
            ) : (
                <button
                    onClick={handleFollow}
                    className="bg-gray-600 hover:bg-gray-700 text-white font-bold py-2 px-6 rounded-lg shadow-md transition duration-200 ease-in-out"
                >
                    Follow
                </button>
            )}
        </div>
    );
};
export default Followers;
