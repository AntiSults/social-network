import { useState } from 'react';
import { User, Group } from "@/app/utils/types"; // Assuming User and Group types are already defined

interface InviteToGroupProps {
    groupId: number; // The selected group ID
    invitedUser: User; // The user to be invited
    currentUser: User | null; // The current logged-in user (group creator)
}

const InviteToGroup: React.FC<InviteToGroupProps> = ({ groupId, invitedUser, currentUser }) => {
    const [error, setError] = useState<string | null>(null);
    const [inviteSent, setInviteSent] = useState<boolean>(false);

    const handleInvite = async () => {
        if (!currentUser) {
            setError('User not logged in.');
            return;
        }

        try {
            await fetch(`http://localhost:8080/groups/invite`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ groupId, invitedUserId: invitedUser.ID, inviterId: currentUser.ID }),
            });
            setInviteSent(true);
        } catch (err) {
            console.error('Failed to send invite', err);
            setError('Failed to send invite');
        }
    };

    return (
        <div>
            {error && <p>{error}</p>}
            {!inviteSent ? (
                <button onClick={handleInvite}>
                    Invite {invitedUser.firstName} to Group
                </button>
            ) : (
                <p>Invite sent to {invitedUser.firstName}!</p>
            )}
        </div>
    );
};

export default InviteToGroup;
