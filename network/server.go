package network

import (
	"errors"
	"fmt"
	"net"
)

type Server struct {
	Name       string
	ListenIP   string
	ListenPort int
	IPVersion  string
	Router     IRouter
}

func NewServer(name string, listenIP string, listenPort int) *Server {
	return &Server{
		Name:       name,
		ListenIP:   listenIP,
		ListenPort: listenPort,
		IPVersion:  "tcp4",
		Router:     nil,
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

	go func() {
		for {
			conn, err := Listen.AcceptTCP()
			if err != nil {
				continue
			}
			c := NewConnection(conn, 1114, s.Router)
			go c.Start()
		}
	}()
}

func (s *Server) Stop() {

}

func (s *Server) Server() {
	s.Start()
	select {}
}

func (s *Server) AddRouter(router IRouter) {
	s.Router = router
}
