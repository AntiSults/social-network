const checkLoginStatus = (): boolean => {
  if (typeof window !== "undefined") {
    const cookies = document.cookie.split("; ");
    const authCookie = cookies.find((cookie) =>
      cookie.startsWith("session_token=")
    );
    return authCookie !== undefined;
  }
  return false;
};

export default checkLoginStatus;
