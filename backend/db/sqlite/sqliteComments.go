package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"social-network/structs"
)

func (d *Database) SaveComment(comment *structs.Comment) error {
	file := sql.NullString{String: comment.File, Valid: comment.File != ""}

	_, err := d.db.Exec(
		"INSERT INTO Comments (PostID, UserID, Content, File, CreatedAt, author_first_name, author_last_name) VALUES (?, ?, ?, ?, ?, ?, ?)",
		comment.PostID, comment.UserID, comment.Content, file, comment.CreatedAt, comment.AuthorFirstName, comment.AuthorLastName,
	)
	if err != nil {
		return fmt.Errorf("failed to save comment: %w", err)
	}
	return nil
}

func (d *Database) GetComments(postID int) ([]structs.Comment, error) {
	rows, err := d.db.Query("SELECT ID, PostID, UserID, Content, File, CreatedAt, author_first_name, author_last_name FROM Comments WHERE PostID = ? ORDER BY CreatedAt ASC", postID)
	if err != nil {
		log.Printf("Database query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var comments []structs.Comment
	for rows.Next() {
		var comment structs.Comment
		var file sql.NullString

		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &file, &comment.CreatedAt, &comment.AuthorFirstName, &comment.AuthorLastName); err != nil {
			return nil, err
		}

		if file.Valid {
			comment.File = file.String
		} else {
			comment.File = ""
		}

		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
