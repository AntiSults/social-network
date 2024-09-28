"use client";

import UserProfileEvents from '@/app/components/UserProfileEvents';
import GroupList from '@/app/components/GroupList';
import CreateEventForm from '@/app/components/CreateEventForm';
import { useUser } from '@/app/context/UserContext';
import { useState } from 'react';

const EventsPage = () => {
    const { user } = useUser();  // Fetch current user from context
    const [selectedGroupId, setSelectedGroupId] = useState<number | null>(null);  // To track selected group for event creation

    if (!user) {
        return <p>Loading...</p>;
    }

    const handleGroupSelect = (groupId: number) => {
        setSelectedGroupId(groupId);  // Set selected group ID when a group is chosen
    };

    return (
        <div>
            <h1>Events</h1>
            {/* Display the list of available groups */}

            <GroupList onSelectGroup={handleGroupSelect} actionType="createEvent" />

            {/* Display the form to create an event if a group is selected */}
            {selectedGroupId && (
                <div>
                    <h2>Create an Event for Group ID: {selectedGroupId}</h2>
                    <CreateEventForm groupId={selectedGroupId} />
                </div>
            )}
            {/* Fetch and display the user's events */}
            <UserProfileEvents />
        </div>
    );
};

export default EventsPage;
