package network

import (
	"errors"
	"sync"
)

type ConnManager struct {
	connections map[int]*Connection
	//保护链接集合的读写锁
	connLock sync.RWMutex
}

// 创建当前链接的方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[int]*Connection),
	}
}

func (cm *ConnManager) AddConnection(conn *Connection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	cm.connections[int(conn.GetConnID())] = conn
}

func (cm *ConnManager) RemoveConnection(conn *Connection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	delete(cm.connections, int(conn.GetConnID()))
}

func (cm *ConnManager) GetConnById(connId uint32) (*Connection, error) {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()
	if cm.connections[int(connId)] == nil {
		return nil, errors.New("没有此链接")
	}
	return cm.connections[int(connId)], nil
}

func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

func (cm *ConnManager) ClearConnMgr() {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	for connId, conn := range cm.connections {
		conn.Stop()
		delete(cm.connections, connId)
	}
}
