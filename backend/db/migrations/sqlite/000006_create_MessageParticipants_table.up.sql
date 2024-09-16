CREATE TABLE IF NOT EXISTS MessageParticipants (
    message_id INTEGER,
    user_id INTEGER,
    FOREIGN KEY (message_id) REFERENCES Messages(ID) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(ID) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS InsertSenderToParticipants
AFTER INSERT ON Messages
FOR EACH ROW
BEGIN
    INSERT INTO MessageParticipants (message_id, user_id)
    VALUES (NEW.ID, NEW.fromuser);
END;

CREATE TRIGGER IF NOT EXISTS InsertRecipientsToParticipants
AFTER INSERT ON MessageRecipients
FOR EACH ROW
BEGIN
    INSERT INTO MessageParticipants (message_id, user_id)
    VALUES (NEW.message_id, NEW.recipient_id);
END;
