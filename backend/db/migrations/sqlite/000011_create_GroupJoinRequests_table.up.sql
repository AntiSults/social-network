CREATE TABLE IF NOT EXISTS GroupJoinRequests (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    GroupID INTEGER NOT NULL,
    UserID INTEGER NOT NULL,
    InviterID INTEGER,
    Status TEXT CHECK (Status IN ('pending', 'accepted', 'rejected')) DEFAULT 'pending',
    RequestType TEXT DEFAULT 'join', -- 'join' or 'invite'
    FOREIGN KEY (GroupID) REFERENCES Groups(ID),
    FOREIGN KEY (UserID) REFERENCES Users(ID)
);

CREATE TRIGGER IF NOT EXISTS DeleteRejectedGroupJoinRequests
AFTER UPDATE ON GroupJoinRequests
FOR EACH ROW
WHEN NEW.Status = 'rejected'
BEGIN
    DELETE FROM GroupJoinRequests WHERE ID = NEW.ID;
END;

