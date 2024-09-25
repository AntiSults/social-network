CREATE TABLE GroupJoinRequests (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    GroupID INTEGER NOT NULL,
    UserID INTEGER NOT NULL,
    Status TEXT DEFAULT 'pending', -- 'pending', 'accepted', 'rejected'
    RequestType TEXT DEFAULT 'join', -- 'join' or 'invite'
    FOREIGN KEY (GroupID) REFERENCES Groups(ID),
    FOREIGN KEY (UserID) REFERENCES Users(ID)
);
