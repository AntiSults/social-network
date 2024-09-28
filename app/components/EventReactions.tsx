

const EventReactions = ({ eventId, userId }: { eventId: number, userId: number | undefined }) => {
    const handleReaction = async (reaction: string) => {
        await fetch(`http://localhost:8080/groups/events-react`, {
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
