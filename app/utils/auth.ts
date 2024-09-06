import { GetServerSidePropsContext } from 'next';
import Cookies from 'js-cookie';

export const serverCookieToken = (req: GetServerSidePropsContext['req']): string | null => {
    if (!req.headers.cookie) return null;

    const tokenCookie = req.headers.cookie.split(';').find(c => c.trim().startsWith('token='));

    if (!tokenCookie) return null;

    const token = tokenCookie.split('=')[1];
    return token || null;
};

export const clientCookieToken = (): string | null => {
    return Cookies.get('session_token') || null;
};
