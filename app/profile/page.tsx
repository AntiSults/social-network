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
        setUser(userData);
      } else {
        console.log("Failed to retrieve user data");
      }
    };
    getUserData();
  }, []);

  return <div>test</div>;
};

export default Profile;
