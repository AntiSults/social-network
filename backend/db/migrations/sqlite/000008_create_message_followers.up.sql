CREATE TABLE followers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    follower_id INTEGER,
    status TEXT CHECK( status IN ('pending', 'accepted') ) DEFAULT 'accepted',
    FOREIGN KEY (user_id) REFERENCES Users(ID),
    FOREIGN KEY (follower_id) REFERENCES Users(ID),
    UNIQUE (user_id, follower_id) -- Unique constraint to avoid duplicate records
);
