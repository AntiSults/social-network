CREATE TABLE Posts (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    UserID INTEGER NOT NULL,
    CONTENT TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    privacy TEXT DEFAULT 'public',
    author_first_name TEXT,
    author_last_name TEXT,
    FOREIGN KEY (UserID) REFERENCES users(ID)
);