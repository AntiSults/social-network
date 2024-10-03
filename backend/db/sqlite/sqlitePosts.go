package sqlite

import (
	"database/sql"
	"fmt"
	"social-network/structs"
)

func (d *Database) SavePost(post *structs.Post) error {
	files := sql.NullString{String: post.Files, Valid: post.Files != ""}

	_, err := d.db.Exec(
		"INSERT INTO posts (userID, content, created_at, privacy, author_first_name, author_last_name, files) VALUES (?, ?, ?, ?, ?, ?, ?)",
		post.UserID, post.Content, post.CreatedAt, post.Privacy, post.AuthorFirstName, post.AuthorLastName, files,
	)

	if err != nil {
		return fmt.Errorf("failed to save post: %w", err)
	}
	return nil
}

func (d *Database) GetPosts(authenticated bool) ([]structs.Post, error) {
	var query string
	if authenticated {
		query = `SELECT ID, UserID, Content, Privacy, created_at, author_first_name, author_last_name, files FROM Posts ORDER BY created_at DESC`
	} else {
		query = `SELECT ID, UserID, Content, Privacy, created_at, author_first_name, author_last_name, files FROM Posts WHERE Privacy = 'public' ORDER BY created_at DESC`
	}

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []structs.Post
	for rows.Next() {
		var post structs.Post
		var files sql.NullString

		if err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.Privacy, &post.CreatedAt, &post.AuthorFirstName, &post.AuthorLastName, &files); err != nil {
			return nil, err
		}

		if files.Valid {
			post.Files = files.String
		} else {
			post.Files = ""
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
