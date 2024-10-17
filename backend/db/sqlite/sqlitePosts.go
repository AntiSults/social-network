package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"social-network/structs"
)

func (d *Database) SavePost(post *structs.Post) error {
	files := sql.NullString{String: post.Files, Valid: post.Files != ""}

	var visibleUsers sql.NullString
	if len(post.VisibleUsers) > 0 {
		jsonUsers, err := json.Marshal(post.VisibleUsers)
		if err != nil {
			return fmt.Errorf("failed to marshal visible_users: %w", err)
		}
		visibleUsers = sql.NullString{String: string(jsonUsers), Valid: true}
	} else {
		visibleUsers = sql.NullString{Valid: false}

	}

	_, err := d.db.Exec(
		"INSERT INTO posts (userID, content, created_at, privacy, GroupID, author_first_name, author_last_name, files, visible_users) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		post.UserID, post.Content, post.CreatedAt, post.Privacy, post.GroupID, post.AuthorFirstName, post.AuthorLastName, files, visibleUsers,
	)

	if err != nil {
		return fmt.Errorf("failed to save post: %w", err)
	}
	return nil
}

func (d *Database) GetPosts(userID int, authenticated bool) ([]structs.Post, error) {
	var query string
	var rows *sql.Rows
	var err error

	if authenticated {
		query = `
		SELECT p.ID, p.UserID, p.Content, p.Privacy, p.created_at, p.GroupID, 
               p.author_first_name, p.author_last_name, p.files, 
               g.name AS group_name, p.visible_users 
		FROM Posts p 
		LEFT JOIN Groups g ON p.GroupID = g.ID 
		WHERE p.Privacy = 'public' 
		OR (p.Privacy = 'private' AND p.UserID = ?) 
		OR (p.Privacy = 'private' AND p.UserID IN (SELECT follower_id FROM followers WHERE user_id = ? AND status = 'accepted')) 
		OR (p.Privacy = 'group' AND p.GroupID IN (SELECT GroupID FROM GroupUsers WHERE UserID = ?))
		OR (p.Privacy = 'almost private' AND (p.UserID = ? OR ? IN (SELECT value FROM json_each(p.visible_users))))
		ORDER BY p.created_at DESC
		`
		rows, err = d.db.Query(query, userID, userID, userID, userID, userID)

	} else {
		query = `
		SELECT p.ID, p.UserID, p.Content, p.Privacy, p.created_at, p.GroupID, 
               p.author_first_name, p.author_last_name, p.files, 
               g.name AS group_name 
		FROM Posts p 
		LEFT JOIN Groups g ON p.GroupID = g.ID 
		WHERE p.Privacy = 'public' 
		ORDER BY p.created_at DESC
		`
		rows, err = d.db.Query(query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []structs.Post
	for rows.Next() {
		var post structs.Post
		var files sql.NullString
		var groupName sql.NullString
		var visibleUsersStr sql.NullString

		if err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.Privacy, &post.CreatedAt, &post.GroupID, &post.AuthorFirstName, &post.AuthorLastName, &files, &groupName, &visibleUsersStr); err != nil {
			return nil, err
		}

		if files.Valid {
			post.Files = files.String
		} else {
			post.Files = ""
		}

		if groupName.Valid {
			post.GroupName = groupName.String
		} else {
			post.GroupName = ""
		}

		if visibleUsersStr.Valid {
			users, err := parseVisibleUsers(visibleUsersStr.String)
			if err != nil {
				return nil, err
			}
			post.VisibleUsers = users
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func parseVisibleUsers(usersStr string) ([]int, error) {
	if usersStr == "" {
		return nil, nil
	}
	var users []int
	err := json.Unmarshal([]byte(usersStr), &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}
