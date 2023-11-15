package network

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

type HandleFunc func(*net.TCPConn, []byte, int) error

type Connection struct {
	Conn        *net.TCPConn
	ConnID      uint32
	State       int
	ExitChannel chan int
	MsgHandler  *MsgHandler
	Ser         *Server
	//无缓冲管道，用于读、写Goroutine之间的消息通信
	msgChan chan []byte
	//链接属性集合
	property map[string]interface{}
	//保护链接属性的锁
	propertyLock sync.RWMutex
}

func NewConnection(conn *net.TCPConn, connID uint32, MsgHandler *MsgHandler, ser *Server) *Connection {
	tcpConn := &Connection{
		Conn:        conn,
		ConnID:      connID,
		State:       0,
		ExitChannel: make(chan int, 1),
		msgChan:     make(chan []byte),
		MsgHandler:  MsgHandler,
		Ser:         ser,
		property:    make(map[string]interface{}),
	}

	tcpConn.Ser.GetConnMgr().AddConnection(tcpConn)
	return tcpConn
}

func (c *Connection) ReadMsg() {
	defer fmt.Println("connReader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {

		GetHeadLen := make([]byte, GetHeadLen())
		if _, err := io.ReadFull(c.Conn, GetHeadLen); err != nil {
			break
		}

		msg, err := Unpack(GetHeadLen)
		if err != nil {
			break
		}

		var data []byte
		if msg.Len > 0 {
			_, err := io.ReadFull(c.Conn, data)
			if err != nil {
				fmt.Println(err)
				break
			}
			msg.SetData(data)
			fmt.Println("客户端发的消息内容：", string(msg.Data), "客户端发的消息ID：", msg.Id, "客户端发的消息长度：", msg.Len)

			c.MsgHandler.SendMsgToTaskQueue(NewRequest(c, msg.Id, msg.Data))

		}

	}
}

func (c *Connection) WriteMsg() {
	fmt.Println("[Writer Gorutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")
	//不断地阻塞等待ch的消息，进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error, ", err)
				return
			}
		case <-c.ExitChannel:
			//代表Reader已经退出，此时Writer也要退出
			return
		}
	}
}

// 获取远程客户端的TCP状态 IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Start() {
	c.State = 1
	go c.ReadMsg()
	go c.WriteMsg()
	c.Ser.CallOnConnStart(c)

}

func (c *Connection) Stop() {
	c.Ser.CallOnConnStop(c)
	if c.State == 0 {
		return
	}
	c.State = 0
	c.Conn.Close()
	c.ExitChannel <- 1
	close(c.ExitChannel)
	close(c.msgChan)
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

	// if _, err := c.Conn.Write(data); err != nil {
	// 	return errors.New("消息发送失败！")
	// }
	//将数据发送给客户端
	c.msgChan <- data
	return nil
}

// 设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

// 获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}

}

// 移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)

}
