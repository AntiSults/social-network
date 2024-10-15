package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"social-network/structs"
	"strings"
	"time"
)

func (d *Database) GetUserIdFromToken(session string) (int, error) {
	var userID int
	err := d.db.QueryRow("SELECT UserID FROM Sessions WHERE SessionToken = ?", session).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, ErrNotExists
		}
		return -1, fmt.Errorf("failed to get user ID from token: %w", err)
	}
	return userID, nil
}
func (d *Database) GetAvatarFromID(id int) (string, error) {
	var avatarPath string
	err := d.db.QueryRow("SELECT AvatarPath FROM Users WHERE ID = ?", id).Scan(&avatarPath)
	if err != nil {
		return "", err
	}
	return avatarPath, nil
}

// DeleteSessionFromDB is clearing sessions from DB
func (d *Database) DeleteSessionFromDB(session string) error {
	stmt, err := d.db.Prepare("DELETE FROM Sessions WHERE SessionToken = ?")
	if err != nil {
		return err
	}
	result, err := stmt.Exec(session)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no rows")
	}
	return nil
}

// GetUser is returning single user for user ID
func (d *Database) GetUser(userID int) (*structs.User, error) {
	var user structs.User
	var nickName sql.NullString
	var aboutMe sql.NullString
	var avatarPath sql.NullString

	// Execute the query
	err := d.db.QueryRow(`
		SELECT ID, Email, FirstName, LastName, DOB, NickName, AboutMe, AvatarPath, Profile_visibility
		FROM Users 
		WHERE ID = ?
	`, userID).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.DOB,
		&nickName,
		&aboutMe,
		&avatarPath,
		&user.ProfileVisibility,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with ID %d not found", userID)
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	user.NickName = nickName.String
	user.AboutMe = aboutMe.String
	if avatarPath.Valid && strings.HasPrefix(avatarPath.String, "../public/") {
		user.AvatarPath = strings.TrimPrefix(avatarPath.String, "../public")
	} else {
		user.AvatarPath = avatarPath.String
	}

	return &user, nil
}

// GetUsersByIDs is returning slice of users struct for slice of ID-s
func (d *Database) GetUsersByIDs(userIDs []int) ([]structs.User, error) {

	if len(userIDs) == 0 {
		return nil, fmt.Errorf("no user IDs provided")
	}

	// Prepare a slice of interface{} to hold the IDs
	args := make([]interface{}, len(userIDs))
	for i, id := range userIDs {
		args[i] = id
	}

	// Dynamically build the query with the appropriate number of placeholders
	query := `
		SELECT ID, Email, FirstName, LastName, DOB, NickName, AboutMe, AvatarPath, Profile_visibility
		FROM Users 
		WHERE ID IN (?` + strings.Repeat(",?", len(userIDs)-1) + `)
	`

	// Execute the query
	rows, err := d.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []structs.User
	for rows.Next() {
		var user structs.User
		var nickName, aboutMe, avatarPath sql.NullString

		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.DOB,
			&nickName,
			&aboutMe,
			&avatarPath,
			&user.ProfileVisibility,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}

		user.NickName = nickName.String
		user.AboutMe = aboutMe.String
		if avatarPath.Valid && strings.HasPrefix(avatarPath.String, "../public/") {
			user.AvatarPath = strings.TrimPrefix(avatarPath.String, "../public")
		} else {
			user.AvatarPath = avatarPath.String
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while fetching users: %w", err)
	}

	return users, nil
}

func (d *Database) GetUserByEmail(email string) (*structs.User, error) {
	var user structs.User
	var nickName sql.NullString
	var aboutMe sql.NullString
	var avatarPath sql.NullString

	// Execute the query
	err := d.db.QueryRow(`
		SELECT ID, Email, Password, FirstName, LastName, DOB, NickName, AboutMe, AvatarPath, Profile_visibility 
		FROM Users 
		WHERE Email = ?
	`, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.DOB,
		&nickName,
		&aboutMe,
		&avatarPath,
		&user.ProfileVisibility,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with e-mail %v not found: %w", email, err)
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	user.NickName = nickName.String
	user.AboutMe = aboutMe.String
	if avatarPath.Valid && strings.HasPrefix(avatarPath.String, "../public/") {
		user.AvatarPath = strings.TrimPrefix(avatarPath.String, "../public")
	} else {
		user.AvatarPath = avatarPath.String
	}

	return &user, nil
}

func (d *Database) SaveUser(user structs.User) error {
	if d == nil || d.db == nil {
		fmt.Println("No database")
		return errors.New("database is not initialized")
	}
	prep, err := d.db.Prepare(`
		INSERT INTO Users (Email, Password, FirstName, LastName, DOB, NickName, AboutMe, AvatarPath)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer prep.Close()

	// Prepare sql.NullString fields for nullable columns
	nickName := sql.NullString{String: user.NickName, Valid: user.NickName != ""}
	aboutMe := sql.NullString{String: user.AboutMe, Valid: user.AboutMe != ""}
	avatarPath := sql.NullString{String: user.AvatarPath, Valid: user.AvatarPath != ""}

	_, err = prep.Exec(
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		user.DOB,
		nickName,
		aboutMe,
		avatarPath,
	)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	fmt.Println("Successfully inserted to db!")
	return nil
}

func (d *Database) SaveSession(userID int, token string, exp time.Time) error {
	formattedExp := exp.UTC().Format("2006-01-02 15:04:05")
	_, err := d.db.Exec(`
		INSERT INTO Sessions (UserID, SessionToken, ExpiresAt) VALUES (?, ?, ?)
	`, userID, token, formattedExp)
	if err != nil {
		// Wrap the error with additional context
		return fmt.Errorf("failed to save session for user %d: %w", userID, err)
	}

	return nil
}

// Search function to query SQLite for matching users
func (d *Database) SearchUsersInDB(query string) ([]structs.User, error) {
	// SQL query to search users by name or email
	stmt := `
        SELECT ID, Email, FirstName, LastName, NickName, AboutMe, AvatarPath, DOB, Profile_visibility 
        FROM Users 
        WHERE FirstName LIKE ? OR LastName LIKE ? OR Email LIKE ?`
	// Run the query with search term placeholders
	rows, err := d.db.Query(stmt, "%"+query+"%", "%"+query+"%", "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []structs.User
	for rows.Next() {
		var user structs.User
		var nickname sql.NullString
		var aboutMe sql.NullString
		var avatarPath sql.NullString

		// Scan the additional fields into the User struct
		if err := rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &nickname, &aboutMe, &avatarPath, &user.DOB, &user.ProfileVisibility); err != nil {
			return nil, err
		}

		// Assign scanned fields to the User struct
		user.NickName = nickname.String
		user.AboutMe = aboutMe.String
		if avatarPath.Valid && strings.HasPrefix(avatarPath.String, "../public/") {
			user.AvatarPath = strings.TrimPrefix(avatarPath.String, "../public")
		} else {
			user.AvatarPath = avatarPath.String
		}

		users = append(users, user)
	}

	return users, nil
}

func (d *Database) UpdateProfileVisibility(userID int, visibility string) error {
	query := `UPDATE Users SET Profile_visibility = ? WHERE ID = ?`

	result, err := d.db.Exec(query, visibility, userID)
	if err != nil {
		return fmt.Errorf("failed to update Profile Visibility for user %d: %w", userID, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve the number of affected rows for user %d: %w", userID, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no user found with ID %d", userID)
	}
	return nil
}
