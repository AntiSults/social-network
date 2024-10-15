CREATE TABLE IF NOT EXISTS Sessions (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    UserID INTEGER,
    SessionToken TEXT NOT NULL,
    ExpiresAt DATETIME NOT NULL,
    FOREIGN KEY(UserID) REFERENCES Users(ID) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS DeleteExpiredSessions
AFTER INSERT ON Sessions
BEGIN
    DELETE FROM Sessions WHERE ExpiresAt <= datetime('now');
END;
