package network

import (
	"errors"
	"fmt"
	"net"
	"river/utils"
)

type Server struct {
	Name        string
	ListenIP    string
	ListenPort  int
	IPVersion   string
	MsgHandler  *MsgHandler
	ConnMgr     *ConnManager
	OnConnStart func(conn *Connection)
	OnConnStop  func(conn *Connection)
}

func NewServer(name string, listenIP string, listenPort int) *Server {
	return &Server{
		Name:       utils.GlobalObject.Name,
		ListenIP:   utils.GlobalObject.IP,
		ListenPort: utils.GlobalObject.Port,
		IPVersion:  utils.GlobalObject.IPVersion,
		MsgHandler: NewMsgHandler(),
		ConnMgr:    NewConnManager(),
	}
}

func CallBackApi(conn *net.TCPConn, buf []byte, cnt int) error {
	if _, err := conn.Write(buf[:cnt]); err != nil {
		fmt.Println("消息回写发生错误！")
		return errors.New("消息回写发生错误！")
	}
	return nil
}

func (s *Server) Start() {
	cnt := 0
	fmt.Println("服务器启动中…………")
	Addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.ListenIP, s.ListenPort))
	if err != nil {
		fmt.Println("地址解析错误…………")
		return
	}
	Listen, err := net.ListenTCP(s.IPVersion, Addr)
	if err != nil {
		fmt.Println("地址监听错误…………")
		return
	}
	fmt.Println("服务器启动成功…………")
	fmt.Println("服务器名字：", s.Name, "服务器监听地址：", s.ListenIP, "服务器监听端口：", s.ListenPort, "服务器版本：", s.IPVersion)

	s.MsgHandler.StartWorkPool()

	go func() {
		for {
			conn, err := Listen.AcceptTCP()
			if err != nil {
				continue
			}
			c := NewConnection(conn, uint32(cnt), s.MsgHandler, s)
			go c.Start()
			cnt += 1
		}
	}()
}

func (s *Server) Stop() {
	s.ConnMgr.ClearConnMgr()
}

func (s *Server) Server() {
	s.Start()
	select {}
}

func (s *Server) AddRouter(msgId int, router IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
}

func (s *Server) GetConnMgr() *ConnManager {
	return s.ConnMgr
}

// 注册OnConnStart 钩子函数的方法
func (s *Server) SetOnConnStart(hookFunc func(connection *Connection)) {
	s.OnConnStart = hookFunc
}

// 注册OnConnStop 钩子函数的方法
func (s *Server) SetOnConnStop(hookFunc func(connection *Connection)) {
	s.OnConnStop = hookFunc
}

// 调用OnConnStart 钩子函数的方法
func (s *Server) CallOnConnStart(connection *Connection) {
	if s.OnConnStart != nil {
		fmt.Println("------Call OnConnStart()------")
		s.OnConnStart(connection)
	}
}

// 调用OnConnStop 钩子函数的方法
func (s *Server) CallOnConnStop(connection *Connection) {
	if s.OnConnStop != nil {
		fmt.Println("------Call OnConnStop()------")
		s.OnConnStop(connection)
	}
}
