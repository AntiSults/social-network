

const EventReactions = ({ eventId, userId }: { eventId: number, userId: number }) => {
    const handleReaction = async (reaction: string) => {
        await fetch(`/api/events/${eventId}/react`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ userId, reaction }),
        });
    };

    return (
        <div>
            <button onClick={() => handleReaction('going')}>Going</button>
            <button onClick={() => handleReaction('not going')}>Not Going</button>
        </div>
    );
};

export default EventReactions;
