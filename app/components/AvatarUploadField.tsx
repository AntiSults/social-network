import React, { useState } from "react";
import Image from 'next/image';
import { useUser } from "@/app/context/UserContext";

interface Props {
  onFileSelect: (file: File) => void;
}

const AvatarUploadField = ({ onFileSelect }: Props) => {
  const [preview, setPreview] = useState<string | null>(null);
  const { user } = useUser();

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      setPreview(URL.createObjectURL(file));
      onFileSelect(file);
    }
  };

  return (
    <>
      <label className="form-control w-full max-w-xs">
        <div className="label">
          <span className="label-text">Choose your avatar (Optional)</span>
        </div>
        <input
          className="file-input file-input-bordered w-full"
          type="file"
          accept="image/*"
          onChange={handleFileChange}
        />
        {preview ? (
          <div className="avatar">
            <div className="w-24 rounded">
              <Image src={preview}
                alt="Avatar Preview"
                width={400}
                height={400}
              />
            </div>
          </div>
        ) : user && user.avatarPath ? (
          <div className="avatar">
            <div className="w-24 rounded">
              <Image src={user.avatarPath}
                alt="Current Avatar"
                width={400}
                height={400}
              />
            </div>
          </div>
        ) : null}
      </label>
    </>
  );
};

export default AvatarUploadField;
