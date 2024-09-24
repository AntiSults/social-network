"use client"
import { useState, useEffect } from 'react';
import { useUser } from '../context/UserContext';
import { Group } from "../utils/types"

const GroupList = () => {
    const { user } = useUser(); // Current logged-in user context
    const [groups, setGroups] = useState<Group[]>([]);
    const [error, setError] = useState<string | null>(null);
    const currentUser = user
    useEffect(() => {
        const fetchGroups = async () => {
            try {
                const response = await fetch(`http://localhost:8080/groups`);
                const data = await response.json();
                console.log("Group Data", data)
                setGroups(data);
            } catch (err) {
                console.error('Failed to fetch groups', err);
                setError('Failed to fetch groups');
            }
        };
        fetchGroups();
    }, []);

    const handleJoinGroup = async (groupId: number) => {
        try {
            await fetch(`http://localhost:8080/groups/join`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ groupId, userId: currentUser?.ID }),
            });
            // Optionally refresh the groups or show a success message
        } catch (err) {
            console.error('Failed to join group', err);
        }
    };
    return (
        <div>
            <h2>All Groups</h2>
            {error && <p>{error}</p>}
            {groups.length === 0 ? (
                <p>No groups available.</p>
            ) : (
                <ul>
                    {groups.map(group => {
                        if (!group || !group.id) {
                            return null;
                        }
                        return (
                            <li key={`group-${group.id}`}>
                                <div>
                                    <h3>{group.name}</h3>
                                    <p>{group.description}</p>
                                    <button onClick={() => handleJoinGroup(group.id)}>Join Group</button>
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