"use client";

import UserProfileEvents from '@/app/components/UserProfileEvents';
import GroupList from '@/app/components/GroupList';
import CreateEventForm from '@/app/components/CreateEventForm';
import { useUser } from '@/app/context/UserContext';
import { useState } from 'react';
import NavBar from "@/app/components/NavBar";

const EventsPage = () => {
    const { user } = useUser();
    const [selectedGroupId, setSelectedGroupId] = useState<number | null>(null);

    if (!user) {
        return (
            <div className="min-h-screen">
                <NavBar logged={false} />
                <p className="text-center text-gray-600">Please login to see Group Events!</p>
            </div>
        );
    }
    const handleGroupSelect = (groupId: number) => {
        setSelectedGroupId(groupId);
    };
    return (
        <div>
            <NavBar logged={true} />
            <h1 className="text-3xl font-bold mb-6 text-center">Events</h1>
            {selectedGroupId && (
                <div>
                    <h2>Create an Event for Group ID: {selectedGroupId}</h2>
                    <CreateEventForm groupId={selectedGroupId} />
                </div>
            )}
            <GroupList onSelectGroup={handleGroupSelect} actionType="createEvent" />
            <UserProfileEvents />
        </div>
    );
};

export default EventsPage;
