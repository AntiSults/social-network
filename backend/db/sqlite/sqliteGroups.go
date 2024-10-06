package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
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

func (d *Database) InviteUserToGroup(GroupID, UserID, InviterID int) error {
	_, err := d.db.Exec(`
        INSERT INTO GroupJoinRequests (GroupID, UserID, InviterID, Status, RequestType)
        VALUES (?, ?, ?, 'pending', 'invite')`, GroupID, UserID, InviterID)
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

func (d *Database) GetPendingGroupRequests(creatorID int) ([]structs.GroupJoinRequest, error) {
	rows, err := d.db.Query(`
        SELECT 
            g.Name, g.ID AS GroupID, u.ID AS UserID, u.FirstName, u.LastName, r.Status
        FROM 
            GroupJoinRequests r
        JOIN Groups g ON r.GroupID = g.ID
        JOIN Users u ON r.UserID = u.ID
        WHERE 
            g.CreatorID = ? AND r.Status = 'pending'`, creatorID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []structs.GroupJoinRequest
	for rows.Next() {
		var request structs.GroupJoinRequest
		if err := rows.Scan(&request.GroupName, &request.GroupID, &request.UserID, &request.FirstName, &request.LastName, &request.Status); err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}

	return requests, nil
}

func (d *Database) GetPendingGroupInvites(userID int) ([]structs.GroupJoinRequest, error) {
	rows, err := d.db.Query(`
        SELECT 
            g.Name AS GroupName, g.ID AS GroupID, u.ID AS UserID, u.FirstName AS InviterFirstName, u.LastName AS InviterLastName, r.Status
        FROM 
            GroupJoinRequests r
        JOIN Groups g ON r.GroupID = g.ID
        JOIN Users u ON r.InviterID = u.ID -- Get inviter's name
        WHERE 
            r.UserID = ? AND r.Status = 'pending' AND r.RequestType = 'invite'`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invitations []structs.GroupJoinRequest
	for rows.Next() {
		var invite structs.GroupJoinRequest
		if err := rows.Scan(&invite.GroupName, &invite.GroupID, &invite.UserID, &invite.FirstName, &invite.LastName, &invite.Status); err != nil {
			return nil, err
		}
		invitations = append(invitations, invite)
	}

	return invitations, nil
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

func (d *Database) GetGroupsByUser(userID int) ([]structs.Group, error) {
	var groups []structs.Group

	rows, err := d.db.Query(`
        SELECT 
            g.ID, g.Name, g.Description, g.CreatorID,
            IFNULL(GROUP_CONCAT(gu.UserID), '') AS members
        FROM 
            Groups g
        JOIN 
            GroupUsers gu ON gu.GroupID = g.ID
        WHERE 
            gu.UserID = ?
        GROUP BY 
            g.ID;
    `, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var group structs.Group
		var memberList string
		if err := rows.Scan(&group.ID, &group.Name, &group.Description, &group.CreatorID, &memberList); err != nil {
			return nil, err
		}

		group.Members = convertCSVToIntSlice(memberList)
		if memberList == "" {
			group.Members = []int{}
		}
		groups = append(groups, group)
	}

	return groups, nil
}

// GetGroupIDForUser checks if a user is a member of a group and returns the GroupID.
func (d *Database) GetGroupIDForUser(userID int, groupID int) (int, error) {
	var gid int
	err := d.db.QueryRow("SELECT GroupID FROM GroupUsers WHERE UserID = ? AND GroupID = ?", userID, groupID).Scan(&gid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil // Group not found
		}
		return 0, err
	}
	return gid, nil // Returns GroupID
}
