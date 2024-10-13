import React, { useState, useEffect } from 'react';
import { User } from '@/app/utils/types';
interface Props {
    user: User;
}
const PendingRequests: React.FC<Props> = ({ user }) => {
    const [pendingRequests, setPendingRequests] = useState<User[]>([]);
    useEffect(() => {
        const fetchPendingRequests = async () => {
            try {

                const response = await fetch(`http://localhost:8080/followers/pending?userId=${user.ID}`);
                const data = await response.json();

                setPendingRequests(data);
            } catch (error) {
                console.error("Error fetching pending follow requests:", error);
            }
        };
        if (user?.ID) {
            fetchPendingRequests();
        }
    }, [user]);
    const handleAccept = async (userId: number, followerId: number) => {
        try {
            const response = await fetch(`http://localhost:8080/followers/accept`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ userId, followerId }),
            });
            if (response.ok) {
                // Remove the accepted request from pending requests
                setPendingRequests((prev) => prev.filter((req) => req.ID !== followerId));
            } else {
                console.error("Failed to accept follow request");
            }
        } catch (error) {
            console.error("Error accepting follow request:", error);
        }
    };
    const handleReject = async (userId: number, followerId: number) => {
        try {
            const response = await fetch(`http://localhost:8080/followers/reject`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ userId, followerId }),
            });
            if (response.ok) {
                // Remove the rejected request from pending requests
                setPendingRequests((prev) => prev.filter((req) => req.ID !== followerId));
            } else {
                console.error("Failed to reject follow request");
            }
        } catch (error) {
            console.error("Error rejecting follow request:", error);
        }
    };
    if (pendingRequests.length === 0) {
        return null;
    }
    return (
        <div className="flex flex-col items-center mt-10">
            <h2 className="text-2xl font-bold mb-4">Pending Follow Requests</h2>
            <ul className="w-full max-w-lg">
                {pendingRequests.map((request) => (
                    <li

                        key={request.ID}
                        className="flex justify-between items-center bg-white shadow-md rounded-lg p-4 mb-4"
                    >
                        <div>
                            <p className="font-semibold">
                                {request.firstName} {request.lastName}
                            </p>
                        </div>
                        <div className="flex space-x-4">
                            <button
                                onClick={() => handleAccept(user.ID, request.ID)}
                                className="bg-green-500 text-white py-2 px-4 rounded-lg hover:bg-green-600 transition"
                            >
                                Accept
                            </button>
                            <button
                                onClick={() => handleReject(user.ID, request.ID)}
                                className="bg-red-500 text-white py-2 px-4 rounded-lg hover:bg-red-600 transition"
                            >
                                Reject
                            </button>
                        </div>
                    </li>
                ))}
            </ul>
        </div>
    );

};
export default PendingRequests;
