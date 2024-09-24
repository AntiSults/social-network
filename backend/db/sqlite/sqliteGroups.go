package sqlite

import (
	"fmt"
	"social-network/structs"
)

func (d *Database) CreateGroup(name, description string, creator int) error {
	query := "INSERT INTO Groups (Name, Description, CreatorID) VALUES (?, ?, ?)"
	stmt, err := d.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("unable to prepare query: %v", err)
	}
	_, err = stmt.Exec(name, description, creator)
	if err != nil {
		return fmt.Errorf("unable to execute query: %v", err)
	}
	return nil
}
func (d *Database) GetAllGroups() ([]structs.Group, error) {
	rows, err := d.db.Query("SELECT ID, Name, Description, CreatorID FROM Groups")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []structs.Group
	for rows.Next() {
		var group structs.Group
		if err := rows.Scan(&group.ID, &group.Name, &group.Description, &group.CreatorID); err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	return groups, nil
}

func (d *Database) AddUserToGroup(GroupID, UserID int) error {
	query := "INSERT INTO GroupUsers (GroupID, UserID) VALUES (?, ?)"

	_, err := d.db.Exec(query, GroupID, UserID)
	if err != nil {
		return fmt.Errorf("failed to add user to group: %w", err)
	}
	return nil
}
