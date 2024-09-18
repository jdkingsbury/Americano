package drivers

import (
	"fmt"
	"net/url"
	"strings"
)

type DBConnMsg struct {
	Notification string
	Error        error
	DB           Database
}

type Database interface {
	Connect(url string) error
	TestConnection(url string) error
	CloseConnection() error
	ExecuteQuery(query string) (columns []string, rows [][]string, err error)
}

func ConnectToDatabase(dbURL string) (DBConnMsg, error) {
	parsedURL, err := url.Parse(dbURL)
	if err != nil {
		return DBConnMsg{Error: fmt.Errorf("Failed to parse URL: %w", err)}, err
	}

	var db Database
	var notification string

	// Determines the database based on URL scheme
	switch strings.ToLower(parsedURL.Scheme) {
	case "postgres", "postgresql":
		notification = fmt.Sprintf("Connecting to PostgreSQL database...")
	case "mysql":
		notification = fmt.Sprintf("Connecting to MySQL database...")
	case "sqlite":
		notification = fmt.Sprintf("Connecting to Sqlite database...")
		db = &SQLite{}
	default:
		return DBConnMsg{Error: fmt.Errorf("Unsupported database scheme: %s", parsedURL.Scheme)}, nil
	}

	// Calls connect to establish a connection
	err = db.Connect(dbURL)
	if err != nil {
		return DBConnMsg{Error: fmt.Errorf("Failed to connect to the database: %w", err)}, err
	}

	return DBConnMsg{Notification: notification, DB: db}, nil
}
