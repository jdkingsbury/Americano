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
	Provider   string
}

// Tests database connection
func (db *SQLite) TestConnection(url string) error {
	return db.Connect(url)
}

func (db *SQLite) Connect(url string) error {
	if url == "" {
		return errors.New("The database file path cannot be empty.")
	}

	if _, err := os.Stat(url); os.IsNotExist(err) {
		return fmt.Errorf("The database file %s does not exist.", url)
	}

	conn, err := sql.Open("sqlite3", url)
	if err != nil {
		return err
	}

	// Assign the connection to the SQLite struct
	db.Connection = conn
  defer func()  {
    if err != nil {
      db.Connection.Close()
    } 
  }()

	// Test Connection
	if err := db.Connection.Ping(); err != nil {
		return fmt.Errorf("Failed to ping the database: %w", err)
	}

	return nil
}

func main() {
	var file string = "activities.db"
	sqlitedb := SQLite{Provider: "sqlite3"}

	if err := sqlitedb.TestConnection(file); err != nil {
		fmt.Printf("Failed to connect to the db: %v\n", err)
	} else {
		fmt.Println("Successfully connected to the database")
	}

  if sqlitedb.Connection != nil {
    defer sqlitedb.Connection.Close()
  }
}
