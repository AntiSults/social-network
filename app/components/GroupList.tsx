"use client";
import JoinGroup from './JoinGroup';
import { useState, useEffect } from 'react';
import { useUser } from '@/app/context/UserContext';
import { User, Group } from "@/app/utils/types";

const GroupList = () => {
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
        // Check if current user is in the members list
        return group.members?.includes(currentUser?.ID ?? -1);
    };

    const isGroupCreator = (group: Group) => {
        return group.creator_id === currentUser?.ID;
    };

    return (
        <div>
            <h2>All Groups</h2>
            {error && <p>{error}</p>}
            {groups.length === 0 ? (
                <p>No groups available.</p>
            ) : (
                <ul>
                    {groups.map((group) => {
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
