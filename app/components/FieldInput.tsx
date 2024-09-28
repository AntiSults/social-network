import React from "react";

interface Props {
  name?: string;
  type: string;
  placeholder: string;
  required: boolean;
  value: string;
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  className?: string; // Accept className as a prop
}

const FieldInput = ({
  name,
  type,
  placeholder,
  required,
  value,
  onChange,
  className = "", // Default to empty string
}: Props) => {
  return (
    <label className={`input input-bordered flex items-center gap-2 ${className}`}> {/* Merge className with default */}
      {name}
      <input
        type={type}
        className={
          type !== "textarea"
            ? `grow ${className}`
            : `grow textarea textarea-ghost textarea-xs w- max-w-xs ${className}`
        }
        placeholder={placeholder}
        required={required}
        value={value}
        onChange={onChange}
      />
    </label>
  );
};

export default FieldInput;

