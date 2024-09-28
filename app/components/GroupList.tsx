"use client";
import JoinGroup from './JoinGroup';
import { useState, useEffect } from 'react';
import { useUser } from '@/app/context/UserContext';
import { User, Group } from "@/app/utils/types";

interface GroupListProps {
    onSelectGroup: (groupId: number) => void; // Prop to pass selected group for event creation
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
        <div>
            <h2>All Groups</h2>
            {error && <p>{error}</p>}
            {safeGroups.length === 0 ? (
                <p>No groups available.</p>
            ) : (
                <ul>
                    {safeGroups.map((group) => {
                        if (!group || !group.id) {
                            return null;
                        }
                        const alreadyInGroup = isUserInGroup(group);
                        const isCreator = isGroupCreator(group);

                        return (
                            <li key={`group-${group.id}`}>
                                <div>
                                    <h3>{group.name}</h3>
                                    <p>{group.description}</p>

                                    {alreadyInGroup || isCreator ? (
                                        <p>{isCreator ? 'You are the creator' : 'Already a member'}</p>
                                    ) : (
                                        <JoinGroup groupId={group.id} currentUser={user as User | null} />
                                    )}

                                    {/* Button to select group for creating an event */}
                                    {(isCreator || alreadyInGroup) && (
                                        <button
                                            onClick={() => onSelectGroup(group.id)}
                                            className="bg-blue-500 text-white px-4 py-2 mt-2"
                                        >
                                            Select Group for Event Creation
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
