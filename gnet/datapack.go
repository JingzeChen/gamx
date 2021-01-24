package gnet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"gamx/giface"
	"gamx/utils"
)

type DataPack struct{}

func (d *DataPack)GetHeaderLen() uint32{
	//msgId 4 bytes(uint32) and data len 4 bytes(uint32).
	return 8
}

func (d *DataPack)Pack(m giface.IMessage) ([]byte, error){
	//Create a data buffer to store binary packed data.
	dataBuffer := bytes.NewBuffer([]byte{})

	//write message length to binary data.
	if err := binary.Write(dataBuffer, binary.LittleEndian, m.GetMsgLen()); err != nil{
		fmt.Println("failed to write data len to binary buffer.")
		return nil, err
	}

	//write message id to binary data.
	if err := binary.Write(dataBuffer, binary.LittleEndian, m.GetMsgId()); err != nil {
		fmt.Println("failed to write message id to binary buffer.")
		return nil, err
	}

	//write message body to binary data.
	if err := binary.Write(dataBuffer, binary.LittleEndian, m.GetData()); err != nil {
		fmt.Println("failed to write message data to binary buffer")
		return nil, err
	}

	return dataBuffer.Bytes(), nil
}

func (d *DataPack)UnPack(binaryBuffer []byte) (giface.IMessage, error){
	//read data from binary Buffer.
	dataBuffer := bytes.NewReader(binaryBuffer)

	msg := &Message{}

	//read message len from binary buffer.
	if err := binary.Read(dataBuffer, binary.LittleEndian, &msg.Len); err != nil{
		fmt.Println("failed to read message len from binary data.")
		return nil, err
	}

	//read message id from binary buffer.
	if err := binary.Read(dataBuffer, binary.LittleEndian, &msg.Id); err != nil {
		fmt.Println("failed to read message id from binary data.")
		return nil, err
	}

	if(utils.GlobalObject.MaxPackageSize > 0 && msg.Len > utils.GlobalObject.MaxPackageSize){
		return nil, errors.New("too Large message size")
	}

	return msg, nil
}

func NewDataPack() (*DataPack){
	return &DataPack{}
}