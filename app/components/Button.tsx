import React from "react";

interface Props {
  type?: "button" | "submit" | "reset";
  name: string;
}

const Button = ({ type, name }: Props) => {
  return (
    <button type={type} className="btn">
      {name}
    </button>
  );
};

export default Button;
