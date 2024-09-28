"use client"
import { useState } from 'react';

const CreateEventForm = ({ groupId }: { groupId: number }) => {
    const [title, setTitle] = useState('');
    const [description, setDescription] = useState('');
    const [eventDate, setEventDate] = useState('');

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        await fetch(`http://localhost:8080/groups/events`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ title, description, eventDate }),
        });
    };

    return (
        <form onSubmit={handleSubmit}>
            <input
                type="text"
                placeholder="Event Title"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                required
            />
            <textarea
                placeholder="Event Description"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
            />
            <input
                type="datetime-local"
                value={eventDate}
                onChange={(e) => setEventDate(e.target.value)}
                required
            />
            <button type="submit">Create Event</button>
        </form>
    );
};

export default CreateEventForm;
