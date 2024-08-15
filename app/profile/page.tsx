"use client";

import React, { useEffect, useState } from "react";

const Profile = () => {
  const [user, setUser] = useState(Object);
  useEffect(() => {
    const getUserData = async () => {
      const response = await fetch("http://localhost:8080/getUserData", {
        method: "GET",
        credentials: "include",
      });
      if (response.ok) {
        const userData = await response.json();
        console.log(userData);
        const regex = /\/uploads\/.*/;
        const paths = userData.avatarPath.match(regex);
        const avatarUrl = paths ? paths[0] : null;
        userData.avatarPath = avatarUrl;
        console.log(userData.profileVisibility);
        setUser(userData);
      } else {
        console.log("Failed to retrieve user data");
      }
    };
    getUserData();
  }, []);

  return (
    <>
      {user && (
        <div className="flex flex-col bg-base-300 rounded-box w-fit">
          <div className="topWrapper flex flex-row">
            <div className="basis-2/3 flex flex-col p-8">
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
                  <img
                    src={
                      user.avatarPath ? user.avatarPath : "/default_avatar.jpg"
                    }
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
