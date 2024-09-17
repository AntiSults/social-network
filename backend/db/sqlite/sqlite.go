package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"social-network/structs"
	"strconv"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	sqlitemigration "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	sqlite3 "github.com/mattn/go-sqlite3"
)

var (
	ErrDuplicate = errors.New("record already exists")
	ErrNotExists = errors.New("row does not exist")
	Db           *Database
)

type Database struct {
	db *sql.DB
}

func ConnectAndMigrateDb(migrationsPath string) (*Database, error) {
	// Open SQLite database connection
	db, err := OpenDatabase()
	if err != nil {
		return nil, err
	}

	fmt.Println("Database opened successfully")

	// Create a new SQLite driver instance
	driver, err := sqlitemigration.WithInstance(db.db, &sqlitemigration.Config{})
	if err != nil {
		db.db.Close() // Close the db if driver creation fails
		return nil, fmt.Errorf("failed to create SQLite driver: %w", err)
	}

	// Create a new migration instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"sqlite3",
		driver,
	)
	if err != nil {
		db.db.Close() // Close the db if migration instance creation fails
		return nil, fmt.Errorf("failed to create migration instance: %w", err)
	}

	// Apply all migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		db.db.Close() // Close the db if migration fails
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	// Assign to global variable
	Db = db

	return db, nil
}

// Close closes the database connection.
func (d *Database) Close() error {
	return d.db.Close()
}

// OpenDatabase opens the database and returns a Database instance.
func OpenDatabase() (*Database, error) {
	db, err := sql.Open("sqlite3", "./db/database.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	return &Database{db: db}, nil
}

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
		SELECT Email, FirstName, LastName, DOB, NickName, AboutMe, AvatarPath, Profile_visibility
		FROM Users 
		WHERE ID = ?
	`, userID).Scan(
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
	user.AvatarPath = avatarPath.String

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
		user.AvatarPath = avatarPath.String

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
	user.AvatarPath = avatarPath.String

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

func (d *Database) SavePost(post *structs.Post) error {
	files := sql.NullString{String: post.Files, Valid: post.Files != ""}

	_, err := d.db.Exec(
		"INSERT INTO posts (userID, content, created_at, privacy, author_first_name, author_last_name, files) VALUES (?, ?, ?, ?, ?, ?, ?)",
		post.UserID, post.Content, post.CreatedAt, post.Privacy, post.AuthorFirstName, post.AuthorLastName, files,
	)

	if err != nil {
		return fmt.Errorf("failed to save post: %w", err)
	}
	fmt.Println("Post successfully inserted to db!")
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
		var recipientIDs sql.NullString // Use sql.NullString to handle NULL values

		if err := rows.Scan(&message.ID, &message.Content, &message.Created, &message.SenderID, &recipientIDs); err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}

		// Initialize the RecipientID slice
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

// func (d *Database) FetchParticipantUserIDs(messageIDs []int) ([]int, error) {
// 	// Prepare query for fetching participant user IDs
// 	query := `
// 		SELECT DISTINCT user_id
// 		FROM MessageParticipants
// 		WHERE message_id IN (?` + strings.Repeat(",?", len(messageIDs)-1) + `)`

// 	// Convert messageIDs to []interface{} for the query arguments
// 	args := make([]interface{}, len(messageIDs))
// 	for i, id := range messageIDs {
// 		args[i] = id
// 	}

// 	rows, err := d.db.Query(query, args...)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to fetch participant user IDs: %w", err)
// 	}
// 	defer rows.Close()

// 	var userIDs []int
// 	for rows.Next() {
// 		var userID int
// 		if err := rows.Scan(&userID); err != nil {
// 			return nil, fmt.Errorf("failed to scan user ID: %w", err)
// 		}
// 		userIDs = append(userIDs, userID)
// 	}
// 	return userIDs, nil
// }

// ChatCommon is returning messages and correspondent recipients users
// func (d *Database) ChatCommon(userID int) (structs.ChatMessage, error) {
// 	// Step 1: Fetch messages where userID is a participant
// 	messages, err := d.FetchMessages(userID)
// 	if err != nil {
// 		return structs.ChatMessage{}, err
// 	}

// 	// Step 2: Extract message IDs from the fetched messages
// 	var messageIDs []int
// 	for _, message := range messages {
// 		messageIDs = append(messageIDs, message.ID)
// 	}

// 	// Step 3: Fetch all unique userIDs from MessageParticipants
// 	userIDs, err := d.FetchParticipantUserIDs(messageIDs)
// 	if err != nil {
// 		return structs.ChatMessage{}, err
// 	}

// 	// Include the current userID (to ensure their details are fetched as well)
// 	userIDs = append(userIDs, userID)

// 	// Step 4: Fetch user details for all users involved in the messages
// 	users, err := d.GetUsersByIDs(userIDs)
// 	if err != nil {
// 		return structs.ChatMessage{}, err
// 	}

// 	// Step 5: Return the combined ChatMessage result
// 	return structs.ChatMessage{
// 		Message: messages,
// 		User:    users,
// 	}, nil
// }
