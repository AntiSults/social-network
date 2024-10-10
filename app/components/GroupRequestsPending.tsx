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
        return null;
    }
    return (
        <div className="flex flex-col items-center mt-10">
            <h2 className="text-2xl font-bold mb-4">Pending Group Join Requests</h2>

            {error && <p>{error}</p>}
            {requests.length === 0 ? (
                <p>No pending requests</p>
            ) : (
                <ul className="w-full max-w-lg">
                    {requests.map(request => (
                        <li key={`${request.group_id}-${request.user_id}`}
                            className="flex justify-between items-center bg-white shadow-md rounded-lg p-4 mb-4"
                        >
                            <div>
                                <p className="font-semibold">
                                    {`${request.first_name} ${request.last_name} wants to join ${request.group_name}`}
                                </p>

                            </div>
                            <button
                                onClick={() => handleRequestDecision(request.group_id, request.user_id, true)}
                                className="bg-green-500 text-white py-2 px-4 rounded-lg hover:bg-green-600 transition"
                            >
                                Accept
                            </button>
                            <button
                                onClick={() => handleRequestDecision(request.group_id, request.user_id, false)}
                                className="bg-red-500 text-white py-2 px-4 rounded-lg hover:bg-red-600 transition"
                            >
                                Reject
                            </button>
                        </li>
                    ))}
                </ul>
            )}
        </div>
    );
};

export default PendingGroupRequests;
