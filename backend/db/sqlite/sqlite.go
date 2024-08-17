package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"social-network/structs"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
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
	driver, err := sqlite3.WithInstance(db.db, &sqlite3.Config{})
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
func (d *Database) GetId_Password(email string) (string, string, error) {

	var userID string
	var hashedPw string
	err := d.db.QueryRow("SELECT ID, Password FROM Users WHERE Email = ?", email).Scan(&userID, &hashedPw)
	if err != nil {
		return "", "", sql.ErrNoRows
	}
	return userID, hashedPw, nil

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
		SELECT Email, FirstName, LastName, DOB, NickName, AboutMe, AvatarPath 
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

func (d *Database) InsertUserToDatabase(user structs.User) error {
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

func (d *Database) SaveSession(userID string, token uuid.UUID, exp time.Time) error {
	_, err := d.db.Exec(`
		INSERT INTO Sessions (UserID, SessionToken, ExpiresAt) VALUES (?, ?, ?)
	`, userID, token.String(), exp)
	if err != nil {
		// Wrap the error with additional context
		return fmt.Errorf("failed to save session for user %s: %w", userID, err)
	}

	return nil
}
