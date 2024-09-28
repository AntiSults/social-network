import Image from 'next/image';
import React from "react";
import { useUser } from "@/app/context/UserContext";

interface Props {
  logged: boolean;
  logpage?: boolean;
}

const NavBar = ({ logged, logpage = false }: Props) => {
  const { user } = useUser();

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
