import Image from 'next/image';
import React, { useEffect, useState } from "react";
import { useUser } from "@/app/context/UserContext";
import { useNotificationWS } from "@/app/hooks/UseNotify";


interface Props {
  logged: boolean;
  logpage?: boolean;
}

interface User {
  ID: number;
  firstName: string;
  lastName: string;
}

const NavBar = ({ logged, logpage = false }: Props) => {
  const { user } = useUser();

  const [notifications, setNotifications] = useState<User>();
  useNotificationWS(setNotifications);

  const handleLogOut = async () => {
    const response = await fetch("http://localhost:8080/logout", {
      method: "POST",
      credentials: "include",
    });
    if (response.ok) {
      window.location.href = "/";
    } else {
      const data = await response.json();
      console.log("Logout failed: ", data.message);
    }
  };

  //this one logging out if you close the tab
  useEffect(() => {
    let isNavigating = false;
    // Detect navigation within the app
    const handleBeforeNavigate = () => {
      isNavigating = true;
    };
    const handleBeforeUnload = async (event: BeforeUnloadEvent) => {
      if (!isNavigating) {
        // Only log out if it is not navigating but closing the tab
        await fetch("http://localhost:8080/logout", {
          method: "POST",
          credentials: "include",
        });
      }
    };
    window.addEventListener("beforeunload", handleBeforeUnload);
    window.addEventListener("click", handleBeforeNavigate);  // Tracks clicks on links/buttons for navigation

    return () => {
      window.removeEventListener("beforeunload", handleBeforeUnload);
      window.removeEventListener("click", handleBeforeNavigate);
    };
  }, []);
  useEffect(() => {
    console.log("Updated notifications:", notifications);
  }, [notifications]);

  return (
    <div>
      <div className="navbar bg-base-300 rounded-box">
        <div className="flex-1 px-2 lg:flex-none">
          <a href="/" className="text-lg font-bold">
            Home
          </a>
          <a href={`/users/${user?.ID}`} className="ml-4 text-lg font-bold">
            User Page
          </a>
          <a href={!logged && !user ? "/posts" : `/users/${user?.ID}/posts`} className="ml-4 text-lg font-bold">
            Posts
          </a>
          <a href={`/users/${user?.ID}/chat`} className="ml-4 text-lg font-bold">
            Chat
          </a>
          <a href={`/users/${user?.ID}/groups`} className="ml-4 text-lg font-bold">
            Groups
          </a>
          <a href={`/users/${user?.ID}/events`} className="ml-4 text-lg font-bold">
            Group-Events
          </a>
        </div>
        <div className="flex flex-1 justify-end px-2">
          <div className="flex items-stretch">
            {!logged && !logpage && (
              <a href="/login" className="btn btn-ghost rounded-btn">
                Login
              </a>
            )}
            {logged && user && (
              <div className="dropdown dropdown-end">
                <div tabIndex={0} role="button">
                  <div className="avatar">
                    <div className="w-14 rounded-full">
                      <Image
                        src={user.avatarPath ? user.avatarPath : "/default_avatar.jpg"}
                        alt="User Avatar"
                        width={150}
                        height={150}
                      />
                    </div>
                  </div>
                  {/* Flashing bulb or notification icon */}
                  {notifications && (
                    <div className="ml-4 relative">
                      <span className="animate-ping absolute inline-flex h-4 w-4 rounded-full bg-red-400 opacity-75"></span>
                      <span className="relative inline-flex rounded-full h-4 w-4 bg-red-500"></span>
                    </div>
                  )}
                </div>
                <ul
                  tabIndex={0}
                  className="menu dropdown-content bg-base-100 rounded-box z-[1] mt-4 w-52 p-2 shadow"
                >
                  <li>
                    <a href="/profile">My Profile</a>
                  </li>
                  <li>
                    <a onClick={handleLogOut}>Log out</a>
                  </li>
                  {/* Notifications Dropdown */}
                  {notifications && (
                    <>
                      <hr />
                      <li className="font-bold">Pending Requests</li>

                      <li key={notifications.ID}>
                        <a>
                          {notifications.firstName} {notifications.lastName}
                        </a>
                      </li>

                    </>
                  )}
                </ul>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default NavBar;
