package sqlite

import "social-network/structs"

func (d *Database) CreateEvent(groupID int, title, description string, eventDate string) error {
	// Use STRFTIME to convert the incoming format when inserting
	_, err := d.db.Exec(`
        INSERT INTO GroupEvents (GroupID, Title, Description, EventDate)
        VALUES (?, ?, ?, STRFTIME('%Y-%m-%d %H:%M:%S', ?) )`, groupID, title, description, eventDate)
	return err
}

func (d *Database) ReactToEvent(eventID, userID int, reaction string) error {
	_, err := d.db.Exec(`
        INSERT INTO EventReactions (EventID, UserID, Reaction)
        VALUES (?, ?, ?)
        ON CONFLICT(EventID, UserID)
        DO UPDATE SET Reaction = ?`, eventID, userID, reaction, reaction)
	return err
}

func (d *Database) GetAllEvents(UserID int) ([]structs.Event, error) {
	var events []structs.Event

	query := `
        SELECT e.ID, e.Title, e.Description, e.EventDate, e.GroupID
        FROM GroupEvents e
        INNER JOIN GroupUsers gu ON gu.GroupID = e.GroupID
        WHERE gu.UserID = ?;
    `

	rows, err := d.db.Query(query, UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event structs.Event
		if err := rows.Scan(&event.ID, &event.Title, &event.Description, &event.EventDate, &event.GroupID); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
