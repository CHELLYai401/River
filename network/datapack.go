package network

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// 获取包的头的长度的方法
func GetHeadLen() uint32 {
	//Datalen uint32(4字节) + ID uint32(4字节)
	return 8
}

// 封包方法
func Pack(msg *Message) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return nil, errors.New("消息ID封装失败")
	}

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, errors.New("消息长度封装失败")
	}

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgData()); err != nil {
		return nil, errors.New("消息内容封装失败")
	}

	return dataBuff.Bytes(), nil

}

// 拆包方法
func Unpack(binaryData []byte) (*Message, error) {
	//创建一个从输入二进制的ioReader
	dataBuff := bytes.NewReader(binaryData)
	msg := &Message{}
	//读datalen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	//读msgid
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Len); err != nil {
		return nil, err
	}
	//判断datalen是否已经超出了我们允许的最大包长度
	if 1024 > 0 && msg.Len > 1024 {
		return nil, errors.New("too Large msg data recv")
	}

	return msg, nil
}
