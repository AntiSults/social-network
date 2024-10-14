import React, { useEffect, useState } from 'react';
import { User } from '@/app/utils/types';
import Image from 'next/image';
import Followers from '@/app/components/Followers';
import FollowList from './FollowLists';


interface Props {
    foundUser: User;  // User being viewed
    currentUser: User | null;  // Logged-in user
}

const FoundUserProfile: React.FC<Props> = ({ foundUser, currentUser }) => {
    const [isFollower, setIsFollower] = useState<boolean>(false);

    const [canViewFullProfile, setCanViewFullProfile] = useState<boolean>(false);

    useEffect(() => {

        const checkFollowerStatus = async () => {
            // Fetch follower status
            if (currentUser && foundUser) {
                fetch(`http://localhost:8080/followers/status?userId=${foundUser.ID}&followerId=${currentUser.ID}`)
                    .then((res) => res.json())
                    .then((data) => {
                        setIsFollower(data.isFollowing);
                    });
            }
        };
        if (currentUser && foundUser) {
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

        <div className="bg-white shadow-md rounded-lg p-8 max-w-lg w-full text-center">
            <h1 className="text-2xl font-bold mb-4">
                {`${foundUser.firstName} ${foundUser.lastName}'s Profile`}
            </h1>
            {canViewFullProfile ? (
                <div className="flex flex-col items-center">
                    <Image
                        src={foundUser.avatarPath || "/default_avatar.jpg"}
                        alt={`${foundUser.firstName}'s Avatar`}
                        width={250}
                        height={250}
                        className="rounded-full shadow-lg"
                    />
                    <p className="text-gray-600 mt-4">
                        Berth date: {foundUser.dob || "No details provided"}
                    </p>
                    <p className="text-gray-600 mt-4">
                        Nick: {foundUser.nickname || "No details provided"}
                    </p>
                    <p className="text-gray-600 mt-4">
                        About me: {foundUser.aboutMe || "No details provided"}
                    </p>
                    <p className="text-gray-600 mt-4">
                        Profile: {foundUser.profileVisibility || "No details provided"}
                    </p>
                    <p className="text-gray-600 mt-4">
                        Email: {foundUser.email || "No details provided"}
                    </p>
                    {<FollowList user={foundUser} />}
                </div>

            ) : (
                <p>This profile is private. You can only see limited information.</p>
            )}
        </div>

    );
};

export default FoundUserProfile;
