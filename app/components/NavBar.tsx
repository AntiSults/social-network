import React, { useEffect, useState } from "react";

interface Props {
  logged: boolean;
  logpage?: boolean;
}

const NavBar = ({ logged, logpage = false }: Props) => {
  const [avatarPath, setAvatarPath] = useState("");

  useEffect(() => {
    const fetchAvatar = async () => {
      const response = await fetch("http://localhost:8080/getAvatarPath", {
        method: "GET",
        credentials: "include",
      });
      if (response.ok) {
        const data = await response.json();
        const regex = /\/uploads\/.*/;
        const paths = data.avatarPath.match(regex);
        const avatarUrl = paths ? paths[0] : null;

        if (avatarUrl) {
          const img = new Image();
          img.onload = () => setAvatarPath(avatarUrl);
          img.onerror = () => setAvatarPath("/default_avatar.jpg");
          img.src = avatarUrl;
        } else {
          setAvatarPath("/default_avatar.jpg");
        }
      } else {
        console.log("failed to fetch avatar");
      }
    };
    if (logged) {
      fetchAvatar();
    }
  }, [logged]);

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
        </div>
        <div className="flex flex-1 justify-end px-2">
          <div className="flex items-stretch">
            {!logged && !logpage && (
              <a href="/login" className="btn btn-ghost rounded-btn">
                Login
              </a>
            )}
            <div className="dropdown dropdown-end">
              <div tabIndex={0} role="button">
                {logged && (
                  <div className="avatar">
                    <div className="w-12 rounded-full">
                      <img src={avatarPath} />
                    </div>
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
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default NavBar;
