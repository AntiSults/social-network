CREATE TABLE IF NOT EXISTS MessageParticipants (
    message_id INTEGER,
    user_id INTEGER,
    FOREIGN KEY (message_id) REFERENCES Messages(ID) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(ID) ON DELETE CASCADE
);


