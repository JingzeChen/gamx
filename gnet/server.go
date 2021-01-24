package gnet

import (
	"errors"
	"gamx/giface"
	"gamx/utils"
	"net"
)
import "fmt"


type Server struct{
	//server's name
	Name string
	//IP's version
	IPVersion string
	//IP's address
	IP string
	//port
	Port uint32
	//Router Handler Module.
	MsgHandler giface.IMsgHandler
	//Server tcp connection manager module.
	ConnMgr giface.IConnManager
	//register client's method before creating connection.
	OnConnStart func(c giface.IConnection)
	//register client's method after closing connection.
	OnConnStop func(c giface.IConnection)
}

func CallbackToCalient(conn *net.TCPConn, buf []byte, cnt int) error {
	if _, err := conn.Write(buf[:cnt]); err != nil {
		fmt.Println("failed to write back data to client err: ", err)
		return errors.New("Callback failed.")
	}
	return nil
}

func (s *Server)Start(){
	//start a server to listen and handle client's request.
	fmt.Println("[START] server name: ", utils.GlobalObject.Name, " MaxPackageSize: ", utils.GlobalObject.MaxPackageSize, " MaxConnSize: ", utils.GlobalObject.MacConnSize)
	fmt.Printf("[START] a server: name:%s, IPVersion:%s, IP:%s, Port:%d\n", s.Name, s.IPVersion, s.IP, s.Port)

	go func(){

		//0.Start a request message work pool.
		s.MsgHandler.StartWorkPool()

		//1. get a TCP connection address.
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("failed to get a TCP address: ", err)
			return
		}

		//2.listen to tcp connection.
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("faied to listen to a TCP address: ", err)
			return
		}
		fmt.Println("Start a tcp Server succ: ", s.Name, " now listening")

		//init conn id with 0.
		var connId uint32 = 0

		//3.block and listen to client's connection, and do something.
		for {
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("faild to connect client: ", err)
				continue
			}

			if s.ConnMgr.Len() > utils.GlobalObject.MacConnSize {
				fmt.Println("current connections number is larger than MaxConnSize.")
				conn.Close()
				continue
			}

			//bind Conn with client's operation and execute.
			dealConn := NewConnection(s, conn, connId, s.MsgHandler)
			connId++
			go dealConn.Start()
		}
	}()
}

func (s *Server)Stop(){
	//stop a server and release server's resources.
	s.ConnMgr.ClearConn()
}

func (s *Server)Serve(){
	//start a server.
	s.Start()

	//TODO do some extra stuff.

	//block Server's state.
	select{}
}

func (s *Server)AddRouter(msgId uint32, router giface.IRouter){
	s.MsgHandler.AddRouter(msgId, router)
	fmt.Println("Add a Router to Server successfully.")
}

func (s *Server)GetConnMgr() giface.IConnManager{
	return s.ConnMgr
}

func (s *Server)SetOnConnStart(hook func(c giface.IConnection)){
	fmt.Println("Succeed to set register Start method before creating connection.")
	s.OnConnStart = hook
}

func (s *Server)SetOnConnStop(hook func(c giface.IConnection)){
	fmt.Println("Succeed to set register Stop method after closing connection.")
	s.OnConnStop = hook
}

func (s *Server)CallOnConnStart(conn giface.IConnection){
	if s.OnConnStart != nil {
		s.OnConnStart(conn)
		fmt.Println("call on connection start register method successfully.")
	} else {
		fmt.Println("on connection start register method is nil.")
	}
}

func (s *Server)CallOnConnStop(conn giface.IConnection){
	if s.OnConnStop != nil {
		s.OnConnStop(conn)
		fmt.Println("call on connection stop register method successfully.")
	} else {
		fmt.Println("on connection stop register method is nil.")
	}
}

//initialize a server
func NewServer() giface.IServer {
	s := &Server{
		Name:utils.GlobalObject.Name,
		IPVersion:"tcp4",
		IP:utils.GlobalObject.Host,
		Port:utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandler(),
		ConnMgr: NewConnManager(),
		OnConnStart: nil,
		OnConnStop: nil,
	}

	return s
}
