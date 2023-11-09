package network

type Request struct {
	Conn *Connection
	Data []byte
}

func NewRequest(conn *Connection, data []byte) *Request {
	return &Request{
		Conn: conn,
		Data: make([]byte, 0),
	}
}

func (rq *Request) GetRequestConn() *Connection {
	return rq.Conn
}

func (rq *Request) GetRequestData() []byte {
	return rq.Data
}
