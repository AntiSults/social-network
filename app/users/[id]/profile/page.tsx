"use client";

import React, { useState, useEffect } from 'react';
import { useUser } from '@/app/context/UserContext';
import NavBar from '@/app/components/NavBar';
import Image from 'next/image';

const Profile = () => {
  const { user, updateUser } = useUser();
  const [profileVisibility, setProfileVisibility] = useState("public");

  // Set the initial profile visibility state from user data
  useEffect(() => {
    if (user?.profileVisibility) {
      setProfileVisibility(user.profileVisibility);
    }
  }, [user]);

  const handleVisibilityChange = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const newVisibility = event.target.checked ? "private" : "public";
    setProfileVisibility(newVisibility);

    try {
      const response = await fetch('http://localhost:8080/user-update', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ userId: user?.ID, profileVisibility: newVisibility }),
      });

      if (response.ok) {
        updateUser({ profileVisibility: newVisibility }); // Update the user in context
      } else {
        console.error("Failed to update profile visibility");
      }
    } catch (error) {
      console.error("Error updating profile visibility:", error);
    }
  };

  return (
    <>
      <NavBar logged={true} />
      {user && (
        <div className="flex flex-col bg-base-300 rounded-box w-fit">
          <div className="topWrapper flex flex-row">
            <div className="basis-2/3 flex flex-col p-8">
              <div className="basis-1/4">UserID: {user.ID}</div>
              <div className="basis-1/4">Nick: {user.nickname}</div>
              <div className="basis-1/4">They call me: {user.firstName} {user.lastName}</div>
              <div className="basis-1/4">Email: {user.email}</div>
              <div className="basis-1/4">Made in: {
                user.dob ? new Date(user.dob).toLocaleDateString('en-US', {
                  year: 'numeric',
                  month: 'long',
                  day: 'numeric',
                }) : "No details provided"
              }</div>
            </div>
            <div className="wrapper">
              <div className="basis-1/3 avatar p-8">
                <div className="w-24 rounded-full">
                  <Image
                    src={user.avatarPath ? `${user.avatarPath}` : "/default_avatar.jpg"}
                    alt={`${user.firstName}'s Avatar`}
                    width={250}
                    height={250}
                    className="rounded-full shadow-lg"
                  />
                </div>
              </div>
              <div className="form-control p-8">
                <label className="label cursor-pointer">
                  <span className="label-text">Private</span>
                  <input
                    type="checkbox"
                    checked={profileVisibility === "private"}
                    onChange={handleVisibilityChange}
                    className="checkbox"
                  />
                </label>
              </div>
            </div>
          </div>
          {user.aboutMe !== "" && (
            <div className="botWrapper p-8">
              <div>About me: {user.aboutMe}</div>
            </div>
          )}
        </div>
      )}
    </>
  );
};

export default Profile;
