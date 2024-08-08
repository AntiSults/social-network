package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func ConnectAndMigrateDb(dbPath string, migrationsPath string) (*sql.DB, error) {
    // Open SQLite database connection
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, err
    }

    // Create a new SQLite driver instance
    driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
    if err != nil {
        return nil, err
    }

    // Create a new migration instance
    m, err := migrate.NewWithDatabaseInstance(
        "file://"+migrationsPath,
        "sqlite3",
        driver)
    if err != nil {
        return nil, fmt.Errorf("%w, open", err)
    }

    // Apply all migrations
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return nil, err
    }
	
    return db, nil
}


//TEST FUNCTION REMOVE LATER
func printTables(db *sql.DB) error {
    rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table';")
    if err != nil {
        return fmt.Errorf("failed to query tables: %w", err)
    }
    defer rows.Close()

    fmt.Println("Tables in the database:")
    for rows.Next() {
        var tableName string
        if err := rows.Scan(&tableName); err != nil {
            return fmt.Errorf("failed to scan table name: %w", err)
        }
        fmt.Println(tableName)
    }

    if err := rows.Err(); err != nil {
        return fmt.Errorf("rows iteration error: %w", err)
    }

    return nil
}

