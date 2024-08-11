"use client";
//TEST FUNCTION DELETE LATER
import { useEffect } from "react";
import { useRouter } from "next/navigation";

const TestLoggingIn = () => {
  const router = useRouter();

  useEffect(() => {
    const checkAuth = async () => {
      const response = await fetch("http://localhost:8080/testLoggedIn", {
        method: "GET",
        credentials: "include",
      });

      if (!response.ok) {
        router.push("/login");
      }
    };

    checkAuth();
  }, [router]);

  return <div>Protected Content</div>;
};

export default TestLoggingIn;
