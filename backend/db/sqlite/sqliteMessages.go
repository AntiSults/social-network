package sqlite

import (
	"errors"
	"fmt"
	"social-network/structs"

	sqlite3 "github.com/mattn/go-sqlite3"
)

// SaveMessage is sasving chat messages into Messages table.
func (d *Database) SaveMessage(message *structs.Message) (*structs.Message, error) {
	// Step 1: Insert the message into the Messages table with sender and recipient
	res, err := d.db.Exec(
		"INSERT INTO Messages (time_created, content, fromuser, toUser) VALUES(?,?,?,?)",
		message.Created, message.Content, message.SenderID, message.RecipientID, // Use toUser for the recipient
	)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	message.ID = int(id)

	return message, nil
}

// FetchMessages is returning all messages sent and received with userID
func (d *Database) FetchMessages(userID int) ([]structs.Message, error) {
	// Fetch messages where the user is either the sender or the recipient
	rows, err := d.db.Query(`
        SELECT m.ID, m.content, m.time_created, m.fromuser, m.toUser
        FROM Messages m
        WHERE m.fromuser = ? OR m.toUser = ?
    `, userID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch messages: %w", err)
	}
	defer rows.Close()

	var messages []structs.Message

	for rows.Next() {
		var message structs.Message

		if err := rows.Scan(&message.ID, &message.Content, &message.Created, &message.SenderID, &message.RecipientID); err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}

		messages = append(messages, message)
	}
	return messages, nil
}

func (d *Database) SaveGroupMessage(message *structs.Message) (*structs.Message, error) {
	// Step 1: Insert the message into the GroupMessages table
	res, err := d.db.Exec(
		"INSERT INTO GroupMessages (content, fromuser, toUserID, groupID, time_created) VALUES(?,?,?,?,?)",
		message.Content, message.SenderID, message.RecipientID, message.GroupID, message.Created,
	)
	if err != nil {
		return nil, err
	}

	// Get the last inserted message ID
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	message.ID = int(id)

	return message, nil
}

func (d *Database) FetchGroupMessages(groupID int, userID int) ([]structs.Message, error) {
	rows, err := d.db.Query(`
        SELECT ID, content, fromuser, toUserID, groupID, time_created
        FROM GroupMessages
        WHERE groupID = ? AND (fromuser = ? OR toUserID = ?)
        ORDER BY time_created ASC
    `, groupID, userID, userID)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch group messages: %w", err)
	}
	defer rows.Close()

	var messages []structs.Message

	for rows.Next() {
		var message structs.Message
		if err := rows.Scan(&message.ID, &message.Content, &message.SenderID, &message.RecipientID, &message.GroupID, &message.Created); err != nil {
			return nil, fmt.Errorf("failed to scan group message: %w", err)
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (d *Database) GetGroupUsers(UserID, GroupID int) ([]int, error) {
	var users []int
	query := `
        SELECT UserID FROM GroupUsers
        WHERE GroupID = ? AND UserID != ?
    `
	rows, err := d.db.Query(query, GroupID, UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		users = append(users, userID)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// func (d *Database) GetGroupsWithMembersByUser(userID int) ([]structs.Group, error) {
// 	var groups []structs.Group

// 	rows, err := d.db.Query(`
//         SELECT
//             g.ID, g.Name, g.Description, g.CreatorID,
//     		IFNULL(GROUP_CONCAT(u.ID), '') AS members
//         FROM
//             Groups g
//         LEFT JOIN
//             GroupUsers gu ON gu.GroupID = g.ID
//         LEFT JOIN
//             Users u ON u.ID = gu.UserID
//         WHERE
//             g.CreatorID = ? OR g.ID IN (
//                 SELECT GroupID FROM GroupUsers WHERE UserID = ?
//             )
//         GROUP BY
//             g.ID;
//     `, userID, userID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var group structs.Group
// 		var memberList string // This will hold the CSV of user IDs
// 		if err := rows.Scan(&group.ID, &group.Name, &group.Description, &group.CreatorID, &memberList); err != nil {
// 			return nil, err
// 		}

// 		// Convert the CSV member list to a slice of integers
// 		group.Members = convertCSVToIntSlice(memberList)
// 		if memberList == "" {
// 			group.Members = []int{}
// 		}
// 		groups = append(groups, group)
// 	}

// 	return groups, nil
// }
