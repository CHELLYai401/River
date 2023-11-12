package network

import (
	"errors"
	"fmt"
	"net"
	"river/utils"
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
		Name:       utils.GlobalObject.Name,
		ListenIP:   utils.GlobalObject.IP,
		ListenPort: utils.GlobalObject.Port,
		IPVersion:  utils.GlobalObject.IPVersion,
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
	fmt.Println("服务器名字：", s.Name, "服务器监听地址：", s.ListenIP, "服务器监听端口：", s.ListenPort, "服务器版本：", s.IPVersion)

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
