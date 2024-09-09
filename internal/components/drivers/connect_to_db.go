package drivers

import (
	"fmt"
	"net/url"
	"strings"
)

type Database interface {
	Connect(url string) error
	TestConnection(url string) error
	CloseConnection() error
}

func ConnectToDatabase(dbURL string) (Database, error) {
	parsedURL, err := url.Parse(dbURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse URL: %w", err)
	}

	var db Database

  // Determines the database based on URL scheme
	switch strings.ToLower(parsedURL.Scheme) {
	case "postgres", "postgresql":
		fmt.Println("Connecting to PostgreSQL database...")
	case "mysql":
		fmt.Println("Connecting to MySQL database...")
	case "sqlite":
		fmt.Println("Connecting to Sqlite database...")
		db = &SQLite{}
	default:
		return nil, fmt.Errorf("Unsupported database scheme: %s", parsedURL.Scheme)
	}

	// Calls connect to establish a connection
	err = db.Connect(dbURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to the database: %w", err)
	}

	return db, nil
}
