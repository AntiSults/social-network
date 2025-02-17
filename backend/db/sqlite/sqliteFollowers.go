package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"social-network/structs"
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

func (d *Database) GetPendingFollowRequests(userID int) ([]structs.User, error) {
	// Step 1: Query for all pending follower IDs
	query := `
        SELECT follower_id
        FROM followers
        WHERE user_id = ? AND status = 'pending'
    `
	rows, err := d.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying pending follow requests: %v", err)
	}
	defer rows.Close()

	var followerIDs []int

	for rows.Next() {
		var followerID int
		if err := rows.Scan(&followerID); err != nil {
			log.Printf("Error scanning followerId: %v", err)
			continue
		}
		followerIDs = append(followerIDs, followerID)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error processing rows: %v", err)
	}
	// Step 2: Use GetUsersByIDs to get user details for the pending followers
	if len(followerIDs) == 0 {
		// No pending followers, return an empty slice
		return []structs.User{}, nil
	}
	users, err := d.GetUsersByIDs(followerIDs)
	if err != nil {
		return nil, fmt.Errorf("error fetching users for pending follow requests: %v", err)
	}
	return users, nil
}

func (d *Database) UpdateFollowRequestStatus(followingID, followerID int, status string) error {
	query := `UPDATE followers SET status = ? WHERE user_id = ? AND follower_id = ?`
	_, err := d.db.Exec(query, status, followingID, followerID)
	if err != nil {
		return fmt.Errorf("failed to update follow request status: %v", err)
	}
	return nil
}

func (d *Database) GetFollowersSlice(userID int) ([]int, error) {
	// Single query to get both followers and following, with duplicates automatically removed
	query := `
		SELECT follower_id AS user_id
		FROM followers
		WHERE user_id = ? AND status = 'accepted'
		UNION
		SELECT user_id
		FROM followers
		WHERE follower_id = ? AND status = 'accepted'
		`

	// Execute the query
	rows, err := d.db.Query(query, userID, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying followers and following: %v", err)
	}
	defer rows.Close()

	var userIDs []int

	// Process the result set
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			log.Printf("Error scanning userId: %v", err)
			continue
		}
		userIDs = append(userIDs, userID)
	}

	// Check for errors after processing the rows
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error processing rows: %v", err)
	}

	return userIDs, nil
}

// GetFollowers returns a list of users who follow the given user.
func (d *Database) GetFollowers(userID int) ([]int, error) {
	query := `
		SELECT follower_id
		FROM followers
		WHERE user_id = ? AND status = 'accepted'
	`
	rows, err := d.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying followers: %v", err)
	}
	defer rows.Close()

	var followers []int
	for rows.Next() {
		var followerID int
		if err := rows.Scan(&followerID); err != nil {
			log.Printf("Error scanning followerId: %v", err)
			continue
		}
		followers = append(followers, followerID)
	}
	return followers, nil
}

// GetFollowing returns a list of users that the given user is following.
func (d *Database) GetFollowing(userID int) ([]int, error) {
	query := `
		SELECT user_id
		FROM followers
		WHERE follower_id = ? AND status = 'accepted'
	`
	rows, err := d.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying following: %v", err)
	}
	defer rows.Close()

	var following []int
	for rows.Next() {
		var followingID int
		if err := rows.Scan(&followingID); err != nil {
			log.Printf("Error scanning followingId: %v", err)
			continue
		}
		following = append(following, followingID)
	}
	return following, nil
}
