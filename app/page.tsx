import React from "react";
import Link from "next/link";

const Home = () => {
  return (
    <div>
      Home
      <Link href="/register">Register</Link>
      <Link href="/login">Login</Link>
    </div>
  );
};

export default Home;
