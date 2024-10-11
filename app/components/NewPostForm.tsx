import React, { useState, useEffect } from "react";
import { Group } from "@/app/utils/types";
import { getUserGroups } from "@/app/lib/api";



interface Props {
  onPostCreated: (content: string, privacy: string, file?: File | null, groupId?: number | null) => void;
}

const NewPostForm: React.FC<Props> = ({ onPostCreated }) => {
  const [content, setContent] = useState("");
  const [privacy, setPrivacy] = useState("public"); // "public" by default
  const [file, setFile] = useState<File | null>(null);
  const [groups, setGroups] = useState<Group[]>([]);
  const [selectedGroup, setSelectedGroup] = useState<number | null>(null);

  useEffect(() => {
    const fetchGroups = async () => {
      try {
        const data = await getUserGroups();
        setGroups(data);
      } catch (error) {
        console.error("Error fetching groups:", error);
      }
    };

    fetchGroups();
  }, []);

  const handleSubmit = (event: React.FormEvent) => {
    event.preventDefault();
    if (content.trim()) {
      if (privacy === "group" && !selectedGroup) {
        alert("Please select a group to post to.");
        return;
      }

      const groupId = privacy === "group" ? selectedGroup : null;
      onPostCreated(content, privacy, file, groupId);
      setPrivacy("public");
      setFile(null);
      setSelectedGroup(null);
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
        <option value="group">Group</option>

      </select>

      {groups.length > 0 && privacy === "group" && (
        <select
          value={selectedGroup || ""}
          onChange={(e) => setSelectedGroup(Number(e.target.value))}
          className="w-full mt-3 p-2 border rounded-md border-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <option value="">Select a group</option>
          {groups.map((group) => (
            <option key={group.id} value={group.id}>
              {group.name}
            </option>
          ))}
        </select>
      )}

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
