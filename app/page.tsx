"use client";

import React, { useEffect, useState } from "react";
import Link from "next/link";
import NavBar from "./components/NavBar";
import checkLoginStatus from "./utils/checkLoginStatus";

const Home = () => {
  const [isLoggedIn, setIsLoggedIn] = useState(false);

  useEffect(() => {
    setIsLoggedIn(checkLoginStatus());
  });

  return (
    <>
      {NavBar(isLoggedIn)}
      <div>
        <div>
          <Link href="/testLoggedIn">Test if user is logged in [click]</Link>
        </div>
      </div>
    </>
  );
};

export default Home;
