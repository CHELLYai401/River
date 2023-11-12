package network

type Message struct {
	Id   uint32
	Len  uint32
	Data []byte
}

func NewMessage(msgId uint32, msgData []byte) *Message {
	return &Message{
		Id:   msgId,
		Len:  uint32(len(msgData)),
		Data: msgData,
	}
}

func (m *Message) GetMsgID() uint32 {
	return m.Id
}

func (m *Message) GetMsgLen() uint32 {
	return m.Len
}

func (m *Message) GetMsgData() []byte {
	return m.Data
}

// 设置消息的ID
func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}

// 设置消息的内容
func (m *Message) SetData(data []byte) {
	m.Data = data
}

// 设置消息的长度
func (m *Message) SetDataLen(len uint32) {
	m.Len = len
}
