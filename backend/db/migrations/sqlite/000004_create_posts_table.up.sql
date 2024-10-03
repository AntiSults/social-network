CREATE TABLE IF NOT EXISTS Posts (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    UserID INTEGER NOT NULL,
    CONTENT TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    privacy TEXT DEFAULT 'public',
    GroupID INTEGER DEFAULT NULL,
    author_first_name TEXT,
    author_last_name TEXT,
    files TEXT,
    FOREIGN KEY (UserID) REFERENCES users(ID)
);