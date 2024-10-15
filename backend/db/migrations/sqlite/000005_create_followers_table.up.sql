CREATE TABLE IF NOT EXISTS followers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    follower_id INTEGER,
    status TEXT CHECK( status IN ('pending', 'accepted', 'rejected') ) DEFAULT 'pending',
    FOREIGN KEY (user_id) REFERENCES Users(ID) ON DELETE CASCADE,
    FOREIGN KEY (follower_id) REFERENCES Users(ID) ON DELETE CASCADE,
    UNIQUE (user_id, follower_id) -- Unique constraint to avoid duplicate records
);

CREATE TRIGGER IF NOT EXISTS DeleteRejectedFollowers
AFTER UPDATE ON followers
FOR EACH ROW
WHEN NEW.status = 'rejected'
BEGIN
    DELETE FROM followers WHERE user_id = NEW.user_id AND follower_id = NEW.follower_id;
END;