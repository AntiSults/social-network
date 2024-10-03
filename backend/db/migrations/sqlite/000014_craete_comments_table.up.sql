CREATE TABLE Comments (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    PostID INTEGER NOT NULL,
    UserID INTEGER NOT NULL,
    Content TEXT NOT NULL,
    File TEXT,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    author_first_name TEXT,
    author_last_name TEXT,
    FOREIGN KEY (PostID) REFERENCES Posts(ID),
    FOREIGN KEY (UserID) REFERENCES Users(ID)
);
