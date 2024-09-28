CREATE TABLE EventReactions (
    EventID INTEGER NOT NULL,
    UserID INTEGER NOT NULL,
    Reaction TEXT CHECK (Reaction IN ('going', 'not going')),
    FOREIGN KEY (EventID) REFERENCES GroupEvents(ID),
    FOREIGN KEY (UserID) REFERENCES Users(ID)
);