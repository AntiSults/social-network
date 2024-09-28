"use client";

import { useEffect, useState } from 'react';
import { Event } from '@/app/utils/types';
import EventReactions from '@/app/components/EventReactions';
import CreateEventForm from '@/app/components/CreateEventForm';
import { useUser } from '@/app/context/UserContext';

const UserProfileEvents = () => {
    const { user } = useUser();
    const [events, setEvents] = useState<Event[] | null>(null); // Events state
    const [error, setError] = useState<string | null>(null); // Error state

    useEffect(() => {
        const fetchEvents = async () => {
            if (!user?.ID) return; // Ensure user is available

            try {
                const res = await fetch(`http://localhost:8080/groups/events?userID=${user.ID}`);
                if (!res.ok) {
                    throw new Error('Failed to fetch events');
                }
                const data: Event[] = await res.json(); // Fetch data as Event array
                setEvents(data && Array.isArray(data) ? data : []); // Set events or empty array
            } catch (err) {
                console.error('Error fetching events:', err);
                setError('Failed to load events');
            }
        };

        fetchEvents();
    }, [user]);

    if (error) {
        return <p>{error}</p>;
    }

    if (events === null) {
        return <p>Loading events...</p>;
    }

    if (events.length === 0) {
        return <p>No events available</p>;
    }

    return (
        <div>
            <h2>Your Events</h2>
            <ul>
                {events.map(event => (
                    <li key={event.id}>
                        <h3>{event.title}</h3>
                        <p>{event.description}</p>
                        <p>Date: {new Date(event.eventDate).toLocaleString()}</p>
                        {/* Pass groupId and userId to EventReactions */}
                        <EventReactions eventId={event.id} />
                        {/* Option to create a new event in the same group */}
                        <CreateEventForm groupId={event.groupId} />
                    </li>
                ))}
            </ul>
        </div>
    );
};

export default UserProfileEvents;

