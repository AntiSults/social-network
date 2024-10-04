CREATE TRIGGER IF NOT EXISTS DeleteExpiredEvents
AFTER INSERT ON GroupEvents
BEGIN
    DELETE FROM GroupEvents WHERE EventDate < datetime('now');
END;