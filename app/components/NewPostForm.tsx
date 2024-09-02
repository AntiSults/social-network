import React, { useState } from "react";

interface NewPostFormProps {
  onPostCreated: (content: string) => void;
}

const NewPostForm: React.FC<NewPostFormProps> = ({ onPostCreated }) => {
  const [content, setContent] = useState("");

  const handleSubmit = (event: React.FormEvent) => {
    event.preventDefault();
    if (content.trim()) {
      onPostCreated(content);
      setContent("");
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
