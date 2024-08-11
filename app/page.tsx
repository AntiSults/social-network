import React from "react";
import Link from "next/link";

const Home = () => {
  return (
    <div>
      Home
      <div>
        <Link href="/register">Register</Link>
      </div>
      <div>
        <Link href="/login">Login</Link>
      </div>
      <div>
        <Link href="/testLoggedIn">Test if user is logged in </Link>
      </div>
    </div>
  );
};

export default Home;
