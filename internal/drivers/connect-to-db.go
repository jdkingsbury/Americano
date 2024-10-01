package drivers

import (
	"fmt"
	"net/url"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jdkingsbury/americano/msgtypes"
)

type QueryResultMsg struct {
	Columns []string
	Rows    [][]string
	Error   error
}

type Database interface {
	Connect(url string) error
	CloseConnection() error
	ExecuteQuery(query string) QueryResultMsg
	GetDatabaseName() (string, error)
	GetTables() ([]string, error)
}

func ConnectToDatabase(dbURL string) (Database, tea.Msg) {
	parsedURL, err := url.Parse(dbURL)
	if err != nil {
		return nil, msgtypes.NewErrMsg(fmt.Errorf("Failed to parse URL: %w", err))
	}

	var db Database
	var notification string

	// Determines the database based on URL scheme
	switch strings.ToLower(parsedURL.Scheme) {
	case "postgres", "postgresql":
		notification = "Connecting to PostgreSQL database..."
	case "mysql":
		notification = "Connecting to MySQL database..."
	case "sqlite":
		notification = "Connecting to SQLite database..."
		db = &SQLite{} // Ensure you have your SQLite type defined
	default:
		return nil, msgtypes.NewErrMsg(fmt.Errorf("Unsupported database scheme: %s", parsedURL.Scheme))
	}

	// Calls connect to establish a connection
	err = db.Connect(dbURL)
	if err != nil {
		return nil, msgtypes.NewErrMsg(fmt.Errorf("Failed to connect to the database: %w", err))
	}

	return db, msgtypes.NewNotificationMsg(notification)
}
