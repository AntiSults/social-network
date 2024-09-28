package sqlite

import "social-network/structs"

func (d *Database) CreateEvent(groupID int, title, description string, eventDate string) error {
	_, err := d.db.Exec(`
        INSERT INTO GroupEvents (GroupID, Title, Description, EventDate)
        VALUES (?, ?, ?, ?)`, groupID, title, description, eventDate)
	return err
}

func (d *Database) ReactToEvent(eventID, userID int, reaction string) error {
	_, err := d.db.Exec(`
        INSERT INTO EventReactions (EventID, UserID, Reaction)
        VALUES (?, ?, ?) ON CONFLICT(EventID, UserID)
        DO UPDATE SET Reaction = ?`, eventID, userID, reaction, reaction)
	return err
}
func (d *Database) GetAllEvents(UserID int) ([]structs.Event, error) {
	var event []structs.Event
	return event, nil
}
