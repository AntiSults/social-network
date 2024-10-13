import { useState } from 'react';
import { User } from '@/app/utils/types';

interface Props {
    groupId: number;
    currentUser: User | null;
}

const JoinGroup: React.FC<Props> = ({ groupId, currentUser }) => {
    const [error, setError] = useState<string | null>(null);
    const [requestSent, setRequestSent] = useState<boolean>(false); // State to track if the request was sent

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
            setRequestSent(true); // Set state to true after successful request
        } catch (err) {
            console.error('Failed to send join request', err);
            setError('Failed to send join request');
        }
    };

    return (
        <div>
            {error && <p>{error}</p>}
            {!requestSent ? (
                <button onClick={handleJoinRequest} disabled={!currentUser}>
                    Join Group Request
                </button>
            ) : (
                <p>Join request sent!</p> // Display confirmation text after the request is sent
            )}
        </div>
    );
};

export default JoinGroup;

