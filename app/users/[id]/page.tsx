"use client";
import Image from 'next/image';
import React from "react";
import { useUser } from "@/app/context/UserContext";
import NavBar from "@/app/components/NavBar";
import SearchBar from "@/app/components/SearchBar";
import Followers from "@/app/components/Followers";
import PendingRequests from "@/app/components/PendingRequests";
import PendingGroupRequests from "@/app/components/GroupRequestsPending";
const ProfilePage = () => {
    const { user, selectedUser } = useUser();
    const profileUser = selectedUser || user;
    if (!profileUser) {
        return <p>Loading...</p>;
    }
    return (
        <div className="min-h-screen bg-gray-50">
            <NavBar logged={true} />
            <div className="flex flex-col items-center mt-10">
                <div className="w-full max-w-md mb-10">
                    <SearchBar />
                </div>
                <div className="bg-white shadow-md rounded-lg p-8 max-w-lg w-full text-center">
                    <h1 className="text-2xl font-bold mb-4">
                        {`${profileUser.firstName} ${profileUser.lastName}'s Profile`}
                    </h1>
                    <div className="flex flex-col items-center">
                        <Image
                            src={profileUser.avatarPath || "/default_avatar.jpg"}
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
            </div>
            {/* Pending Group Join Requests (Only if viewing own profile) */}
            {profileUser?.ID === user?.ID && <PendingGroupRequests />}
            {/* Pending Follow Requests (Only if viewing own profile) */}
            {profileUser?.ID === user?.ID && <PendingRequests user={user} />}
            {/* Follow / Unfollow Component */}
            {profileUser?.ID !== user?.ID && (
                <Followers profileUser={profileUser} user={user} />
            )}
        </div>
    );
};
export default ProfilePage;
