package drivers

import (
	"fmt"
	"net/url"
	"strings"
)

func ConnectToDatabase(dbURL string) error {
	parsedURL, err := url.Parse(dbURL)
	if err != nil {
		return fmt.Errorf("Failed to parse URL: %w", err)
	}

	switch strings.ToLower(parsedURL.Scheme) {
	case "postgres", "postgresql":
    fmt.Println("Connecting to PostgreSQL database...")
	case "mysql":
    fmt.Println("Connecting to MySQL database...")
	case "sqlite":
    fmt.Println("Connecting to Sqlite database...")
	default:
		return fmt.Errorf("Unsupported database scheme: %s", parsedURL.Scheme)
	}

	return nil
}
