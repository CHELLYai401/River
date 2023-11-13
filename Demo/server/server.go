package main

import (
	"fmt"
	"river/network"
)

type Hello struct {
	network.BaseRouter
}

func (h *Hello) PreHandle(requset network.Request) {
	requset.Conn.SendMsg(1, []byte("处理业务之前"))
	fmt.Println("处理业务之前")
}

func (h *Hello) Handle(requset network.Request) {
	requset.Conn.SendMsg(2, []byte("处理业务"))
	fmt.Println("处理业务")
}

func (h *Hello) PostHandle(requset network.Request) {
	requset.Conn.SendMsg(3, []byte("处理业务之后"))
	fmt.Println("处理业务之后")
}

type ByBy struct {
	network.BaseRouter
}

func (b *ByBy) Handle(requset network.Request) {
	requset.Conn.SendMsg(2, []byte("88"))
	fmt.Println("处理业务")
}

func main() {
	s := network.NewServer("test server", "127.0.0.1", 8888)
	s.AddRouter(1, &Hello{})
	s.AddRouter(2, &ByBy{})
	s.Server()
}
