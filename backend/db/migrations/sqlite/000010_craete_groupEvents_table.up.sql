CREATE TABLE IF NOT EXISTS GroupEvents (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    GroupID INTEGER NOT NULL,
    Title TEXT NOT NULL,
    Description TEXT,
    EventDate DATETIME,
    FOREIGN KEY (GroupID) REFERENCES Groups(ID)
);
