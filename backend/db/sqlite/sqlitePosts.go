package sqlite

import (
	"database/sql"
	"fmt"
	"social-network/structs"
)

func (d *Database) SavePost(post *structs.Post) error {
	files := sql.NullString{String: post.Files, Valid: post.Files != ""}

	_, err := d.db.Exec(
		"INSERT INTO posts (userID, content, created_at, privacy, GroupID, author_first_name, author_last_name, files) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		post.UserID, post.Content, post.CreatedAt, post.Privacy, post.GroupID, post.AuthorFirstName, post.AuthorLastName, files,
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
               g.name AS group_name 
		FROM Posts p 
		LEFT JOIN Groups g ON p.GroupID = g.ID 
		WHERE p.Privacy = 'public' 
		OR (p.Privacy = 'private' AND p.UserID = ?) 
		OR (p.Privacy = 'private' AND p.UserID IN (SELECT follower_id FROM followers WHERE user_id = ? AND status = 'accepted')) 
		OR (p.Privacy = 'group' AND p.GroupID IN (SELECT GroupID FROM GroupUsers WHERE UserID = ?))
		ORDER BY p.created_at DESC
		`
		rows, err = d.db.Query(query, userID, userID, userID)

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

		if err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.Privacy, &post.CreatedAt, &post.GroupID, &post.AuthorFirstName, &post.AuthorLastName, &files, &groupName); err != nil {
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

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
