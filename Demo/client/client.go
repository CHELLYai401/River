package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp4", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("服务器链接错误")
		return
	}
	fmt.Println("服务器链接成功！")
	if _, err := conn.Write([]byte("hello!")); err != nil {
		fmt.Println("消息发送失败！")
		return
	}

	for {
		buf := make([]byte, 1024)
		if _, err := conn.Read(buf); err != nil {
			fmt.Println("消息接收失败！")
			return
		}

		fmt.Println(string(buf))
	}
}
