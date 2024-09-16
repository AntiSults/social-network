CREATE TABLE IF NOT EXISTS MessageRecipients(
    message_id INTEGER,
    recipient_id INTEGER,
    PRIMARY KEY (message_id, recipient_id),
    FOREIGN KEY (message_id) REFERENCES Messages(ID) ON DELETE CASCADE,
    FOREIGN KEY (recipient_id) REFERENCES users(ID) ON DELETE CASCADE
);