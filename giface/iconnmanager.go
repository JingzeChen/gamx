package giface

type IConnManager interface{
	//add tcp connection with connId
	Add(conn IConnection, connId uint32)
	//remove connection in ConnManager with connId
	Remove(connId uint32)
	//get connection from ConnManager.
	GetConn(connId uint32) (IConnection, error)
	//clear and release resources of ConnManager.
	ClearConn()
	//get number of connections.
	Len() uint32
}
