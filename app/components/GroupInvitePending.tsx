"use client";
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
                console.log("invite data", data)
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
        console.log("sending back", groupId, userId, accept)
        try {
            await fetch(`http://localhost:8080/groups/handle-invites`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ groupId, userId, accept }),
            });
            setInvitations(invitations.filter(invite => !(invite.group_id == groupId && invite.user_id === userId)));
        } catch (err) {
            console.error('Failed to handle invitation', err);
            setError('Failed to handle invitation');
        }
    };

    if (!invitations.length) {
        return null;
    }

    return (
        <div>
            <h2>Pending Group Invitations</h2>
            {error && <p>{error}</p>}
            {invitations.length === 0 ? (
                <p>No pending invitations</p>
            ) : (
                <ul>
                    {invitations.map(invitation => (
                        <li key={invitation.group_id}>
                            {`You've been invited by ${invitation.first_name} ${invitation.last_name} to join ${invitation.group_name}`}
                            <button onClick={() => handleInvitationDecision(invitation.group_id, user?.ID, true)}>Accept</button>
                            <button onClick={() => handleInvitationDecision(invitation.group_id, user?.ID, false)}>Decline</button>
                        </li>
                    ))}
                </ul>
            )}
        </div>
    );
};

export default PendingGroupInvites;
