import React, { useState } from "react";

interface NewPostFormProps {
  onPostCreated: (content: string, privacy: string, file?: File | null) => void;
}

const NewPostForm: React.FC<NewPostFormProps> = ({ onPostCreated }) => {
  const [content, setContent] = useState("");
  const [privacy, setPrivacy] = useState("public"); // "public" by default
  const [file, setFile] = useState<File | null>(null);

  const handleSubmit = (event: React.FormEvent) => {
    event.preventDefault();
    if (content.trim()) {
      onPostCreated(content, privacy, file);
      setContent("");
      setPrivacy("public");
      setFile(null);
    }
  };

  return (
    <form 
    className="relative mx-auto max-w-lg p-6 bg-white shadow-md rounded-lg mb-4"
    onSubmit={handleSubmit}
    >
      <textarea
        value={content}
        onChange={(e) => setContent(e.target.value)}
        placeholder="What's on your mind?"
        className="w-full p-3 border rounded-md border-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-500"
      />
      <select
        value={privacy}
        onChange={(e) => setPrivacy(e.target.value)}
        className="w-full mt-3 p-2 border rounded-md border-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        <option value="public">Public</option>
        <option value="private">Private</option>
        <option value="almost private">Almost Private</option>
      </select>
      <input
        type="file"
        accept="image/*, .gif"
        onChange={(e) => setFile(e.target.files ? e.target.files[0] : null)}
        className="mt-3"
      />
      <button
        type="submit"
        className="bg-gray-600 text-white px-4 py-2 rounded hover:bg-gray-700"
      >
        Create Post
      </button>
    </form>
  );
};

export default NewPostForm;
