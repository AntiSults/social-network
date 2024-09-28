"use client";
import { useState } from 'react';
import { useUser } from '@/app/context/UserContext';
import GroupList from '@/app/components/GroupList';
import CreateEventForm from '@/app/components/CreateEventForm';

const EventsPage = () => {
    const { user } = useUser();
    const [selectedGroupId, setSelectedGroupId] = useState<number | null>(null);

    const handleSelectGroup = (groupId: number) => {
        setSelectedGroupId(groupId);
    };

    if (!user) {
        return <p>Loading...</p>;
    }

    return (
        <div>
            <h1>Events</h1>

            {/* Display the list of groups */}
            <GroupList onSelectGroup={handleSelectGroup} />

            {/* Display the create event form if a group is selected */}
            {selectedGroupId && (
                <div>
                    <h2>Create an Event for Group ID: {selectedGroupId}</h2>
                    <CreateEventForm groupId={selectedGroupId} />
                </div>
            )}
        </div>
    );
};

export default EventsPage;


