"use client";

import React from "react";
import { useUser } from "@/app/context/UserContext"; // Import the custom hook to access user context
import NavBar from "@/app/components/NavBar";
import Image from 'next/image';

const Profile = () => {
  const { user } = useUser(); // Access user from context

  // If user data is not yet loaded, display a loading message
  // if (!user) {
  //   return <div>Loading...</div>;
  // }

  return (
    <>
      <NavBar logged={true} />
      {user && (
        <div className="flex flex-col bg-base-300 rounded-box w-fit">
          <div className="topWrapper flex flex-row">
            <div className="basis-2/3 flex flex-col p-8">
              {user.ID && (
                <div className="basis-1/4">{user.ID}</div>
              )}
              {user.nickname && (
                <div className="basis-1/4">{user.nickname}</div>
              )}
              <div className="basis-1/4">
                {user.firstName} {user.lastName}
              </div>
              <div className="basis-1/4">{user.email}</div>
              <div className="basis-1/4">{user.dob}</div>
            </div>
            <div className="wrapper">
              <div className="basis-1/3 avatar p-8">
                <div className="w-24 rounded-full ">
                  <Image
                    src={user.avatarPath || "/default_avatar.jpg"}
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
                    defaultChecked={user.profileVisibility === "private"}
                    className="checkbox"
                  />
                </label>
              </div>
            </div>
          </div>
          {user.aboutMe != "" && (
            <div className="botWrapper p-8">
              <div>{user.aboutMe}</div>
            </div>
          )}
        </div>
      )}
    </>
  );
};

export default Profile;
