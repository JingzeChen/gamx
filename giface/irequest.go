package giface

/*
	Request module
*/

type IRequest interface {
	GetConnection() IConnection
	GetData() []byte
	GetMsgId() uint32
}