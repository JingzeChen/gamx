package gnet

import "gamx/giface"

type Request struct{
	conn giface.IConnection
	msg  giface.IMessage
}

func (r *Request)GetConnection() giface.IConnection{
	return r.conn
}

func (r *Request)GetData() []byte{
	return r.msg.GetData()
}

func (r *Request)GetMsgId() uint32{
	return r.msg.GetMsgId()
}

func NewRequest(c giface.IConnection, m giface.IMessage) *Request {
	req := &Request{
		conn:c,
		msg: m,
	}
	return req
}
