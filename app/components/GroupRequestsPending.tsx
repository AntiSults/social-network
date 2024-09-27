"use client";
import { useEffect, useState } from 'react';
import { GroupJoinRequest } from '@/app/utils/types';
import { useUser } from '@/app/context/UserContext';

const PendingGroupRequests = () => {
    const { user } = useUser();
    const [requests, setRequests] = useState<GroupJoinRequest[]>([]);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchPendingRequests = async () => {
            try {
                const response = await fetch(`http://localhost:8080/groups/pending-requests?creatorID=${user?.ID}`);
                const data = await response.json();
                console.log("returned group join requests", data);

                if (data && Array.isArray(data)) {
                    setRequests(data);
                } else if (data === null || data.length === 0) {
                    setRequests([]);
                } else {
                    setError('Unexpected response format: expected an array.');
                    console.error('Unexpected data format', data);
                }
            } catch (err) {
                console.error('Failed to fetch pending requests', err);
                setError('Failed to fetch pending requests');
            }
        };

        if (user?.ID) {
            fetchPendingRequests();
        }
    }, [user]);

    const handleRequestDecision = async (groupId: number, userId: number, accept: boolean) => {
        try {
            await fetch(`http://localhost:8080/groups/handle-request`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ groupId, userId, accept }),
            });
            setRequests(requests.filter(request => !(request.group_id === groupId && request.user_id === userId)));
        } catch (err) {
            console.error('Failed to handle request', err);
            setError('Failed to handle request');
        }
    };
    if (!requests.length) {
        // Skip rendering anything if there are no pending requests
        return null;
    }
    return (
        <div>
            <h2>Pending Group Join Requests</h2>
            {error && <p>{error}</p>}
            {requests.length === 0 ? (
                <p>No pending requests</p>
            ) : (
                <ul>
                    {requests.map(request => (
                        <li key={`${request.group_id}-${request.user_id}`}>
                            {`${request.first_name} ${request.last_name} wants to join ${request.group_name}`}
                            <button onClick={() => handleRequestDecision(request.group_id, request.user_id, true)}>Accept</button>
                            <button onClick={() => handleRequestDecision(request.group_id, request.user_id, false)}>Reject</button>
                        </li>
                    ))}
                </ul>
            )}
        </div>
    );
};

export default PendingGroupRequests;
