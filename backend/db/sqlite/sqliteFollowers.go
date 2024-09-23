package sqlite

import (
	"database/sql"
	"fmt"
	"log"
)

// FollowUser allows a user to follow another user, with a specific status.
func (d *Database) FollowUser(userID int, followerID int, status string) error {
	query := `
    INSERT INTO followers (user_id, follower_id, status)
    VALUES (?, ?, ?)
    ON CONFLICT(user_id, follower_id) DO UPDATE SET status=excluded.status;
    `
	_, err := d.db.Exec(query, userID, followerID, status)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return fmt.Errorf("failed to insert follower: %v", err)
	}
	return nil
}

// UnfollowUser allows a user to unfollow another user by removing the record from the database.
func (d *Database) UnfollowUser(userID int, followerID int) error {
	query := `
    DELETE FROM followers
    WHERE user_id = ? AND follower_id = ?;
    `

	_, err := d.db.Exec(query, userID, followerID)
	if err != nil {
		return fmt.Errorf("failed to unfollow user: %v", err)
	}
	return nil
}

func (d *Database) CheckFollowStatus(userID, followerID int) (bool, bool, error) {
	var status string
	query := `
        SELECT status 
        FROM followers 
        WHERE user_id = ? AND follower_id = ?
    `

	err := d.db.QueryRow(query, userID, followerID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, false, nil // Not following
		}
		return false, false, err // Some other error
	}

	// Return the follow status
	if status == "accepted" {
		return true, false, nil // Following
	} else if status == "pending" {
		return false, true, nil // Follow request pending
	}

	return false, false, nil // Default case
}
