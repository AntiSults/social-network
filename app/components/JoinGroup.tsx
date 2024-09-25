import { useState } from 'react';
import { useUser } from '../context/UserContext'; // For accessing currentUser
import { User, Group } from "../utils/types"; // Import interfaces from types

interface JoinGroupProps {
    groupId: number;
    currentUser: User | null; // Use the User type from your utils/types
}

const JoinGroup: React.FC<JoinGroupProps> = ({ groupId, currentUser }) => {
    const [error, setError] = useState<string | null>(null);

    const handleJoinRequest = async () => {
        if (!currentUser) {
            setError('User not logged in.');
            return;
        }

        try {
            await fetch(`http://localhost:8080/groups/join-request`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ groupId, userId: currentUser.ID }),
            });
            console.log('Join request sent!');
        } catch (err) {
            console.error('Failed to send join request', err);
            setError('Failed to send join request');
        }
    };

    return (
        <div>
            {error && <p>{error}</p>}
            <button onClick={handleJoinRequest} disabled={!currentUser}>
                Join Group Request
            </button>
        </div>
    );
};

export default JoinGroup;
