package server

import (
	"bufio"
	"bytes"
	"encoding/binary"
)

func Encode(message string) []byte {
	dataLen := len(message)
	m := make([]byte, dataLen+4)

	binary.LittleEndian.PutUint32(m, uint32(dataLen))
	copy(m[4:], message)

	return m
}

func Decode(reader *bufio.Reader) (msg string, err error) {
	var dataLen int32
	lenBytes, err := reader.Peek(4)
	if err != nil {
		return
	}

	lenBuff := bytes.NewBuffer(lenBytes)
	err = binary.Read(lenBuff, binary.LittleEndian, &dataLen)
	if err != nil {
		return
	}

	if int32(reader.Buffered()) < dataLen+4 {
		return
	}

	pack := make([]byte, int(dataLen+4))
	_, err = reader.Read(pack)
	if err != nil {
		return
	}
	return string(pack[4:]), nil
}
