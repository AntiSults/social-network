import React from "react";

interface Props {
  type?: "button" | "submit" | "reset";
  name: string;
  className?: string; // Accept className as a prop
}

const Button = ({ type = "button", name, className = "" }: Props) => {
  return (
    <button type={type} className={`btn ${className}`}> {/* Merge className with default */}
      {name}
    </button>
  );
};

export default Button;

