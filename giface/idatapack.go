package giface

type IDataPack interface {
	GetHeaderLen() uint32
	Pack(message IMessage) ([]byte, error)
	UnPack([]byte) (IMessage, error)
}
