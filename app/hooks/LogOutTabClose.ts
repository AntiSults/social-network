import { useEffect } from "react";

const useLogoutOnTabClose = (logoutUrl: string) => {
    useEffect(() => {
        let isNavigating = false;

        const handleBeforeNavigate = () => {
            isNavigating = true;
        };
        const handleBeforeUnload = async () => {
            if (!isNavigating) {
                // Only log out if it is not navigating but closing the tab
                await fetch(logoutUrl, {
                    method: "POST",
                    credentials: "include",
                });
            }
        };
        window.addEventListener("beforeunload", handleBeforeUnload);
        window.addEventListener("click", handleBeforeNavigate); // Tracks clicks for navigation
        return () => {
            window.removeEventListener("beforeunload", handleBeforeUnload);
            window.removeEventListener("click", handleBeforeNavigate);
        };
    }, [logoutUrl]);
};

export default useLogoutOnTabClose;
