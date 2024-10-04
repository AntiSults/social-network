"use client";
import React, { useEffect, useState } from "react";
import { useUser } from "@/app/context/UserContext"; // Importing useUser from UserContext
import NavBar from "./components/NavBar";
import checkLoginStatus from "./utils/checkLoginStatus";
import Image from "next/image";
import Link from "next/link";

const Home = () => {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const { user: currentUser } = useUser(); // Accessing the current user from UserContext

  useEffect(() => {
    setIsLoggedIn(checkLoginStatus());
  }, []);

  return (
    <>
      <NavBar logged={isLoggedIn} />
      <div className="flex flex-col items-center justify-center h-screen text-center">
        <h1 className="text-4xl font-bold">HI!</h1>
        <p className="text-2xl mt-4">
          This is a learning project: Social-Network of coding school
        </p>
        <div className="mt-4">
          <Image
            src="/image/kj.png"
            alt="Kood/Johvi Logo"
            width={100}
            height={100}
            className="mx-auto"
          />
        </div>
        <p className="text-2xl mt-4">Kood/Johvi</p>

        {currentUser ? (
          <p className="text-2xl mt-4">
            {currentUser?.firstName} {currentUser?.lastName} is logged in
          </p>
        ) : (
          <p className="text-2xl mt-4">
            Please <Link href="/login" className="text-blue-500 underline">login/register</Link> to continue
          </p>

        )}
      </div >
    </>
  );
};

export default Home;


