package gnet

type Message struct{
	Id uint32
	Len uint32
	Data []byte
}

func (m *Message)GetMsgId() uint32{
	return m.Id
}

func (m *Message)GetMsgLen() uint32{
	return m.Len
}

func (m *Message)GetData() []byte{
	return m.Data
}

func (m *Message)SetMsgId(id uint32){
	m.Id = id
}

func (m *Message)SetDataLen(len uint32){
	m.Len = len
}

func (m *Message)SetData(data []byte){
	m.Data = data
}

func NewMsgPackage(data []byte, id uint32) *Message {
	return &Message{
		Id: id,
		Len: uint32(len(data)),
		Data: data,
	}
}