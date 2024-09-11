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

// Execute db query
func (db *SQLite) ExecuteQuery(query string) (columns []string, rows [][]string, err error) {
	// Execute the query
	rowsResult, err := db.Connection.Query(query)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rowsResult.Close()

	// Get column names
	columns, err = rowsResult.Columns()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get columns: %w", err)
	}

	// Process rows
	for rowsResult.Next() {
		// Create a slice to hold row values
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		// Scan row values into pointers
		if err := rowsResult.Scan(valuePtrs...); err != nil {
			return nil, nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Convert values to strings
		row := make([]string, len(columns))
		for i, value := range values {
			if value == nil {
				row[i] = "NULL"
			} else {
				row[i] = fmt.Sprintf("%v", value)
			}
		}
		rows = append(rows, row)
	}

	// Check for errors from iterating over rows
	if err := rowsResult.Err(); err != nil {
		return nil, nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return columns, rows, nil
}

// Close connection to sqlite database
func (db *SQLite) CloseConnection() error {
	if db.Connection != nil {
		return db.Connection.Close()
	}

	return nil
}
