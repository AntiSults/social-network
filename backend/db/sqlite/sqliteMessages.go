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

func (d *Database) SaveGroupMessage(message *structs.GroupMessage) (*structs.GroupMessage, error) {
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

func (d *Database) FetchGroupMessages(groupID int, userID int) ([]structs.GroupMessage, error) {
	rows, err := d.db.Query(`
        SELECT id, content, fromuser, toUserID, groupID, time_created
        FROM GroupMessages
        WHERE groupID = ? AND (fromuser = ? OR toUserID = ?)
        ORDER BY time_created ASC
    `, groupID, userID, userID)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch group messages: %w", err)
	}
	defer rows.Close()

	var messages []structs.GroupMessage

	for rows.Next() {
		var message structs.GroupMessage
		if err := rows.Scan(&message.ID, &message.Content, &message.SenderID, &message.RecipientID, &message.GroupID, &message.Created); err != nil {
			return nil, fmt.Errorf("failed to scan group message: %w", err)
		}
		messages = append(messages, message)
	}
	return messages, nil
}
