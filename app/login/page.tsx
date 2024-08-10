"use client";

import React, { useState } from "react";
import FieldInput from "../components/FieldInput";
import { useRouter } from "next/navigation";
import Button from "../components/Button";

const LoginPage = () => {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    const formData = new FormData();
    formData.append("email", email);
    formData.append("password", password);

    try {
      const response = await fetch("http://localhost:8080/login", {
        method: "POST",
        body: formData,
      });

      if (response.ok) {
        router.push("/");
      } else {
        console.error("Form submission error:", await response.text());
      }
    } catch (error) {
      console.error("Network error:", error);
    }
  };

  return (
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
  );
};

export default LoginPage;
