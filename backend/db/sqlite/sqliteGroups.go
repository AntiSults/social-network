package sqlite

import (
	"fmt"
	"log"
	"social-network/structs"
	"strconv"
	"strings"
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

func (d *Database) AddUserToGroup(GroupID, UserID int) error {
	query := "INSERT INTO GroupUsers (GroupID, UserID) VALUES (?, ?)"

	_, err := d.db.Exec(query, GroupID, UserID)
	if err != nil {
		return fmt.Errorf("failed to add user to group: %w", err)
	}
	return nil
}

func (d *Database) GetGroupsWithMembers() ([]structs.Group, error) {
	var groups []structs.Group

	rows, err := d.db.Query(`
        SELECT 
            g.ID, g.Name, g.Description, g.CreatorID, 
    		IFNULL(GROUP_CONCAT(u.ID), '') AS members 
        FROM 
            Groups g 
        LEFT JOIN 
            GroupUsers gu ON gu.GroupID = g.ID
        LEFT JOIN 
            Users u ON u.ID = gu.UserID
        GROUP BY 
            g.ID;
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var group structs.Group
		var memberList string // This will hold the CSV of user IDs
		if err := rows.Scan(&group.ID, &group.Name, &group.Description, &group.CreatorID, &memberList); err != nil {
			return nil, err
		}

		// Convert the CSV member list to a slice of integers
		group.Members = convertCSVToIntSlice(memberList)
		if memberList == "" {
			group.Members = []int{}
		}
		groups = append(groups, group)
	}

	return groups, nil
}

func (d *Database) RequestToJoinGroup(GroupID, UserID int) error {
	_, err := d.db.Exec(`
        INSERT INTO GroupJoinRequests (GroupID, UserID, Status, RequestType)
        VALUES (?, ?, 'pending', 'join')`, GroupID, UserID)
	return err
}

func (d *Database) InviteUserToGroup(GroupID, UserID int) error {
	_, err := d.db.Exec(`
        INSERT INTO GroupJoinRequests (GroupID, UserID, Status, RequestType)
        VALUES (?, ?, 'pending', 'invite')`, GroupID, UserID)
	return err
}

func (d *Database) HandleGroupRequest(GroupID, UserID int, accept bool) error {
	status := "rejected"
	if accept {
		status = "accepted"
	}
	_, err := d.db.Exec(`
        UPDATE GroupJoinRequests
        SET Status = ?
        WHERE GroupID = ? AND UserID = ?`, status, GroupID, UserID)
	if err != nil {
		return err
	}
	if accept {
		return d.AddUserToGroup(GroupID, UserID)
	}
	return nil
}

func (d *Database) GetPendingGroupJoin(groupID int) ([]structs.User, error) {
	// Step 1: Query for all pending follower IDs
	query := `
        SELECT UserID
        FROM GroupJoinRequests
        WHERE GroupID = ? AND status = 'pending'
    `

	rows, err := d.db.Query(query, groupID)
	if err != nil {
		return nil, fmt.Errorf("error querying pending follow requests: %v", err)
	}
	defer rows.Close()

	var userIDs []int

	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			log.Printf("Error scanning followerId: %v", err)
			continue
		}
		userIDs = append(userIDs, userID)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error processing rows: %v", err)
	}

	// Step 2: Use GetUsersByIDs to get user details for the pending followers
	if len(userIDs) == 0 {
		// No pending followers, return an empty slice
		return []structs.User{}, nil
	}

	users, err := d.GetUsersByIDs(userIDs)
	if err != nil {
		return nil, fmt.Errorf("error fetching users for pending follow requests: %v", err)
	}

	return users, nil
}

func convertCSVToIntSlice(csv string) []int {
	if csv == "" {
		return []int{}
	}
	strList := strings.Split(csv, ",")
	var intList []int
	for _, str := range strList {
		num, _ := strconv.Atoi(str)
		intList = append(intList, num)
	}
	return intList
}
