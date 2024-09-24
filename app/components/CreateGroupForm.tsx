"use client";

import { useState } from 'react';
import { useUser } from '../context/UserContext';

const CreateGroupForm = () => {
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState(false);
    const { user } = useUser();
    const currentUser = user

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        try {
            const response = await fetch('http://localhost:8080/groups', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    name,
                    description,
                    creator_id: currentUser?.ID,
                }),
            });

            if (response.ok) {
                setSuccess(true);
                setName('');
                setDescription('');
            } else {
                const errorData = await response.json();
                setError(errorData.error || 'Error creating group');
            }
        } catch (err) {
            setError('Failed to create group');
        }
    };

    return (
        <div>
            {success && <p>Group created successfully!</p>}
            {error && <p>{error}</p>}
            <form onSubmit={handleSubmit}>
                <label>
                    Group Name:
                    <input
                        type="text"
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        required
                    />
                </label>
                <label>
                    Description:
                    <textarea
                        value={description}
                        onChange={(e) => setDescription(e.target.value)}
                        required
                    />
                </label>
                <button type="submit">Create Group</button>
            </form>
        </div>
    );
};

export default CreateGroupForm;
