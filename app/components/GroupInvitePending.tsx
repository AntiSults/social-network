import { useEffect, useState } from 'react';
import { GroupJoinRequest } from '@/app/utils/types';
import { useUser } from '@/app/context/UserContext';

const PendingGroupInvites = () => {
    const { user } = useUser();
    const [invitations, setInvitations] = useState<GroupJoinRequest[]>([]);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchPendingInvitations = async () => {
            try {
                const response = await fetch(`http://localhost:8080/groups/pending-invites?userID=${user?.ID}`);
                const data = await response.json();

                if (data && Array.isArray(data)) {
                    setInvitations(data);
                } else if (data === null || data.length === 0) {
                    setInvitations([]);
                } else {
                    setError('Unexpected response format: expected an array.');
                    console.error('Unexpected data format', data);
                }
            } catch (err) {
                console.error('Failed to fetch pending invitations', err);
                setError('Failed to fetch pending invitations');
            }
        };

        if (user?.ID) {
            fetchPendingInvitations();
        }
    }, [user]);

    const handleInvitationDecision = async (groupId: number, userId: number | undefined, accept: boolean) => {

        try {
            await fetch(`http://localhost:8080/groups/handle-invites`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ groupId, userId, accept }),
            });
            setInvitations(invitations.filter(invite => invite.group_id !== groupId));
        } catch (err) {
            console.error('Failed to handle invitation', err);
            setError('Failed to handle invitation');
        }
    };
    if (!invitations.length) {
        return null;
    }
    return (
        <div className="flex flex-col items-center mt-10">
            <h2 className="text-2xl font-bold mb-4">Pending Group Invitations</h2>
            {error && <p>{error}</p>}
            {invitations.length === 0 ? (
                <p>No pending invitations</p>
            ) : (
                <ul className="w-full max-w-lg">
                    {invitations.map(invitation => (
                        <li
                            key={invitation.group_id}
                            className="flex justify-between items-center bg-white shadow-md rounded-lg p-4 mb-4"
                        >
                            <div>
                                <p className="font-semibold">
                                    {`You've been invited by ${invitation.first_name} ${invitation.last_name} to join ${invitation.group_name}`}
                                </p>
                            </div>
                            <button
                                onClick={() => handleInvitationDecision(invitation.group_id, user?.ID, true)}
                                className="bg-green-500 text-white py-2 px-4 rounded-lg hover:bg-green-600 transition"
                            >
                                Accept
                            </button>
                            <button
                                onClick={() => handleInvitationDecision(invitation.group_id, user?.ID, false)}
                                className="bg-red-500 text-white py-2 px-4 rounded-lg hover:bg-red-600 transition"
                            >
                                Decline
                            </button>

                        </li>
                    ))}
                </ul>
            )}
        </div>
    );
};

export default PendingGroupInvites;
