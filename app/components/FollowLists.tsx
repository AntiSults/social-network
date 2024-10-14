import React, { useEffect, useState } from 'react';
import { User } from '@/app/utils/types';

interface Props {
    user: User | null;
}

const FollowList: React.FC<Props> = ({ user }) => {
    const [following, setFollowing] = useState<User[]>([]);
    const [followers, setFollowers] = useState<User[]>([]);
    const [error, setError] = useState<string | null>(null);
    const [isLoading, setIsLoading] = useState<boolean>(true);

    useEffect(() => {
        const fetchFollowLists = async () => {
            if (!user) {
                setError("User not logged in");
                setIsLoading(false);
                return;
            }

            try {
                const response = await fetch(`http://localhost:8080/followers/followersList?userId=${user.ID}`);

                if (!response.ok) {
                    throw new Error("Failed to fetch follow lists");
                }

                const data: { followers: User[], following: User[] } = await response.json();
                if (!data || (!Array.isArray(data.followers) && data.followers !== null) || (!Array.isArray(data.following) && data.following !== null)) {
                    throw new Error("Invalid data received");
                }

                setFollowers(data.followers || []);
                setFollowing(data.following || []);

                setError(null);
            } catch (err) {
                setError((err as Error).message || "An unexpected error occurred");
            } finally {
                setIsLoading(false);
            }
        };

        fetchFollowLists();
    }, [user]);

    if (isLoading) {
        return <p>Loading...</p>;
    }

    return (
        <div className="flex flex-col items-center mt-6">
            {error && <p className="text-red-500">{error}</p>}
            <div className="w-full max-w-md mb-10">
                <h2 className="text-xl font-bold mb-4">Following</h2>
                {following.length > 0 ? (
                    following.map((user) => (
                        <div key={user.ID} className="bg-white p-4 rounded-lg shadow-md mb-2">
                            {user.firstName} {user.lastName}
                        </div>
                    ))
                ) : (
                    <p className="text-gray-500 text-center">Is not following anyone yet.</p>
                )}
            </div>
            <div className="w-full max-w-md">
                <h2 className="text-xl font-bold mb-4">Followers</h2>
                {followers.length > 0 ? (
                    followers.map((user) => (
                        <div key={user.ID} className="bg-white p-4 rounded-lg shadow-md mb-2">
                            {user.firstName} {user.lastName}
                        </div>
                    ))
                ) : (
                    <p className="text-gray-500 text-center">No followers yet.</p>
                )}
            </div>
        </div>
    );
};

export default FollowList;
