import { useCallback, useState, useEffect } from 'react';
import { useNotificationWS } from '@/app/hooks/UseNotify';
import { useUser } from '@/app/context/UserContext';
import { useRouter } from 'next/navigation';

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

export const Notifications: React.FC<Props> = ({ setHasNotifications }) => {
    const [notifications, setNotifications] = useState<Notification | null>(null);
    const { user: currentUser } = useUser();

    const router = useRouter();
    const handleUserClick = () => {
        const userPage = `/users/${currentUser?.ID}`;

        if (window.location.pathname === userPage) {
            window.location.href = userPage;
        } else {
            router.push(userPage);
        }
    };
    const handleEventClick = () => {
        const userPage = `/users/${currentUser?.ID}/events`;

        if (window.location.pathname === userPage) {
            window.location.href = userPage;
        } else {
            router.push(userPage);
        }
    };
    // Memoize the setNotifications callback
    const handleSetNotifications = useCallback(
        (type: string, user: User | null, group: string | null) => {
            setNotifications({ type, user, group });
        },
        [] // Dependencies: No dependencies so it's memoized once
    );
    useNotificationWS(handleSetNotifications);


    useEffect(() => {
        // Only update if there is an actual change
        setHasNotifications(notifications !== null);
    }, [notifications, setHasNotifications]);
    // useEffect(() => {
    //     console.log("Current notification: ", notifications);
    // }, [notifications]);
    return (
        <>
            <div className="notifications-dropdown">
                <hr />
                <ul>
                    {notifications?.user ? (
                        <>
                            <li key={notifications.user.ID} className="font-bold">
                                <a>You&apos;ve got new {notifications.type}</a>
                                {notifications.group && (
                                    <a> to join the group {notifications.group}</a>
                                )}
                            </li>
                            <li>
                                <a onClick={handleUserClick}>
                                    from {notifications.user.firstName} {notifications.user.lastName}
                                </a>
                            </li>
                        </>
                    ) : (
                        <>
                            {notifications?.group && (
                                <li className="font-bold">
                                    <a>Someone from your group &quot;{notifications.group}&quot;</a>
                                    <a onClick={handleEventClick}>
                                        has created {notifications.type}
                                    </a>
                                </li>
                            )}
                        </>
                    )}
                </ul>
            </div>
        </>
    );
};

