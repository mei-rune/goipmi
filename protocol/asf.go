package protocol

import (
	"encoding/binary"
)

const (
	asfMessageTypePing = 0x80
	asfMessageTypePong = 0x40
	asfIANA            = 0x000011be
)

var (
	asfHeaderSize = 8
)

// --------------------+----------------+-----------------------------------------
// Field               | bytes          | Description
// --------------------+----------------+-----------------------------------------
// ASF Message Header
// --------------------+----------------+-----------------------------------------
// IANA 企业号         |   4            |
// --------------------+----------------+-----------------------------------------
// 消息类型            |   1            | 如 80h 为 ping 消息
// --------------------+----------------+-----------------------------------------
// 消息 Tag            |   1            | 0 - FEh 由远程控制台生成
//                     |                | FFh     表示消息不是 请求/响应 的一部分
// --------------------+----------------+-----------------------------------------
// 保留                |   1            |
// --------------------+----------------+-----------------------------------------
// 数据长度            |   1            |
// --------------------+----------------+-----------------------------------------
// 数据部分
// --------------------+----------------+-----------------------------------------
type asfHeader struct {
	IANAEnterpriseNumber uint32
	MessageType          uint8
	MessageTag           uint8
	Reserved             uint8
	DataLength           uint8
}

func (self *asfHeader) WriteBytes(w *Writer) {
	w.WriteUint32WithOrder(self.IANAEnterpriseNumber, binary.BigEndian)
	w.WriteUint8(self.MessageType)
	w.WriteUint8(self.MessageTag)
	w.WriteUint8(self.Reserved)
	w.WriteUint8(self.DataLength)
}

func (self *asfHeader) ReadBytes(r *Reader) {
	if r.Len() < asfHeaderSize {
		r.SetError(ErrInsufficientBytes)
		return
	}

	bs := r.ReadBytes(asfHeaderSize)
	self.IANAEnterpriseNumber = binary.BigEndian.Uint32(bs)
	self.MessageType = uint8(bs[4])
	self.MessageTag = uint8(bs[5])
	self.Reserved = uint8(bs[6])
	self.DataLength = uint8(bs[7])
}

type asfMessage struct {
	RMCP RMCPHeader
	ASF  asfHeader
	Data Serialiable
}

func (self *asfMessage) WriteBytes(w *Writer) {
	self.RMCP.WriteBytes(w)
	self.ASF.WriteBytes(w)

	old_len := w.Len()

	if nil != self.Data {
		self.Data.WriteBytes(w)
		new_len := w.Len()
		bs := w.Bytes()
		bs[old_len-1] = uint8(new_len - old_len)
	}
}

func (self *asfMessage) ReadBytes(r *Reader) {
	self.RMCP.ReadBytes(r)
	self.ASF.ReadBytes(r)

	if nil != self.Data {
		self.Data.ReadBytes(r)
	} else if self.ASF.DataLength > 0 {
		//self.Data = r.ReadBytes(self.ASF.DataLength)
		panic("data is nil object.")
	}
}

type asfPong struct {
	IANAEnterpriseNumber  uint32
	OEM                   uint32
	SupportedEntities     uint8
	SupportedInteractions uint8
	Reserved              [6]uint8
}

func (self *asfPong) WriteBytes(w *Writer) {
	w.WriteUint32WithOrder(self.IANAEnterpriseNumber, binary.BigEndian)
	w.WriteUint32WithOrder(self.OEM, binary.BigEndian)
	w.WriteUint8(self.SupportedEntities)
	w.WriteUint8(self.SupportedInteractions)
	for i := 0; i < 6; i++ {
		w.WriteUint8(self.Reserved[i])
	}
}

func (self *asfPong) ReadBytes(r *Reader) {
	if r.Len() < 16 {
		r.SetError(ErrInsufficientBytes)
		return
	}

	bs := r.ReadBytes(16)
	self.IANAEnterpriseNumber = binary.BigEndian.Uint32(bs)
	self.OEM = binary.BigEndian.Uint32(bs[4:])

	self.SupportedEntities = uint8(bs[8])
	self.SupportedInteractions = uint8(bs[9])
	copy(self.Reserved[:], bs[10:])
}

func (m *asfPong) valid() bool {
	return m.SupportedEntities&0x80 != 0
}
