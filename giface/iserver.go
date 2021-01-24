package giface

type IServer interface{
	//start a new server.
	Start()
	//stop a server.
	Stop()
	//run a server.
	Serve()
	//add router module.
	AddRouter(msgId uint32, router IRouter)
	//get ConnManager.
	GetConnMgr() IConnManager
	//set client's method before creating connection.
	SetOnConnStart(hook func(c IConnection))
	//Set client's method after closing connection.
	SetOnConnStop(hook func(c IConnection))
	//call client's method before creating connection.
	CallOnConnStart(c IConnection)
	//call client's method after closing connection.
	CallOnConnStop(c IConnection)
}
