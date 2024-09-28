"use client"
import { useState } from 'react';
import FieldInput from './FieldInput';
import Button from "@/app/components/Button";

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
            <FieldInput
                type="text"
                placeholder="Event Title"
                required={true}
                value={title}
                onChange={(e) => setTitle(e.target.value)}
            />
            <FieldInput
                type="text"
                placeholder="Event Description"
                required={true}
                value={description}
                onChange={(e) => setDescription(e.target.value)}
            />
            <FieldInput
                type="datetime-local"
                required={true}
                placeholder=""
                value={eventDate}
                onChange={(e) => setEventDate(e.target.value)}
            />
            <Button type="submit" name="Create Event" />
        </form>
    );
};

export default CreateEventForm;
