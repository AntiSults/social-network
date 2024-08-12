"use client";

import React, { useEffect, useState } from "react";
import FieldInput from "../components/FieldInput";
import { useRouter } from "next/navigation";
import Button from "../components/Button";

const LoginPage = () => {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  //GOTTA USE USEEFFECT TO RECOGNISE ERRORS FROM MIDDLEWARE
  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    const errorParam = params.get("error");
    if (errorParam) {
      setError(errorParam);
    }
  }, []);

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    const formData = new FormData();
    formData.append("email", email);
    formData.append("password", password);

    try {
      const response = await fetch("http://localhost:8080/login", {
        method: "POST",
        credentials: "include",
        body: formData,
      });

      if (response.ok) {
        router.push("/");
      } else {
        // Handle errors by reading the error message from the response
        const data = await response.json();
        console.log(data.message);
        setError(data.message || "Login failed");
      }
    } catch (error) {
      console.error("Network error:", error);
    }
  };

  return (
    <>
      <div>{error && <p style={{ color: "red" }}>{error}</p>}</div>
      <form onSubmit={handleLogin}>
        <FieldInput
          name="Email"
          type="text"
          placeholder="example@example.com"
          required={true}
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        />
        <FieldInput
          name="Password"
          type="password"
          placeholder="example"
          required={true}
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        <Button type="submit" name="Login" />
      </form>
    </>
  );
};

export default LoginPage;
