package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/adler32"
)

type Message struct {
	size     int32
	id       int32
	data     []byte
	checksum uint32
}

func NewMessage(msgId int32, data []byte) *Message {
	msg := &Message{
		size:     int32(len(data)) + 4 + 4,
		id:       msgId,
		data:     data,
		checksum: 0,
	}
	msg.checksum = msg.calc()
	return msg
}

func (m *Message) GetData() []byte {
	return m.data
}

func (m *Message) GetId() int32 {
	return m.id
}

func (m *Message) GetSize() int32 {
	return m.size
}

func (m *Message) Checksum() bool {
	return m.checksum == m.calc()
}

func (m *Message) calc() (checksum uint32) {
	if m == nil {
		return
	}

	data := new(bytes.Buffer)
	err := binary.Write(data, binary.LittleEndian, m.id)
	if err != nil {
		return
	}
	err = binary.Write(data, binary.LittleEndian, m.data)
	if err != nil {
		return
	}
	checksum = adler32.Checksum(data.Bytes())
	return
}

func (m *Message) String() string {
	return fmt.Sprintf("Size=%d ID=%d DataLen=%d Checksum=%d", m.GetSize(), m.GetId(), len(m.GetData()), m.checksum)
}
