"use client";
import React, { useState, useEffect } from "react";
import { User } from "../utils/types";

interface PendingRequestsProps {
    user: User;
}

const PendingRequests: React.FC<PendingRequestsProps> = ({ user }) => {
    const [pendingRequests, setPendingRequests] = useState<User[]>([]); // Keep this as an array to avoid map errors

    useEffect(() => {
        const fetchPendingRequests = async () => {
            try {
                // Fetch pending follow requests for the logged-in user
                const response = await fetch(`http://localhost:8080/followers/pending?userId=${user.ID}`);
                const data = await response.json();

                // Assuming you receive an array of users with pending requests
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
                body: JSON.stringify({ userId, followerId }), // Send both userId and followerId
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
                body: JSON.stringify({ userId, followerId }), // Send both userId and followerId
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
        <div>
            <h2>Pending Follow Requests</h2>
            <ul>
                {pendingRequests.map((request) => (
                    <li key={request.ID}>
                        {request.firstName} {request.lastName}
                        <button onClick={() => handleAccept(user.ID, request.ID)}>Accept</button>
                        <button onClick={() => handleReject(user.ID, request.ID)}>Reject</button>
                    </li>
                ))}
            </ul>
        </div>
    );

};

export default PendingRequests;
