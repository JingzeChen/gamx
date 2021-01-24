package gnet

import (
	"errors"
	"fmt"
	"gamx/giface"
	"sync"
)

type ConnManager struct{
	connections map[uint32]giface.IConnection
	connLock sync.RWMutex
}

func NewConnManager() *ConnManager{
	return &ConnManager{
		connections: make(map[uint32]giface.IConnection),
	}
}

func (cm *ConnManager)Add(conn giface.IConnection, connId uint32){
	//add writer lock to block connections resources.
	cm.connLock.Lock()
	//unlock.
	defer cm.connLock.Unlock()
	if _, ok := cm.connections[connId]; !ok{
		cm.connections[connId] = conn
		fmt.Println("Succeed to Add connection connID: ", connId)
	} else {
		fmt.Println("connection connId: ", connId, " already exited.")
	}
}

func (cm *ConnManager)Remove(connId uint32){
	//add writer lock to block connections resources.
	cm.connLock.Lock()
	//unlock.
	defer cm.connLock.Unlock()
	if _, ok := cm.connections[connId]; ok {
		delete(cm.connections, connId)
		fmt.Println("Succeed to remove connId: ", connId, " from connection manager.")
	} else {
		fmt.Println("connId: ", connId, " does not exit in connection manager.")
	}
}

func (cm *ConnManager)GetConn(connId uint32) (giface.IConnection, error) {
	//add read lock to block connections resources.
	cm.connLock.RLock()
	//unlock.
	defer cm.connLock.RUnlock()
	if conn, ok := cm.connections[connId]; ok {
		fmt.Println("Succeed to get connId: ", connId)
		return conn, nil
	} else {
		fmt.Println("failed to get connection with id: ", connId)
		return nil, errors.New("failed to get connection")
	}
}

func (cm *ConnManager)ClearConn(){
	//add writer lock to block connections resources.
	cm.connLock.Lock()
	//unlock.
	defer cm.connLock.Unlock()
	for connId, conn := range cm.connections {
		conn.Stop()
		delete(cm.connections, connId)
	}

	fmt.Println("Succeed to clear all connections.")
}

func (cm *ConnManager)Len() uint32 {
	return uint32(len(cm.connections))
}