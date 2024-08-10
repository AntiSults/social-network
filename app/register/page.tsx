"use client";

import React, { useState } from "react";
import styles from "./page.module.css";
import FieldInput from "../components/FieldInput";
import Button from "../components/Button";

import { useRouter } from "next/navigation";
import AvatarUploadField from "../components/AvatarUploadField";

const RegisterPage = () => {
  // Using useStates to be able to change the values later, also
  // the reason why all of this code has to be in this function.
  // Couldn't get the useStates working outside of it
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [dob, setDob] = useState("");
  const [nickname, setNickname] = useState("");
  const [aboutMe, setAboutMe] = useState("");
  const [avatar, setAvatar] = useState<File | null>(null);
  const router = useRouter();

  // Get called when button is clicked, sends data to backend
  const HandleRegisterForm = async (e: React.FormEvent) => {
    e.preventDefault();

    const formData = new FormData();
    formData.append("email", email);
    formData.append("password", password);
    formData.append("firstName", firstName);
    formData.append("lastName", lastName);
    formData.append("dob", dob);
    if (avatar) {
      formData.append("avatar", avatar);
    }
    if (nickname != "") {
      formData.append("nickname", nickname);
    }
    if (aboutMe != "") {
      formData.append("aboutMe", aboutMe);
    }

    try {
      const response = await fetch("http://localhost:8080/register", {
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

  // Gets called when an avatar is uploaded
  const HandleFileSelect = (file: File) => {
    setAvatar(file);
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
        <FieldInput
          name="Nickname[Optional]"
          type="text"
          placeholder="Johnny"
          required={false}
          value={nickname}
          onChange={(e) => setNickname(e.target.value)}
        />
        <FieldInput
          name="About me[Optional]"
          type="textarea"
          placeholder="I'm a cool guy!"
          required={false}
          value={aboutMe}
          onChange={(e) => setAboutMe(e.target.value)}
        />
        <AvatarUploadField onFileSelect={HandleFileSelect} />
        <Button type="submit" name="Register" />
      </div>
    </form>
  );
};

export default RegisterPage;
