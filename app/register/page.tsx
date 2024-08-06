"use client";

import React, { useState } from "react";
import styles from "./page.module.css";
import FieldInput from "../components/FieldInput";
import Button from "../components/Button";
import { redirect } from "next/dist/server/api-utils";

import { useRouter } from "next/navigation";

const RegisterPage = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [dob, setDob] = useState("");
  const router = useRouter();

  const HandleRegisterForm = async (e: React.FormEvent) => {
    console.log("Form submitted");
    e.preventDefault();

    const formData = {
      email,
      password,
      firstName,
      lastName,
      dob,
    };

    try {
      const response = await fetch("http://localhost:8080/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(formData),
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
    <form onSubmit={HandleRegisterForm}>
      <div id={styles.register} className="mx-auto">
        <FieldInput
          name="Email"
          type="email"
          placeholder="example@example.com"
          required={true}
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        />
        <FieldInput
          name="Password"
          type="password"
          placeholder="shh secret"
          required={true}
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        <FieldInput
          name="First Name"
          type="text"
          placeholder="John"
          required={true}
          value={firstName}
          onChange={(e) => setFirstName(e.target.value)}
        />
        <FieldInput
          name="Last Name"
          type="text"
          placeholder="Smith"
          required={true}
          value={lastName}
          onChange={(e) => setLastName(e.target.value)}
        />
        <FieldInput
          name="Date of birth"
          type="date"
          placeholder=""
          required={true}
          value={dob}
          onChange={(e) => setDob(e.target.value)}
        />
        <Button type="submit" name="Register" />
      </div>
    </form>
  );
};

export default RegisterPage;
