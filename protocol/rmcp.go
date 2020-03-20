package protocol

import (
	"fmt"
)

const (
	rmcpClassASF  = 0x06
	rmcpClassIPMI = 0x07
	rmcpVersion1  = 0x06
)

var (
	rmcpHeaderSize = 4
)

// --------------------+----------------+-----------------------------------------
// Field               | bytes          | Description
// --------------------+----------------+-----------------------------------------
// RMCP Header
// --------------------+----------------+-----------------------------------------
// Version             |   1            | 06h = v1
// --------------------+----------------+-----------------------------------------
// Reserved            |   1            | 00h
// --------------------+----------------+-----------------------------------------
// Sequence Number     |   1            | 序列号
// --------------------+----------------+-----------------------------------------
// Class of Message    |   1            | 这个表示消息的格式
//                     |                | Bit 7   RMCP ACK 域，
//                     |                |         为 0 时为 Normal RMCP 消息
//                     |                |         为 1 时为 RMCP ACK 消息
//                     |                | Bit 6-5 保留
//                     |                | Bit 4-0 消息类型
//                     |                |          为 0-5 时为保留值
//                     |                |          为 6 时为 ASF 消息
//                     |                |          为 7 时为 goipmi. 消息
//                     |                |          为 8 时为 OEM 消息
//                     |                |          其它值为保留值
// --------------------+----------------+-----------------------------------------
// RMCP Data
//    任意数据
// --------------------+----------------+-----------------------------------------
//
//
//  消息类型
// --------------------+----------------+----------------------+----------------------------
//   ACK/Normal bit    |    Class bit   |       消息类型       |   消息数据
// --------------------+----------------+----------------------+----------------------------
//       ACK           |     ASF        |  RCMP ACK 应答消息   |  无数据，仅有消息头
// --------------------+----------------+----------------------+----------------------------
//       ACK           |   all other    |        未定义        |  不合法的
// --------------------+----------------+----------------------+----------------------------
//      normal         |      ASF       |       ASF 消息       | 看 asf 规范
// --------------------+----------------+----------------------+----------------------------
//      normal         |      OEM       |       OEM 消息       | bytes 0-3 = OEM IANA
//                     |                |                      +----------------------------
//                     |                |                      | bytes 4-n = OEM 消息的数据
// --------------------+----------------+----------------------+----------------------------
//      normal         |      goipmi.      |      goipmi. 消息       |  看 goipmi. 规范
// --------------------+----------------+----------------------+----------------------------
type RMCPHeader struct {
	Version            uint8
	Reserved           uint8
	RMCPSequenceNumber uint8
	Class              uint8
}

func (self *RMCPHeader) WriteBytes(w *Writer) {
	w.WriteUint8(self.Version)
	w.WriteUint8(self.Reserved)
	w.WriteUint8(self.RMCPSequenceNumber)
	w.WriteUint8(self.Class)
}

func (self *RMCPHeader) ReadBytes(r *Reader) {
	if r.Len() < rmcpHeaderSize {
		r.SetError(ErrInsufficientBytes)
		return
	}

	bs := r.ReadBytes(rmcpHeaderSize)
	self.Version = uint8(bs[0])
	self.Reserved = uint8(bs[1])
	self.RMCPSequenceNumber = uint8(bs[2])
	self.Class = uint8(bs[3])
}

func (h *RMCPHeader) unsupportedClass() error {
	return fmt.Errorf("unsupported RMCP class: %d", h.Class)
}
