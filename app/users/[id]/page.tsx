"use client";

import Image from 'next/image';
import React from 'react';
import { useUser } from '@/app/context/UserContext';
import NavBar from '@/app/components/NavBar';
import SearchBar from '@/app/components/SearchBar';
import Followers from '@/app/components/Followers';
import FollowList from '@/app/components/FollowLists';
import PendingRequests from '@/app/components/PendingFollowRequests';
import PendingGroupRequests from '@/app/components/GroupRequestsPending';
import PendingGroupInvites from '@/app/components/GroupInvitePending';

const ProfilePage = () => {
    const { user } = useUser();  // Focus only on the logged-in user

    if (!user) {
        return (
            <div className="min-h-screen">
                <NavBar logged={false} />
                <p className="text-center text-gray-600">Please login to see your Profile!</p>
            </div>
        );
    }

    return (
        <div className="min-h-screen">
            <NavBar logged={true} />
            <div className="flex flex-col items-center mt-10">
                <div className="w-full max-w-md mb-10">
                    <SearchBar />
                </div>
                <div className="bg-white shadow-md rounded-lg p-8 max-w-lg w-full text-center">
                    <h1 className="text-2xl font-bold mb-4">
                        {`${user.firstName} ${user.lastName}'s Profile`}
                    </h1>
                    <div className="flex flex-col items-center">
                        <Image
                            src={user.avatarPath || "/default_avatar.jpg"}
                            alt={`${user.firstName}'s Avatar`}
                            width={250}
                            height={250}
                            className="rounded-full shadow-lg"
                        />
                        <p className="text-gray-600 mt-4">
                            About Me: {user.aboutMe || "No details provided"}
                        </p>
                    </div>
                </div>
            </div>
            {/* Show pending requests only for the logged-in user */}
            <PendingGroupInvites />
            <PendingGroupRequests />
            <PendingRequests user={user} />
            <FollowList user={user} />
        </div>
    );
};

export default ProfilePage;
