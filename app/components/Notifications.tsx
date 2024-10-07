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
}

const Notifications: React.FC<Props> = ({ setHasNotifications }) => {
    const [notifications, setNotifications] = useState<Notification | null>(null);

    // Hook to handle WebSocket notifications
    useNotificationWS((user, type) => setNotifications({ user, type }));

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
                        </li>
                        <li>
                            <a>
                                {notifications.user.firstName} {notifications.user.lastName}
                            </a>
                        </li>
                    </ul>
                </div>

            )}
        </>
    );
};

export default Notifications;

