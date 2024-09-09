package drivers

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	Connection *sql.DB
}

// Parses out sqlite:/// from the url
func normalizeSQLiteURL(url string) (string, error) {
	if len(url) >= 10 && url[:10] == "sqlite:///" {
		return url[10:], nil
	}
	return "", errors.New("Invalid SQLite URL format")
}

// Tests database connection
func (db *SQLite) TestConnection(url string) error {
	return db.Connect(url)
}

// Opens a connection to sqlite database
func (db *SQLite) Connect(url string) error {
	normalizedURL, err := normalizeSQLiteURL(url)
	if err != nil {
		return err
	}

	if normalizedURL == "" {
		return errors.New("The database file path cannot be empty.")
	}

	if _, err := os.Stat(normalizedURL); os.IsNotExist(err) {
		return fmt.Errorf("The database file %s does not exist.", url)
	}

	conn, err := sql.Open("sqlite3", normalizedURL)
	if err != nil {
		return err
	}

	// Assign the connection to the SQLite struct
	db.Connection = conn

	// Test Connection
	if err := db.Connection.Ping(); err != nil {
		db.Connection.Close()
		return fmt.Errorf("Failed to ping the database: %w", err)
	}

	return nil
}

// Close connection to sqlite database
func (db *SQLite) CloseConnection() error {
	if db.Connection != nil {
		return db.Connection.Close()
	}

	return nil
}
