package drivers

type DBConnection struct {
	Name     string
	Host     string
	Port     int
	Username string
	Database string
}

type ConnectionList struct {
	Conections []DBConnection
}

func NewConnectionList() *ConnectionList {
	return &ConnectionList{
		Conections: []DBConnection{},
	}
}

func (cl *ConnectionList) AddConnection(conn DBConnection) {
	cl.Conections = append(cl.Conections, conn)
}
