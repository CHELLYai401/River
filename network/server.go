package network

import (
	"fmt"
	"net"
)

type Server struct {
	Name       string
	ListenIP   string
	ListenPort int
	IPVersion  string
}

func NewServer(name string, listenIP string, listenPort int) *Server {
	return &Server{
		Name:       name,
		ListenIP:   listenIP,
		ListenPort: listenPort,
		IPVersion:  "tcp4",
	}
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
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				continue
			}
			fmt.Println(string(buf[:n]))

			if _, err := conn.Write([]byte("hi!")); err != nil {
				continue
			}
		}
	}()
}

func (s *Server) Stop() {

}

func (s *Server) Server() {
	s.Start()
	select {}
}
