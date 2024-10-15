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
import FoundUserProfile from '@/app/components/FoundUserProfile';

const ProfilePage = () => {
    const { user, selectedUser } = useUser();
    const profileUser = selectedUser || user;

    if (!profileUser) {
        return (
            <div className="min-h-screen">
                <NavBar logged={false} />
                <p className="text-center text-gray-600">Please login to see User Profile!</p>
            </div>
        );
    }
    const isCurrentUser = profileUser.ID === user?.ID;

    return (
        <div className="min-h-screen">
            <NavBar logged={true} />
            <div className="flex flex-col items-center mt-10">
                <div className="w-full max-w-md mb-10">
                    <SearchBar />
                </div>

                {/* If the profileUser is the current user, show their profile details, otherwise show FoundUserProfile */}
                {isCurrentUser ? (
                    <div className="bg-white shadow-md rounded-lg p-8 max-w-lg w-full text-center">
                        <h1 className="text-2xl font-bold mb-4">
                            {`${profileUser.firstName} ${profileUser.lastName}'s Profile`}
                        </h1>
                        <div className="flex flex-col items-center">
                            <Image
                                src={profileUser.avatarPath ? `${profileUser.avatarPath}` : "/default_avatar.jpg"}
                                alt={`${profileUser.firstName}'s Avatar`}
                                width={250}
                                height={250}
                                className="rounded-full shadow-lg"
                            />
                            <p className="text-gray-600 mt-4">
                                About Me: {profileUser.aboutMe || "No details provided"}
                            </p>
                        </div>
                    </div>
                ) : (
                    // Render FoundUserProfile for the found user
                    <FoundUserProfile foundUser={profileUser} currentUser={user} />
                )}
            </div>
            {/* Render other components based on the current user's ID */}
            {isCurrentUser && <PendingGroupInvites />}
            {isCurrentUser && <PendingGroupRequests />}
            {isCurrentUser && <PendingRequests user={user} />}
            {!isCurrentUser && <Followers profileUser={profileUser} user={user} />}
            {isCurrentUser && <FollowList user={user} />}
        </div>
    );
};

export default ProfilePage;
