import React from "react";

interface Props {
  name: string;
  type: string;
  placeholder: string;
  required: boolean;
  value: string;
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
}

const FieldInput = ({
  name,
  type,
  placeholder,
  required,
  value,
  onChange,
}: Props) => {
  return (
    <label className="input input-bordered flex items-center gap-2">
      {name}
      <input
        type={type}
        className={
          type != "textarea"
            ? "grow"
            : "grow textarea textarea-ghost textarea-xs w-full max-w-xs"
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
