package drivers

type Database interface {
	Connect(url string) error
	TestConnection(url string) error
	CloseConnection() error
}
