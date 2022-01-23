package server

import (
	"bufio"
	"bytes"
	"encoding/binary"
)

func Encode(message string) ([]byte, error) {
	//var length = int32(len(message))
	//var pack = new(bytes.Buffer)

	// 写入消息头
	//err := binary.Write(pack, binary.LittleEndian, length)
	//if err != nil {
	//	return nil, err
	//}

	// 写入消息实体
	//err := binary.Write(pack, binary.LittleEndian, []byte(message))
	//if err != nil {
	//	return nil, err
	//}

	//return pack.Bytes(), err

	dataLen := len(message)
	m := make([]byte, dataLen + 4)
	binary.LittleEndian.PutUint32(m, uint32(dataLen))
	copy(m[4:], message)
	return m, nil
}

func Decode(reader *bufio.Reader) (msg string, err error) {
	var dataLen int32
	// 读取消息的长度(阻塞)
	lenBytes, err := reader.Peek(4)  // 读取前4个字节的数据
	if err != nil {
		return
	}
	lenBuff := bytes.NewBuffer(lenBytes)
	err = binary.Read(lenBuff, binary.LittleEndian, &dataLen)
	if err != nil {
		return
	}

	// 判断是否是完整的包(Buffered返回缓冲中现有的可读取的字节数)
	if int32(reader.Buffered()) < dataLen+4 {
		return
	}

	// 读取完整包的消息数据
	pack := make([]byte, int(dataLen+4))
	_, err = reader.Read(pack)
	if err != nil {
		return
	}
	return string(pack[4:]), nil
}
