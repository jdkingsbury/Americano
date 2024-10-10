package tests

import (
	"github.com/jdkingsbury/americano/internal/drivers"
)

type MockDatabase struct {
	ExecutedQuery string
	QueryResult   drivers.QueryResultMsg
}

func (m *MockDatabase) Connect(url string) error {
	return nil
}

func (m *MockDatabase) CloseConnection() error {
	return nil
}

func (m *MockDatabase) ExecuteQuery(query string) drivers.QueryResultMsg {
	m.ExecutedQuery = query
	return m.QueryResult
}

func (m *MockDatabase) GetDatabaseName() (string, error) {
	return "mock_db", nil
}

func (m *MockDatabase) GetTables() ([]string, error) {
	return []string{"mock_table"}, nil
}
