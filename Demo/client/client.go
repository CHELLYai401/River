package main

import (
	"fmt"
	"io"
	"net"
	"river/network"
)

func main() {
	conn, err := net.Dial("tcp4", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("服务器链接错误")
		return
	}
	fmt.Println("服务器链接成功！")

	data, err := network.Pack(network.NewMessage(2, []byte("88")))
	if err != nil {
		return
	}

	if _, err := conn.Write(data); err != nil {
		fmt.Println("消息发送失败！")
		return
	}

	for {
		GetHeadLen := make([]byte, network.GetHeadLen())
		if _, err := io.ReadFull(conn, GetHeadLen); err != nil {
			return
		}

		msgHead, err := network.Unpack(GetHeadLen)
		if err != nil {
			return
		}

		if msgHead.Len > 0 {
			msg := msgHead
			msg.Data = make([]byte, msg.GetMsgLen())
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("服务器发的消息内容：", string(msg.Data), "服务器发的消息ID：", msg.Id, "服务器发的消息长度：", msg.Len)
		}
	}
}
