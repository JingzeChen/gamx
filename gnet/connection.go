package gnet

import (
	"errors"
	"fmt"
	"gamx/giface"
	"gamx/utils"
	"io"
	"net"
	"sync"
)

type Connection struct{
	//tcp server.
	TcpServer giface.IServer
	//tcp socket
	Conn *net.TCPConn
	//tcp connection id
	ConnID uint32
	//connection status
	isClosed bool
	//channel for communication between reader and writer.
	MsgChan chan []byte
	//channel to inform close or open tcp connection
	ExitChan chan bool
	//Router Handler Module.
	MsgHandle giface.IMsgHandler
	//connection property for client.
	property map[string]interface{}
	//rwlock for connection property.
	propertyLock sync.RWMutex
}

//read goroutine to implement reading process.
func (c *Connection)StartReader(){
	fmt.Println("[Start Reader Goroutine.]")
	defer fmt.Println("exit to read data from ConnId: ", c.ConnID)
	defer c.Stop()

	//define buffer to store data from client with 512bytes size.
	for {
		//read data from client constantly.
		dp := NewDataPack()

		headData := make([]byte, dp.GetHeaderLen())
		_, err := io.ReadFull(c.GetTCPConnection(), headData)
		if err != nil {
			fmt.Println("failed to read head data from bin buffer: ", err)
			break
		}

		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("failed to unpack message head data: ", err)
			break
		}

		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			_, err := io.ReadFull(c.GetTCPConnection(), data)
			if err != nil {
				fmt.Println("failed to read message data: ", err)
				c.ExitChan <- true
				continue
			}
		}
		msg.SetData(data)
		req := Request{
			conn: c,
			msg:  msg,
		}
		fmt.Println("is doing MsgHandler req id:", req.GetMsgId(), " data: ", req.GetData())
		if utils.GlobalObject.WorkPoolSize > 0 {
			c.MsgHandle.SendMsgToTaskQueue(&req)
		} else {
			go c.MsgHandle.DoMsgHandler(&req)
		}
	}
}

//writer goroutine to implement writing process to client.
func (c *Connection)StartWriter(){
	fmt.Println("[Start Writing Goroutine]")

	for {
		select {
			case writerBuf := <-c.MsgChan:
				if _, err := c.Conn.Write(writerBuf); err != nil{
					fmt.Println("failed to write back binary msg to client")
					continue
				}
			case <-c.ExitChan:
				fmt.Println("closing writer goroutine.")
				return
		}
	}
}

func (c *Connection)Start(){
	fmt.Println("start tcp connection id: ", c.ConnID)

	//read data from client.
	go c.StartReader()
	go c.StartWriter()
	//register hook method before creating connection.
	c.TcpServer.CallOnConnStart(c)

	for {
		select{
			case <-c.ExitChan:
				return
		}
	}
}

func (c *Connection)Stop(){
	if c.isClosed == true {
		return
	}

	c.isClosed = true
	c.ExitChan <- true

	c.Conn.Close()
	c.TcpServer.CallOnConnStop(c)
	c.TcpServer.GetConnMgr().Remove(c.ConnID)
	close(c.ExitChan)
	close(c.MsgChan)
}

func (c *Connection)GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection)GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection)RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection)SendMsg(data []byte, msgId uint32) error {
	if c.isClosed == true {
		fmt.Println("connection closed.")
		return errors.New("connection is closed.")
	}

	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(data, msgId))
	if err != nil {
		fmt.Println("failed to pack sent msg to binary msg: ", err)
		return errors.New("pack msg failed.")
	}
	fmt.Println("send message to writer.")
	//pass data to writer.
	c.MsgChan <- binaryMsg
	fmt.Println("succed to send message.")
	return nil
}

func (c *Connection)SetProperty(key string, value interface{}){
	fmt.Println("Succeed to set property for key: ", key)
	c.property[key] = value
}

func (c *Connection)GetProperty(key string) (interface{}, error){
	if p, ok := c.property[key]; ok {
		fmt.Println("Succeed to get property for key: ", key)
		return p, nil
	} else {
		fmt.Println("failed to get property for key: ", key)
		return nil, errors.New("Property not found.")
	}
}

func (c *Connection)RemoveProperty(key string) {
	if _, ok := c.property[key]; ok {
		delete(c.property, key)
		fmt.Println("Succeed to remove property for key: ", key)
	} else {
		fmt.Println("property for key: ", key, " not found.")
	}
}

func NewConnection(s giface.IServer, conn *net.TCPConn, connID uint32, mh giface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer: s,
		Conn:conn,
		ConnID:connID,
		isClosed: false,
		MsgChan: make(chan []byte),
		ExitChan: make(chan bool, 1),
		MsgHandle: mh,
		property: make(map[string]interface{}),
	}

	//add current connection to ConnManager.
	c.TcpServer.GetConnMgr().Add(c, c.ConnID)
	return c
}
