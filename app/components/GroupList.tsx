import JoinGroup from './JoinGroup';
import { useState, useEffect } from 'react';
import { useUser } from '@/app/context/UserContext';
import { User, Group } from "@/app/utils/types";

interface GroupListProps {
    onSelectGroup: (groupId: number) => void;  // Pass group ID when selected
    actionType: 'invite' | 'createEvent' | 'chat'; // Add 'chat' for group chat selection
}

const GroupList: React.FC<GroupListProps> = ({ onSelectGroup, actionType }) => {
    const { user } = useUser();  // Current logged-in user context
    const [groups, setGroups] = useState<Group[]>([]);
    const [error, setError] = useState<string | null>(null);
    const [showMembers, setShowMembers] = useState<boolean>(false);
    const [members, setMembers] = useState<User[]>([]);
    const [selectedGroupId, setSelectedGroupId] = useState<number | null>(null);
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

    const fetchGroupMembers = async (groupId: number) => {
        try {
            const response = await fetch(`http://localhost:8080/groups/members?groupId=${groupId}`);
            const data = await response.json();
            setMembers(data);
            setShowMembers(true);
        } catch (err) {
            console.error('Failed to fetch group members', err);
        }
    };

    const handleViewMembers = (groupId: number) => {
        setSelectedGroupId(groupId);
        fetchGroupMembers(groupId);
    };

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
                        const isCreator = isGroupCreator(group) && isUserInGroup(group);

                        return (
                            <li key={`group-${group.id}`} className="p-4 border border-gray-300 rounded-md">
                                <div>
                                    <h3 className="text-lg font-semibold">{group.name}</h3>
                                    <p className="text-gray-600 font-semibold">{group.description}</p>

                                    <button
                                        onClick={() => handleViewMembers(group.id)}
                                        className="bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded-lg"
                                    >
                                        View Members
                                    </button>

                                    {alreadyInGroup || isCreator ? (
                                        <p className="text-gray-500">
                                            {isCreator ? 'You are the creator' : 'Already a member'}
                                        </p>
                                    ) : (
                                        <JoinGroup groupId={group.id} currentUser={user as User | null} />
                                    )}

                                    {/* Button to select group for the desired action */}
                                    {(isCreator || alreadyInGroup) && (
                                        <button
                                            onClick={() => onSelectGroup(group.id)}  // Call onSelectGroup with the group ID
                                            className="bg-gray-600 hover:bg-gray-700 text-white font-bold py-2 px-6 rounded-lg shadow-md transition duration-200 ease-in-out"
                                        >
                                            {actionType === 'createEvent'
                                                ? 'Select Group for Event Creation'
                                                : actionType === 'invite'
                                                    ? 'Select Group for Sending Invite'
                                                    : 'Select Group for Chat'}
                                        </button>
                                    )}
                                </div>
                            </li>
                        );
                    })}
                </ul>
            )}
            {showMembers && (
                <div className="fixed inset-0 bg-gray-500 bg-opacity-75 flex items-center justify-center">
                    <div className="bg-white p-6 rounded-lg">
                        <h3 className="text-lg font-semibold">Group Members</h3>
                        <ul>
                            {members.map((member) => (
                                <li key={member.ID} className="text-gray-700">
                                    {member.firstName} {member.lastName}
                                </li>
                            ))}
                        </ul>
                        <button onClick={() => setShowMembers(false)} className="mt-4 text-blue-500">
                            Close
                        </button>
                    </div>
                </div>
            )}


        </div>
    );
};

export default GroupList;
