"use client";

import { useEffect, useState } from 'react';
import { Event } from '@/app/utils/types';
import EventReactions from '@/app/components/EventReactions';
import CreateEventForm from '@/app/components/CreateEventForm';
import { useUser } from '@/app/context/UserContext';

const UserProfileEvents = () => {
    const { user } = useUser();  // Fetch the current user from context
    const [events, setEvents] = useState<Event[] | null>(null); // State for events
    const [error, setError] = useState<string | null>(null); // State for errors
    const [loading, setLoading] = useState<boolean>(true); // Loading state for fetch

    useEffect(() => {
        const fetchEvents = async () => {
            if (!user?.ID) return;  // Return if no user is available

            try {
                const res = await fetch(`http://localhost:8080/groups/events?userID=${user.ID}`);
                if (!res.ok) {
                    throw new Error('Failed to fetch events');
                }
                const data: Event[] = await res.json();  // Fetch the event data
                setEvents(data && Array.isArray(data) ? data : []);
            } catch (err) {
                console.error('Error fetching events:', err);
                setError('Failed to load events');
            } finally {
                setLoading(false); // Mark loading as complete
            }
        };

        fetchEvents();
    }, [user]);

    if (loading) {
        return <p>Loading events...</p>;
    }

    if (error) {
        return <p>{error}</p>;
    }

    if (events && events.length === 0) {
        return <p>No events available</p>;
    }

    return (
        <div>
            <h2>Your Events</h2>
            <ul>
                {events?.map(event => (
                    <li key={event.id}>
                        <h3>{event.title}</h3>
                        <p>{event.description}</p>
                        <p>Date: {new Date(event.eventDate).toLocaleString()}</p>
                        {/* Event reactions */}
                        <EventReactions eventId={event.id} />
                        {/* Create a new event within the same group */}
                        <CreateEventForm groupId={event.groupId} />
                    </li>
                ))}
            </ul>
        </div>
    );
};

export default UserProfileEvents;
