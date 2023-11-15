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

// 创建链接之后执行的钩子函数
func DoConnectionBegin(conn *network.Connection) {
	fmt.Println("---->DoConnectionBegin is Called ...")
	if err := conn.SendMsg(202, []byte("DoConnection Begin")); err != nil {
		fmt.Println(err)
	}
	//给当前链接设置属性
	fmt.Println("Set conn Name , Hoe……")
	conn.SetProperty("Name", "chelly")
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Name = ", name)
	}
	conn.SetProperty("Age", "24")

}

// 销毁链接之后执行的钩子函数
func DoConnectionLost(conn *network.Connection) {
	fmt.Println("---->DoConnectionLost is Called ...")
	fmt.Println("conn ID = ", conn.GetConnID(), " is Lost...")
	//TODO  回收资源，通知其他玩家
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Name = ", name)
	}
}

func main() {
	s := network.NewServer("test server", "127.0.0.1", 8888)
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)
	s.AddRouter(1, &Hello{})
	s.AddRouter(2, &ByBy{})
	s.Server()
}
