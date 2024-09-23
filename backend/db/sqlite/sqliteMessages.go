package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"social-network/structs"
	"strconv"
	"strings"

	sqlite3 "github.com/mattn/go-sqlite3"
)

// SaveMessage is sasving chat messages into Messages table and filling MessageRecipients table.
func (d *Database) SaveMessage(message *structs.Message) (*structs.Message, error) {
	// Step 1: Insert the message into the Messages table
	res, err := d.db.Exec(
		"INSERT INTO Messages (time_created, content, fromuser) VALUES(?,?,?)",
		message.Created, message.Content, message.SenderID,
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
	// Get the last inserted message ID
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	message.ID = int(id)
	// Step 2: Insert recipients into MessageRecipients table
	for _, recipientID := range message.RecipientID {
		_, err := d.db.Exec(
			"INSERT INTO MessageRecipients (message_id, recipient_id) VALUES(?,?)",
			message.ID, recipientID,
		)
		if err != nil {
			return nil, err
		}
	}
	return message, nil
}

// FetchMessages is returning all messages sent and received with userID
func (d *Database) FetchMessages(userID int) ([]structs.Message, error) {
	// Fetch messages where the user is a participant
	rows, err := d.db.Query(`
        SELECT m.id, m.content, m.time_created, m.fromuser, GROUP_CONCAT(r.recipient_id) AS recipients
        FROM Messages m
        JOIN MessageParticipants p ON m.id = p.message_id
        LEFT JOIN MessageRecipients r ON m.id = r.message_id
        WHERE p.user_id = ?
        GROUP BY m.id
    `, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch messages: %w", err)
	}
	defer rows.Close()

	var messages []structs.Message

	for rows.Next() {
		var message structs.Message
		var recipientIDs sql.NullString

		if err := rows.Scan(&message.ID, &message.Content, &message.Created, &message.SenderID, &recipientIDs); err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}

		message.RecipientID = []int{}

		// If recipientIDs is not NULL or empty, split and convert to []int
		if recipientIDs.Valid && recipientIDs.String != "" {
			recipientIDStrings := strings.Split(recipientIDs.String, ",")
			for _, idStr := range recipientIDStrings {
				id, err := strconv.Atoi(idStr)
				if err != nil {
					return nil, fmt.Errorf("failed to parse recipient ID: %w", err)
				}
				message.RecipientID = append(message.RecipientID, id)
			}
		}
		messages = append(messages, message)
	}
	return messages, nil
}
