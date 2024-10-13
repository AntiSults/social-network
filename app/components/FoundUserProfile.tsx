import Image from 'next/image';
import React, { useState, useEffect, useCallback } from 'react';
import { User } from '@/app/utils/types';
import Followers from '@/app/components/Followers';
import FollowList from '@/app/components/FollowLists';

interface Props {
    foundUser: User;
    currentUser: User;
}

const FoundUserProfile: React.FC<Props> = ({ currentUser, foundUser }) => {
    const [isFollower, setIsFollower] = useState(false);
    const [canViewFullProfile, setCanViewFullProfile] = useState(false);

    // Check follower status and profile visibility
    useEffect(() => {
        const checkFollowerStatus = async () => {
            try {
                const response = await fetch(`http://localhost:8080/followers/check?currentUserId=${currentUser?.ID}&foundUserId=${foundUser.ID}`);
                const data = await response.json();
                setIsFollower(data.isFollower);
            } catch (error) {
                console.error('Error fetching follow status:', error);
            }
        };

        if (currentUser && currentUser.ID !== foundUser.ID) {
            checkFollowerStatus();
        }
    }, [currentUser, foundUser]);

    useEffect(() => {
        if (foundUser.profileVisibility === "public" || (foundUser.profileVisibility === "private" && isFollower)) {
            setCanViewFullProfile(true);
        } else {
            setCanViewFullProfile(false);
        }
    }, [foundUser, isFollower]);

    return (
        <div className="min-h-screen">
            <div className="flex flex-col items-center mt-10">
                <div className="bg-white shadow-md rounded-lg p-8 max-w-lg w-full text-center">
                    <h1 className="text-2xl font-bold mb-4">
                        {`${foundUser.firstName} ${foundUser.lastName}'s Profile`}
                    </h1>
                    <div className="flex flex-col items-center">
                        <Image
                            src={foundUser.avatarPath || "/default_avatar.jpg"}
                            alt={`${foundUser.firstName}'s Avatar`}
                            width={250}
                            height={250}
                            className="rounded-full shadow-lg"
                        />
                        <p className="text-gray-600 mt-4">
                            About Me: {foundUser.aboutMe || "No details provided"}
                        </p>
                    </div>
                    {canViewFullProfile ? (
                        <div className="bg-white shadow-md rounded-lg p-8 max-w-lg w-full text-center mt-4">
                            <p>Email: {foundUser.email}</p>
                            <p>Date of Birth: {foundUser.dob}</p>
                            <p>Nickname: {foundUser.nickname || "Not provided"}</p>
                            {/* You can display followers and following logic here */}
                        </div>
                    ) : (
                        <p className="text-center text-gray-600 mt-4">This profile is private.</p>
                    )}
                    {/* Follow/Unfollow button */}
                    {currentUser.ID !== foundUser.ID && (
                        <button className="mt-4 bg-blue-500 text-white p-2 rounded">
                            {isFollower ? "Unfollow" : "Follow"}
                        </button>
                    )}
                </div>
            </div>

            {/* Conditional rendering based on profile visibility */}
            {canViewFullProfile ? (
                <>
                    {/* Components for followers, following, posts, etc. */}
                    <Followers profileUser={foundUser} user={currentUser} />
                    <FollowList user={foundUser} />
                </>
            ) : (
                <div className="text-center text-gray-600 mt-4">
                    This profile is private. Follow to see more details.
                </div>
            )}
        </div>
    );
};

export default FoundUserProfile;
