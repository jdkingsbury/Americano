package models

type Database interface {
	Connect(url string) error
	TestConnection(url string) error
	CloseConnection() error
	ExecuteQuery(query string) (columns []string, rows [][]string, err error)
}

