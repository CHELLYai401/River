package network

type Request struct {
	Conn *Connection
	Data []byte
	Msg  *Message
}

func NewRequest(conn *Connection, msgId uint32, data []byte) *Request {
	return &Request{
		Conn: conn,
		Data: data,
		Msg:  NewMessage(msgId, data),
	}
}

func (rq *Request) GetRequestConn() *Connection {
	return rq.Conn
}

func (rq *Request) GetRequestData() []byte {
	return rq.Data
}
