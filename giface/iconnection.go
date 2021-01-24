package giface

import "net"

//define a tcp connection module.
type IConnection interface {
	//start a TCP connect
	Start()

	//stop a connection
	Stop()

	//get TCP socket.
	GetTCPConnection() *net.TCPConn

	//get conn ID
	GetConnID() uint32

	//get remote client tcp address.
	RemoteAddr() net.Addr

	//send data
	SendMsg(data []byte, id uint32) error

	//Add connection property.
	SetProperty(key string, value interface{})

	//Remove connection property.
	RemoveProperty(key string)

	//Get connection property.
	GetProperty(key string) (interface{}, error)
}
