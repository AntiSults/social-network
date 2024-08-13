"use client";
//TEST FUNCTION DELETE LATER
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";

const TestLoggingIn = () => {
  const router = useRouter();
  const [isLoggedIn, setIsLoggedIn] = useState(false);

  useEffect(() => {
    const checkAuth = async () => {
      const response = await fetch("http://localhost:8080/testLoggedIn", {
        method: "GET",
        credentials: "include", // Ensure cookies are sent with the request
      });

      if (!response.ok) {
        router.push("/login");
      } else {
        setIsLoggedIn(true);
      }
    };
    checkAuth();
  }, [router]);

  if (!isLoggedIn) {
    return <div>Loading...</div>;
  } else {
    return <div>Protected Data</div>;
  }
};

export default TestLoggingIn;
