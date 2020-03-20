package protocol

import (
	"encoding/binary"
	"errors"

	"io"
)

const (
	//ipmiBodySize    = binary.Size(ipmiBody{})
	//ipmiSessionSize = binary.Size(ipmiSession{})
	ipmiBufSize = 1024
)

var ErrNotIPMIV2 = errors.New("IPMI message version isn't v2.")
var NullAuthCode = make([]byte, 16)
var ErrInvalidAuthCode = errors.New("AuthCode is invalid format.")

type IPMIV1Header struct {
	AuthType  uint8
	Sequence  uint32
	SessionID uint32
	AuthCode  []byte
	Length    uint8
}

func (self *IPMIV1Header) WriteBytes(w *Writer) {
	w.WriteUint8(self.AuthType)
	w.WriteUint32(self.Sequence)
	w.WriteUint32(self.SessionID)
	if self.AuthType != AuthTypeNone {
		if 16 != len(self.AuthCode) {
			if nil == self.AuthCode {
				w.WriteBytes(NullAuthCode)
			} else {
				w.SetError(ErrInvalidAuthCode)
			}
		} else {
			w.WriteBytes(self.AuthCode)
		}
	}
	w.WriteUint8(self.Length)
}

func (self *IPMIV1Header) ReadBytes(r *Reader) {
	self.AuthType = r.ReadUint8()
	self.Sequence = r.ReadUint32()
	self.SessionID = r.ReadUint32()
	if self.AuthType != AuthTypeNone {
		self.AuthCode = r.ReadBytes(16)
	}
	self.Length = r.ReadUint8()
}

type IPMIV2Header struct {
	AuthType     uint8 // 在 v2 中一定为 AuthTypeFormatIPMIV2
	PayloadType  uint8
	OemIana      uint32
	OemPayloadId uint16
	SessionID    uint32
	Sequence     uint32
	Length       uint16
}

func (self *IPMIV2Header) WriteBytes(w *Writer) {
	w.WriteUint8(uint8(AuthTypeFormatIPMIV2))
	w.WriteUint8(self.PayloadType)
	if self.PayloadType == uint8(PayloadOEMExplicit) {
		w.WriteUint32(self.OemIana)
		w.WriteUint16(self.OemPayloadId)
	}
	w.WriteUint32(self.SessionID)
	w.WriteUint32(self.Sequence)
	w.WriteUint16(self.Length)
}

func (self *IPMIV2Header) ReadBytes(r *Reader) {
	self.AuthType = r.ReadUint8()
	if AuthTypeFormatIPMIV2 != self.AuthType {
		r.SetError(ErrNotIPMIV2)
		return
	}
	self.PayloadType = r.ReadUint8()
	if self.PayloadType == uint8(PayloadOEMExplicit) {
		self.OemIana = r.ReadUint32()
		self.OemPayloadId = r.ReadUint16()
	}
	self.SessionID = r.ReadUint32()
	self.Sequence = r.ReadUint32()
	self.Length = r.ReadUint16()
}

func checksum(b ...uint8) uint8 {
	var c uint8
	for _, x := range b {
		c += x
	}
	return -c
}

func checksumBytes(b []byte) uint8 {
	var c uint8
	for _, x := range b {
		c += x
	}
	return -c
}

func binaryWrite(writer io.Writer, data interface{}) {
	err := binary.Write(writer, binary.LittleEndian, data)
	if err != nil {
		// shouldn't happen to a bytes.Buffer
		panic(err)
	}
}
