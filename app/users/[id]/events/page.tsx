"use client"
import NavBar from "@/app/components/NavBar";
import { useUser } from "@/app/context/UserContext";
import { useEffect, useState } from 'react';
import EventReactions from '@/app/components/EventReactions';



const UserEvents = () => {

    const { user } = useUser();
    const [events, setEvents] = useState([]);

    useEffect(() => {
        const fetchEvents = async () => {
            const res = await fetch(`/api/users/${user?.ID}/events`);
            const data = await res.json();
            setEvents(data);
        };
        fetchEvents();
    }, [user?.ID]);

    return (
        <div>
            <h2>Your Events</h2>
            <ul>
                {events.map(event => (
                    <li key={event.id}>
                        <h3>{event.title}</h3>
                        <p>{event.description}</p>
                        <EventReactions eventId={event.id} />
                    </li>
                ))}
            </ul>
        </div>
    );
};

export default UserProfileEvents;
