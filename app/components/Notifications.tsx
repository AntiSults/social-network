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
    user: User;
    type: string;
    group: string | null
}

const Notifications: React.FC<Props> = ({ setHasNotifications }) => {
    const [notifications, setNotifications] = useState<Notification | null>(null);

    // Hook to handle WebSocket notifications
    useNotificationWS((user, type, group?) => setNotifications({ user, type, group: group || null }));

    useEffect(() => {
        // Notify NavBar about the presence of notifications
        setHasNotifications(!!notifications);
    }, [notifications, setHasNotifications]);

    return (
        <>
            {/* Notifications Dropdown */}
            {notifications && (
                <div className="notifications-dropdown">
                    <hr />
                    <ul>
                        <li key={notifications.user.ID} className="font-bold">
                            <a>You&apos;ve got {notifications.type}</a>
                            {notifications.group && (
                                <a> to join group of {notifications.group}</a>
                            )}
                        </li>
                        <li>
                            <a>
                                from {notifications.user.firstName} {notifications.user.lastName}
                            </a>
                        </li>
                    </ul>
                </div>

            )}
        </>
    );
};

export default Notifications;

