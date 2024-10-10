import { useState } from 'react';
import { useUser } from '@/app/context/UserContext';

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
        <div className="bg-white shadow-md rounded-lg p-6 mb-4">
            {success && <p className="text-green-500 font-semibold">Group created successfully!</p>}
            {error && <p className="text-red-500 font-semibold">{error}</p>}
            <form onSubmit={handleSubmit} className="flex flex-col space-y-4">
                <label className="flex flex-col">
                    <span className="font-semibold">Group Name:</span>
                    <input
                        type="text"
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        required
                        className="border border-gray-300 rounded-md p-2"
                    />
                </label>
                <label className="flex flex-col">
                    <span className="font-semibold">Description:</span>
                    <textarea
                        value={description}
                        onChange={(e) => setDescription(e.target.value)}
                        required
                        className="border border-gray-300 rounded-md p-2"
                    />
                </label>
                <div className="flex justify-center">
                    <button
                        type="submit"
                        className="bg-gray-600 hover:bg-gray-700 text-white font-bold py-2 px-6 rounded-lg shadow-md transition duration-200 ease-in-out"
                    >
                        Create Group
                    </button>
                </div>
            </form>
        </div>
    );
};

export default CreateGroupForm;
