"use client";
import JoinGroup from './JoinGroup';
import { useState, useEffect } from 'react';
import { useUser } from '@/app/context/UserContext';
import { User, Group } from "@/app/utils/types";

interface GroupListProps {
    onSelectGroup: (groupId: number) => void; // Prop to pass selected group for invitation
}

const GroupList: React.FC<GroupListProps> = ({ onSelectGroup }) => {
    const { user } = useUser(); // Current logged-in user context
    const [groups, setGroups] = useState<Group[]>([]);
    const [error, setError] = useState<string | null>(null);
    const currentUser = user;

    useEffect(() => {
        const fetchGroups = async () => {
            try {
                const response = await fetch(`http://localhost:8080/groups`);
                const data = await response.json();
                setGroups(data);
            } catch (err) {
                console.error('Failed to fetch groups', err);
                setError('Failed to fetch groups');
            }
        };
        fetchGroups();
    }, []);

    const isUserInGroup = (group: Group) => {
        return group.members?.includes(currentUser?.ID ?? -1);
    };

    const isGroupCreator = (group: Group) => {
        return group.creator_id === currentUser?.ID;
    };
    const safeGroups = Array.isArray(groups) ? groups : [];

    return (
        <div className="bg-white shadow-md rounded-lg p-6 mb-4">
            <h2 className="text-xl font-bold mb-4">All Groups</h2>
            {error && <p>{error}</p>}
            {safeGroups.length === 0 ? (
                <p className="text-gray-500">No groups available.</p>
            ) : (
                <ul className="space-y-4">
                    {groups.map((group) => {
                        if (!group || !group.id) {
                            return null;
                        }
                        const alreadyInGroup = isUserInGroup(group);
                        const isCreator = isGroupCreator(group);

                        return (
                            <li key={`group-${group.id}`} className="p-4 border border-gray-300 rounded-md">
                                <div>
                                    <h3 className="text-lg font-semibold">{group.name}</h3>
                                    <p className="text-gray-600  font-semibold">{group.description}</p>

                                    {alreadyInGroup || isCreator ? (
                                        <p className="text-gray-500">
                                            {isCreator ? 'You are the creator' : 'Already a member'}
                                        </p>
                                    ) : (
                                        <JoinGroup groupId={group.id} currentUser={user as User | null} />
                                    )}
                                    {/* Button to open the group */}
                                    <button
                                        onClick={() => {/* Add logic to open the group */}}
                                        className="bg-gray-600 hover:bg-gray-700 text-white font-bold py-2 px-6 rounded-lg shadow-md transition duration-200 ease-in-out"
                                    >
                                        Open Group
                                    </button>

                                    {/* Button to select group for inviting a user */}
                                    {(isCreator || alreadyInGroup) && (
                                        <button
                                            onClick={() => {onSelectGroup(group.id)}}
                                            className="bg-gray-600 hover:bg-gray-700 text-white font-bold py-2 px-6 rounded-lg shadow-md transition duration-200 ease-in-out"
                                        >
                                            Select Group for Invitation
                                        </button>
                                    )}
                                </div>
                            </li>
                        );
                    })}
                </ul>
            )}
        </div>
    );
};

export default GroupList;
