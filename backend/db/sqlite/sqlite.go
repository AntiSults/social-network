package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"social-network/structs"
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

	if db != nil {
		fmt.Println("Database opened successfully")
	}

	// Create a new SQLite driver instance
	driver, err := sqlitemigration.WithInstance(db.db, &sqlitemigration.Config{})
	if err != nil {
		db.Close() // Close the db if driver creation fails
		return nil, fmt.Errorf("failed to create SQLite driver: %w", err)
	}

	// Create a new migration instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"sqlite3",
		driver,
	)
	if err != nil {
		db.Close() // Close the db if migration instance creation fails
		return nil, fmt.Errorf("failed to create migration instance: %w", err)
	}

	// Apply all migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		db.Close() // Close the db if migration fails
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
func (d *Database) SaveMessage(message *structs.Message) (*structs.Message, error) {

	res, err := d.db.Exec("INSERT INTO Messages (time_created, content, foruser, fromuser) VALUES(?,?,?,?)", message.Created, message.Content, message.RecipientID, message.SenderID)
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

func (d *Database) SavePost(post *structs.Post) error {
	_, err := d.db.Exec(
		"INSERT INTO posts (userID, content, created_at, privacy, author_first_name, author_last_name) VALUES (?, ?, ?, ?, ?, ?)",
		post.UserID, post.Content, post.CreatedAt, post.Privacy, post.AuthorFirstName, post.AuthorLastName,
	)
	if err != nil {
		return fmt.Errorf("failed to save post: %w", err)
	}
	return nil
}

func (d *Database) GetPosts(showAll bool) ([]structs.Post, error) {
	var rows *sql.Rows
	var err error

	if showAll {
		rows, err = d.db.Query("SELECT id, userID, content, created_at, privacy, author_first_name, author_last_name FROM posts ORDER BY created_at DESC")
	} else {
		rows, err = d.db.Query("SELECT id, userID, content, created_at, privacy, author_first_name, author_last_name FROM posts WHERE privacy = 'public' ORDER BY created_at DESC")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch posts: %w", err)
	}
	defer rows.Close()

	var posts []structs.Post
	for rows.Next() {
		var post structs.Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.CreatedAt, &post.Privacy, &post.AuthorFirstName, &post.AuthorLastName); err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to fetch posts: %w", err)
	}

	return posts, nil
}

func (d *Database) FetchMessages(sender int) ([]structs.Message, error) {
	rows, err := d.db.Query(
		`SELECT id, content, strftime('%d-%m-%Y, %H:%M', time_created) AS time_created, foruser, fromuser
	FROM messages 
	WHERE fromuser = ?`, sender)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []structs.Message
	for rows.Next() {
		var message structs.Message
		if err := rows.Scan(&message.ID, &message.Content, &message.Created, &message.RecipientID, &message.SenderID); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}
