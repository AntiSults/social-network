"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useUser } from "@/app/context/UserContext"; // Import the custom hook to access user context

const TestLoggingIn = () => {
  const router = useRouter();
  const { user } = useUser(); // Access user from context

  useEffect(() => {
    // If user is not available, redirect to the login page
    if (!user) {
      router.push("/login");
    }
  }, [user, router]);

  // Display loading state until the user context is loaded
  if (!user) {
    return <div>Loading...</div>;
  }

  return <div>Protected Data</div>;
};

export default TestLoggingIn;

