package network

import (
	"errors"
	"fmt"
	"net"
)

type HandleFunc func(*net.TCPConn, []byte, int) error

type Connection struct {
	Conn        *net.TCPConn
	ConnID      uint32
	State       int
	ExitChannel chan int
	HandleAPI   HandleFunc
}

func NewConnection(conn *net.TCPConn, connID uint32, handleFunc HandleFunc) *Connection {
	return &Connection{
		Conn:        conn,
		ConnID:      connID,
		State:       0,
		HandleAPI:   handleFunc,
		ExitChannel: make(chan int, 1),
	}
}

func (c *Connection) ReadMsg() {
	defer fmt.Println("connReader is exit, remote addr is ", c.Conn.RemoteAddr().String())
	defer c.Stop()

	for {
		buf := make([]byte, 1024)
		n, err := c.Conn.Read(buf)
		if err != nil {
			break
		}
		fmt.Println(string(buf[:n]))
		//调用当前链接所绑定的HandleAPI
		if err := c.HandleAPI(c.Conn, buf, n); err != nil {
			fmt.Println("ConnID", c.ConnID, "handle is error", err)
			break
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

func (c *Connection) SendMsg(data []byte) error {
	if c.State == 0 {
		return errors.New("链接失效！")
	}
	if _, err := c.Conn.Write(data); err != nil {
		return errors.New("消息发送失败！")
	}
	return nil
}
