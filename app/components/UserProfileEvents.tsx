import { useEffect, useState } from 'react';
import { Event, GroupMemberReaction } from '@/app/utils/types';
import EventReactions from '@/app/components/EventReactions';
import { useUser } from '@/app/context/UserContext';

const UserProfileEvents = () => {
    const { user } = useUser();
    const [events, setEvents] = useState<Event[] | null>(null);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);

    // Maintain a map where each eventId stores the corresponding members with reactions
    const [membersWithReactions, setMembersWithReactions] = useState<{ [eventId: number]: GroupMemberReaction[] | null }>({});

    useEffect(() => {
        const fetchEvents = async () => {
            if (!user?.ID) return;

            try {
                const res = await fetch(`http://localhost:8080/groups/events?userID=${user.ID}`);
                if (!res.ok) {
                    throw new Error('Failed to fetch events');
                }
                const data: Event[] = await res.json();
                setEvents(data && Array.isArray(data) ? data : []);
            } catch (err) {
                console.error('Error fetching events:', err);
                setError('Failed to load events');
            } finally {
                setLoading(false);
            }
        };

        fetchEvents();
    }, [user]);

    const fetchMembersWithReactions = async (eventID: number, groupID: number) => {
        if (membersWithReactions[eventID]) return; // Avoid duplicate fetch

        try {
            const res = await fetch(`http://localhost:8080/groups/members-with-reactions?eventID=${eventID}&groupID=${groupID}`);
            if (!res.ok) {
                throw new Error('Failed to fetch members with reactions');
            }
            const data: GroupMemberReaction[] = await res.json();
            setMembersWithReactions(prev => ({
                ...prev,
                [eventID]: data // Store members for the specific eventID
            }));
        } catch (err) {
            console.error('Error fetching members with reactions:', err);
            setError('Failed to load members with reactions');
        }
    };

    if (loading) return <p>Loading events...</p>;
    if (error) return <p>{error}</p>;
    if (!events?.length) return <p>No events available</p>;

    return (
        <div>
            <h2>Your Events</h2>
            <ul>
                {events.map(event => (
                    <li key={event.id}>
                        <h3>{event.title}</h3>
                        <p>{event.description}</p>
                        <p>Date: {new Date(event.eventDate).toLocaleString()}</p>

                        {/* Button to fetch and show members with reactions */}
                        <button onClick={() => fetchMembersWithReactions(event.id, event.groupId)}>
                            View Event Members with Reactions
                        </button>

                        {/* Conditionally render members list for this event */}
                        {membersWithReactions[event.id] && (
                            <ul>
                                {membersWithReactions[event.id]?.map(member => (
                                    <li key={member.userID}>
                                        {member.fname} {member.lname} - {member.reaction}
                                    </li>
                                ))}
                            </ul>
                        )}

                        <EventReactions eventId={event.id} />
                    </li>
                ))}
            </ul>
        </div>
    );
};

export default UserProfileEvents;
