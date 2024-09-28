"use client"
import { useState } from 'react';
import FieldInput from './FieldInput';
import Button from "@/app/components/Button";

const CreateEventForm = ({ groupId }: { groupId: number }) => {
    const [title, setTitle] = useState('');
    const [description, setDescription] = useState('');
    const [eventDate, setEventDate] = useState('');
    const [submitError, setSubmitError] = useState<string | null>(null);
    const [success, setSuccess] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setSubmitError(null);  // Clear any previous errors

        try {
            const response = await fetch(`http://localhost:8080/groups/events`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ groupId, title, description, eventDate }),
            });

            if (response.ok) {
                setTitle('');
                setDescription('');
                setEventDate('');
                setSuccess(true);
            } else {
                const errorMessage = await response.text();
                setSubmitError(errorMessage || 'Failed to create event');
            }
        } catch (err) {
            console.error('Error submitting form:', err);
            setSubmitError('An error occurred while creating the event');
        }
    };

    return (
        <form onSubmit={handleSubmit}>
            {submitError && <p style={{ color: 'red' }}>{submitError}</p>}
            {success && <p style={{ color: 'green' }}>Event created successfully!</p>}

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
