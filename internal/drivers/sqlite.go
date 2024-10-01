package drivers

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jdkingsbury/americano/msgtypes"
	_ "github.com/mattn/go-sqlite3"
)

// TODO: Check to see how we plan to display or log different errors and notifications

type SQLite struct {
	Connection    *sql.DB
	connectionUrl string
}

// Parses out sqlite:/// from the url
func normalizeSQLiteURL(url string) (string, error) {
	if len(url) >= 10 && url[:10] == "sqlite:///" {
		return url[10:], nil
	}
	// return "", msgtypes.NewErrMsg(err, "Invalid SQLite URL format")
	return "", errors.New("Invalid SQLite URL format")
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
	db.connectionUrl = normalizedURL

	// Test Connection
	if err := db.Connection.Ping(); err != nil {
		db.Connection.Close()
		return msgtypes.NewErrMsg(fmt.Errorf("Failed to ping the database: %w", err))
	} else {
		msgtypes.NewNotificationMsg("Successfully connected to database.")
	}

	return nil
}

// Execute db query
func (db *SQLite) ExecuteQuery(query string) QueryResultMsg {
	var columns []string
	var rows [][]string

	// Execute the query
	rowsResult, err := db.Connection.Query(query)
	if err != nil {
		return QueryResultMsg{Error: fmt.Errorf("failed to execute query: %w", err)}
	}
	defer rowsResult.Close()

	// Get column names
	columns, err = rowsResult.Columns()
	if err != nil {
		return QueryResultMsg{Error: fmt.Errorf("failed to get columns: %w", err)}
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
			return QueryResultMsg{Error: fmt.Errorf("failed to scan row: %w", err)}
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
		return QueryResultMsg{Error: fmt.Errorf("error iterating over rows: %w", err)}
	}

	return QueryResultMsg{Columns: columns, Rows: rows, Error: nil}
}

// Close connection to sqlite database
func (db *SQLite) CloseConnection() error {
	if db.Connection != nil {
		return db.Connection.Close()
	}

	return nil
}

func (db *SQLite) GetDatabaseName() (string, error) {
	if db.connectionUrl == "" {
		return "", errors.New("no database connection")
	}

	dbName := filepath.Base(db.connectionUrl)

	return dbName, nil
}

// Fetch table names from SQLite database
func (db *SQLite) GetTables() ([]string, error) {
	rows, err := db.Connection.Query("SELECT name FROM sqlite_master WHERE type='table';")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, table)
	}

	return tables, nil
}
