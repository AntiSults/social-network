import { imageConfigDefault } from "next/dist/shared/lib/image-config";
import React, { useState } from "react";

interface Props {
  onFileSelect: (file: File) => void;
}

const AvatarUploadField = ({ onFileSelect }: Props) => {
  const [preview, setPreview] = useState<string | null>(null);

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      setPreview(URL.createObjectURL(file));
      console.log(URL.createObjectURL(file));
      onFileSelect(file);
    }
  };

  return (
    <>
      <label>
        Choose your avatar (Optional)
        <input
          className="avatarInput"
          type="file"
          accept="image"
          onChange={handleFileChange}
        />
        ;
        {preview && (
          <img
            src={preview}
            alt="Avatar Preview"
            style={{ width: 100, height: 100, borderRadius: "50%" }}
          />
        )}
      </label>
    </>
  );
};

export default AvatarUploadField;
