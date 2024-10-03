CREATE TABLE IF NOT EXISTS Groups (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    Name TEXT NOT NULL,
    Description TEXT NOT NULL,
    CreatorID INTEGER REFERENCES Users(ID) ON DELETE CASCADE,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
);