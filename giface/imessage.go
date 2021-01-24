package giface

type IMessage interface{

	GetMsgId() uint32

	GetMsgLen() uint32

	GetData() []byte

	SetMsgId(id uint32)

	SetDataLen(len uint32)

	SetData(data []byte)
}
