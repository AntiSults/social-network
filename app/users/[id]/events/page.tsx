"use client";
import UserProfileEvents from '@/app/components/UserProfileEvents';
import { useUser } from '@/app/context/UserContext';

const EventsPage = () => {
    const { user } = useUser();

    if (!user) {
        return <p>Loading...</p>;
    }

    return (
        <div>
            <h1>Events</h1>
            <UserProfileEvents />
        </div>
    );
};

export default EventsPage;

