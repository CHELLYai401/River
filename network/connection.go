package network

import (
	"errors"
	"fmt"
	"io"
	"net"
)

type HandleFunc func(*net.TCPConn, []byte, int) error

type Connection struct {
	Conn        *net.TCPConn
	ConnID      uint32
	State       int
	ExitChannel chan int
	Router      IRouter
}

func NewConnection(conn *net.TCPConn, connID uint32, router IRouter) *Connection {
	return &Connection{
		Conn:        conn,
		ConnID:      connID,
		State:       0,
		ExitChannel: make(chan int, 1),
		Router:      router,
	}
}

func (c *Connection) ReadMsg() {
	defer fmt.Println("connReader is exit, remote addr is ", c.Conn.RemoteAddr().String())
	defer c.Stop()

	for {

		GetHeadLen := make([]byte, GetHeadLen())
		if _, err := io.ReadFull(c.Conn, GetHeadLen); err != nil {
			return
		}

		msgHead, err := Unpack(GetHeadLen)
		if err != nil {
			return
		}

		if msgHead.Len > 0 {
			msg := msgHead
			msg.Data = make([]byte, msg.GetMsgLen())
			_, err := io.ReadFull(c.Conn, msg.Data)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("客户端发的消息内容：", string(msg.Data), "客户端发的消息ID：", msg.Id, "客户端发的消息长度：", msg.Len)

			r := NewRequest(c, 1, msg.Data)

			if c.Router != nil {
				c.Router.PreHandle(r)
				c.Router.Handle(r)
				c.Router.PostHandle(r)
			}
		}

	}
}

func (c *Connection) Start() {
	c.State = 1
	go c.ReadMsg()
}

func (c *Connection) Stop() {
	if c.State == 0 {
		return
	}
	c.State = 0
	c.Conn.Close()
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) GetConnAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SendMsg(msgid uint32, data []byte) error {
	if c.State == 0 {
		return errors.New("链接失效！")
	}
	data, err := Pack(NewMessage(msgid, data))
	if err != nil {
		return err
	}

	if _, err := c.Conn.Write(data); err != nil {
		return errors.New("消息发送失败！")
	}
	return nil
}
