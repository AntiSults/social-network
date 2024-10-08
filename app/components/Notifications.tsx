import { useState, useEffect } from 'react';
import { useNotificationWS } from '@/app/hooks/UseNotify';

interface User {
    ID: number;
    firstName: string;
    lastName: string;
}

interface Props {
    setHasNotifications: (hasNotifications: boolean) => void;
}

interface Notification {
    type: string;
    user: User | null;
    group: string | null
}

const Notifications: React.FC<Props> = ({ setHasNotifications }) => {
    const [notifications, setNotifications] = useState<Notification | null>(null);

    // Hook to handle WebSocket notifications
    useNotificationWS((type, user?, group?) => setNotifications({ type, user: user || null, group: group || null }));

    useEffect(() => {
        // Notify NavBar about the presence of notifications
        setHasNotifications(!!notifications);
    }, [notifications, setHasNotifications]);
    useEffect(() => {
        console.log("Current notification: ", notifications);
    }, [notifications]);
    return (
        <>
            <div className="notifications-dropdown">
                <hr />
                <ul>
                    {/* Check if it's a user-based notification */}
                    {notifications?.user ? (
                        <>
                            <li key={notifications.user.ID} className="font-bold">
                                <a>You&apos;ve got {notifications.type}</a>
                                {/* Display group-related info if applicable */}
                                {notifications.group && (
                                    <a> to join the group {notifications.group}</a>
                                )}
                            </li>
                            <li>
                                <a>
                                    from {notifications.user.firstName} {notifications.user.lastName}
                                </a>
                            </li>
                        </>
                    ) : (
                        // If no user data, display a generic notification
                        <>
                            {notifications?.group && (
                                <li className="font-bold">
                                    <a>Someone from your group &quot;{notifications.group}&quot;</a>
                                    <a> has created {notifications.type}</a>
                                </li>
                            )}
                        </>
                    )}
                </ul>
            </div>
        </>
    );
};

export default Notifications;

