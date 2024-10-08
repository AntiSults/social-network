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
func (d *Database) GetMembersWithReactions(EventID, GroupID int) ([]structs.GroupMemberReaction, error) {

	var members []structs.GroupMemberReaction
	query := `
		SELECT gu.UserID, u.FirstName, u.LastName, COALESCE(er.Reaction, 'pending') AS Reaction
		FROM GroupUsers gu
		JOIN Users u ON gu.UserID = u.ID
		LEFT JOIN EventReactions er ON er.UserID = gu.UserID AND er.EventID = ?
		WHERE gu.GroupID = ?
	`
	rows, err := d.db.Query(query, EventID, GroupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var member structs.GroupMemberReaction
		if err := rows.Scan(&member.UserId, &member.FirstName, &member.LastName, &member.Reaction); err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return members, nil
}

func (d *Database) GetGroupUserIDs(groupID int) ([]int, error) {
	query := "SELECT UserID FROM GroupUsers WHERE GroupID = ?"
	rows, err := d.db.Query(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var userIDs []int
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return userIDs, nil
}
