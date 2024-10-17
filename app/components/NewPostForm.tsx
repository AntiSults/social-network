import React, { useEffect, useState } from 'react';
import { Group, User } from '@/app/utils/types';
import { getUserGroups } from '@/app/lib/api';

interface Props {
    onPostCreated: (content: string, privacy: string, file?: File | null, groupId?: number | null, visibleUsers?: number[]) => void;
    user: User | null;
}

const NewPostForm: React.FC<Props> = ({ onPostCreated, user }) => {
    const [content, setContent] = useState("");
    const [privacy, setPrivacy] = useState("public");
    const [file, setFile] = useState<File | null>(null);
    const [groups, setGroups] = useState<Group[]>([]);
    const [followers, setFollowers] = useState<User[]>([]);
    const [selectedGroup, setSelectedGroup] = useState<number | null>(null);
    const [selectedFollowers, setSelectedFollowers] = useState<number[]>([]);

    useEffect(() => {
      console.log("Current user:", user);
        const fetchGroups = async () => {
            try {
                const data = await getUserGroups();
                setGroups(data);
            } catch (error) {
                console.error("Error fetching groups:", error);
            }
        };

        const fetchFollowers = async () => {
          if (user) {
              try {
                  const response = await fetch(`http://localhost:8080/followers/followersList?userId=${user.ID}`);
                  console.log("Request to fetch followers made.");
                  if (!response.ok) {
                      throw new Error("Failed to fetch followers");
                  }
                  const data = await response.json();
                  console.log("Followers data received:", data);
                  setFollowers(data.followers || []);
                  console.log("Updated followers state:", data.followers);
              } catch (error) {
                  console.error("Error fetching followers:", error);
              }
          }
      };

        fetchGroups();
        fetchFollowers();
    }, [user]);

    const handleFollowerChange = (followerId: number) => {
        setSelectedFollowers((prevSelected) => {
            if (prevSelected.includes(followerId)) {
                return prevSelected.filter((id) => id !== followerId);
            } else {
                return [...prevSelected, followerId];
            }
        });
    };

    const handleSubmit = (event: React.FormEvent) => {
        event.preventDefault();
        if (content.trim()) {
            
            if (privacy === "group" && !selectedGroup) {
                alert("Please select a group to post to.");
                return;
            }
            if (privacy === "almost private" && selectedFollowers.length === 0) {
                alert("Please select at least one follower to share with.");
                return;
            }
    
            const groupId = privacy === "group" ? selectedGroup : null;
    
            try {
                onPostCreated(content, privacy, file, groupId, selectedFollowers);
            } catch (error) {
                console.error("Error creating post:", error);
            }
    
            setPrivacy("public");
            setFile(null);
            setSelectedGroup(null);
            setSelectedFollowers([]);
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

            {groups && groups.length > 0 && privacy === "group" && (
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

            {privacy === "almost private" && (
                <div className="mt-3">
                    <h3 className="font-bold">Select Followers:</h3>
                    {followers.length > 0 ? (
                        followers.map((user) => (
                            <div key={user.ID} className="flex items-center">
                                <input
                                    type="checkbox"
                                    checked={selectedFollowers.includes(user.ID)}
                                    onChange={() => handleFollowerChange(user.ID)}
                                    className="mr-2"
                                />
                                {user.firstName} {user.lastName}
                            </div>
                        ))
                    ) : (
                        <p className="text-gray-500">No followers available.</p>
                    )}
                </div>
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