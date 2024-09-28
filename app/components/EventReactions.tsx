"use client";

import { useState } from 'react';
import { useUser } from '@/app/context/UserContext';


const EventReactions = ({ eventId }: { eventId: number }) => {
    const { user } = useUser(); // Get the logged-in user from context
    const [error, setError] = useState<string | null>(null);
    const [status, setStatus] = useState<string | null>(null); // Feedback message for user action
    const [loading, setLoading] = useState(false); // Loading state for API requests

    const handleReaction = async (reaction: string) => {
        if (!user) {
            setError('You must be logged in to respond to an event.');
            return;
        }

        setLoading(true);
        setError(null);
        setStatus(null);

        try {
            const response = await fetch(`http://localhost:8080/groups/events-react`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    userId: user.ID,
                    eventId,
                    reaction
                }),
            });

            if (!response.ok) {
                throw new Error('Failed to register your response.');
            }

            setStatus(`You have marked yourself as "${reaction}" for this event.`);
        } catch (err) {
            setError("No response registered");
        } finally {
            setLoading(false);
        }
    };

    return (
        <div>
            <h3>Will you attend?</h3>
            {error && <p className="text-red-500">{error}</p>}
            {status && <p className="text-green-500">{status}</p>}

            <div className="mt-2">
                <button
                    onClick={() => handleReaction('going')}
                    className={`bg-blue-500 text-white px-4 py-2 mr-2 ${loading ? 'opacity-50 cursor-not-allowed' : ''}`}
                    disabled={loading}
                >
                    Going
                </button>

                <button
                    onClick={() => handleReaction('not going')}
                    className={`bg-gray-500 text-white px-4 py-2 ${loading ? 'opacity-50 cursor-not-allowed' : ''}`}
                    disabled={loading}
                >
                    Not Going
                </button>
            </div>
        </div>
    );
};

export default EventReactions;
