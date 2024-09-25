"use client";
import React, { useState, useEffect } from "react";
import { User, Group } from "../utils/types";
interface PendingRequestsProps {
    group: Group;
}
const PendingRequests: React.FC<PendingRequestsProps> = ({ group }) => {
    const [requests, setRequests] = useState<Group[]>([]);

    useEffect(() => {
        const fetchRequests = async () => {
            try {
                const response = await fetch(`http://localhost:8080/groups/pending-requests?groupId=${group.id}`);
                const data = await response.json();
                setRequests(data);
            } catch (err) {
                console.error('Error fetching pending group join request', err);
            }

            fetchRequests();
        }, [user]);

    const handleDecision = (userId: number, accept: boolean) => {
        handleRequestDecision(group.id, userId, accept);
    };

    const handleRequestDecision = async (groupId: number, userId: number, accept: boolean) => {
        try {
            await fetch(`http://localhost:8080/groups/handle-request`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ groupId, userId, accept }),
            });
            console.log('Request handled!');
        } catch (err) {
            console.error('Failed to handle request', err);
        }
    };

    return (
        <div>
            <h3>Pending Requests</h3>
            <ul>
                {requests.map(req => (
                    <li key={req.userId}>
                        User {req.userId} wants to join/invited
                        <button onClick={() => handleDecision(req.userId, true)}>Accept</button>
                        <button onClick={() => handleDecision(req.userId, false)}>Reject</button>
                    </li>
                ))}
            </ul>
        </div>
    );
};