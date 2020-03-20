package protocol

import (
	"encoding/binary"
	"fmt"
)

type PayloadType uint8

const (
	PayloadIPMI                PayloadType = 0
	PayloadSOL                 PayloadType = 1
	PayloadOEMExplicit         PayloadType = 2
	PayloadOpenSessionRequest  PayloadType = 0x10
	PayloadOpenSessionResponse PayloadType = 0x11
	PayloadRAKPMessage1        PayloadType = 0x12
	PayloadRAKPMessage2        PayloadType = 0x13
	PayloadRAKPMessage3        PayloadType = 0x14
	PayloadRAKPMessage4        PayloadType = 0x15
)

func (pt PayloadType) Value() PayloadType {
	return pt & 0x3F
}

func (pt PayloadType) Encryption() bool {
	return pt&0x80 != 0
}

func (pt PayloadType) Authenticated() bool {
	return pt&0x40 != 0
}

// AuthCapabilitiesRequest per section 22.13
type AuthCapabilitiesRequest struct {
	ChannelNumber uint8
	PrivLevel     uint8
}

// AuthCapabilitiesResponse per section 22.13
type AuthCapabilitiesResponse struct {
	ChannelNumber   uint8
	AuthTypeSupport uint8
	Status          uint8
	Reserved        uint8
	OemReserved     uint8
	OEMID           OemID
	OEMAux          uint8
}

// AuthType
const (
	AuthTypeNone = iota
	AuthTypeMD2
	AuthTypeMD5
	authTypeReserved
	AuthTypePassword
	AuthTypeOEM
	AuthTypeFormatIPMIV2
)

// SessionChallengeRequest per section 22.16
type SessionChallengeRequest struct {
	AuthType uint8
	Username [16]uint8
}

// SessionChallengeResponse per section 22.16
type SessionChallengeResponse struct {
	TemporarySessionID uint32
	Challenge          [16]byte
}

// ActivateSessionRequest per section 22.17
type ActivateSessionRequest struct {
	AuthType  uint8
	PrivLevel uint8
	AuthCode  [16]uint8
	InSeq     [4]uint8
}

// ActivateSessionResponse per section 22.17
type ActivateSessionResponse struct {
	AuthType   uint8
	SessionID  uint32
	InboundSeq uint32
	MaxPriv    uint8
}

// SessionPrivilegeLevelRequest per section 22.18
type SessionPrivilegeLevelRequest struct {
	PrivLevel uint8
}

// SessionPrivilegeLevelResponse per section 22.18
type SessionPrivilegeLevelResponse struct {
	NewPrivilegeLevel uint8
}

// CloseSessionRequest per section 22.19
type CloseSessionRequest struct {
	SessionID uint32
}

// CloseSessionResponse per section 22.19
type CloseSessionResponse struct {
}

type OpenSessionRequest struct {
	MessageTag uint8
	PrivLevel  uint8
	//Reserved   uint16
	SessionID uint32 // 0xA0A2A3A4

	// Authentication payload (9-16)
	Authentication AuthenticationPayload

	// Integrity payload (17-24)
	Integrity IntegrityPayload

	// Confidentiality payload (25-32)
	Confidentiality ConfidentialityPayload
}

func (self *OpenSessionRequest) String() string {
	return fmt.Sprintf("OpenSessionRequest:	MessageTag=%d, PrivLevel=%d, SessionID=%d, Authentication=%s, Integrity=%s, Confidentiality=%s",
		self.MessageTag, self.PrivLevel, self.SessionID, self.Authentication, self.Integrity, self.Confidentiality)
}

func (self *OpenSessionRequest) WriteBytes(w *Writer) {
	w.WriteUint8(self.MessageTag)
	w.WriteUint8(self.PrivLevel)
	w.WriteUint8(0) // Reserved
	w.WriteUint8(0) // Reserved
	w.WriteUint32(self.SessionID)
	self.Authentication.WriteBytes(w)
	self.Integrity.WriteBytes(w)
	self.Confidentiality.WriteBytes(w)
}

func (self *OpenSessionRequest) ReadBytes(r *Reader) {
	if r.Len() < 8 {
		r.SetError(ErrInsufficientBytes)
		return
	}

	bs := r.ReadBytes(8)
	self.MessageTag = uint8(bs[0])
	self.PrivLevel = uint8(bs[1])
	//self.Reserved = binary.LittleEndian.Uint16(bs[2:])
	self.SessionID = binary.LittleEndian.Uint32(bs[4:])

	self.Authentication.ReadBytes(r)
	self.Integrity.ReadBytes(r)
	self.Confidentiality.ReadBytes(r)
}

type OpenSessionResponse struct {
	MessageTag uint8
	StatusCode uint8
	PrivLevel  uint8
	//Reserved   uint8
	SessionID  uint32
	MSessionID uint32

	// Authentication payload (13-20)
	Authentication AuthenticationPayload

	// Integrity payload (21-28)
	Integrity IntegrityPayload

	// Confidentiality payload (29-36)
	Confidentiality ConfidentialityPayload
}

func (self *OpenSessionResponse) String() string {
	return fmt.Sprintf("OpenSessionResponse:	MessageTag=%d, StatusCode=%d PrivLevel=%d, SessionID=%d, Authentication=%s, Integrity=%s, Confidentiality=%s",
		self.MessageTag, self.StatusCode, self.PrivLevel, self.SessionID, self.Authentication, self.Integrity, self.Confidentiality)
}

func (self *OpenSessionResponse) WriteBytes(w *Writer) {
	w.WriteUint8(self.MessageTag)
	w.WriteUint8(self.StatusCode)
	w.WriteUint8(self.PrivLevel)
	w.WriteUint8(0) // Reserved

	w.WriteUint32(self.SessionID)
	w.WriteUint32(self.MSessionID)
	self.Authentication.WriteBytes(w)
	self.Integrity.WriteBytes(w)
	self.Confidentiality.WriteBytes(w)
}

func (self *OpenSessionResponse) ReadBytes(r *Reader) {
	if r.Len() < 12 {
		r.SetError(ErrInsufficientBytes)
		return
	}

	bs := r.ReadBytes(12)
	self.MessageTag = uint8(bs[0])
	self.StatusCode = uint8(bs[1])
	self.PrivLevel = uint8(bs[2])

	self.SessionID = binary.LittleEndian.Uint32(bs[4:])
	self.MSessionID = binary.LittleEndian.Uint32(bs[8:])

	self.Authentication.ReadBytes(r)
	self.Integrity.ReadBytes(r)
	self.Confidentiality.ReadBytes(r)
}

type RakpMessage1 struct {
	MessageTag uint8
	//Reserved1  [3]uint8
	SessionID uint32 // BMC session ID return by the previous open session response

	Rand      [16]byte // (9-24)
	PrivLevel uint8    // 25
	// Reserved2      uint16
	UserNameLength uint8
	UserName       [16]byte // (29-44)
}

func (self *RakpMessage1) String() string {
	return fmt.Sprintf("RakpMessage1:	MessageTag=%d, SessionID=%d, Rand=%x, PrivLevel=%d, UserNameLength=%d, UserName=%s",
		self.MessageTag, self.SessionID, self.Rand, self.PrivLevel, self.UserNameLength, self.UserName[:self.UserNameLength])
}

func (self *RakpMessage1) WriteBytes(w *Writer) {
	w.WriteUint8(self.MessageTag)
	w.WriteUint8(0) // Reserved1
	w.WriteUint8(0) // Reserved1
	w.WriteUint8(0) // Reserved1

	w.WriteUint32(self.SessionID)

	w.WriteBytes(self.Rand[:])

	w.WriteUint8(self.PrivLevel)
	w.WriteUint16(0) // Reserved2

	w.WriteUint8(self.UserNameLength)
	w.WriteBytes(self.UserName[:self.UserNameLength])
}

func (self *RakpMessage1) ReadBytes(r *Reader) {
	if r.Len() < 28 {
		r.SetError(ErrInsufficientBytes)
		return
	}

	bs := r.ReadBytes(28)
	self.MessageTag = uint8(bs[0])
	self.SessionID = binary.LittleEndian.Uint32(bs[4:])

	for i := 0; i < len(self.Rand); i++ {
		self.Rand[i] = bs[8+i]
	}

	self.PrivLevel = uint8(bs[24])
	self.UserNameLength = uint8(bs[27])

	copy(self.UserName[:], r.ReadBytes(int(self.UserNameLength)))
}

type RakpMessage2 struct {
	MessageTag uint8
	StatusCode uint8
	//Reserved  uint16
	SessionID uint32 // BMC session ID return by the previous open session response

	Rand        [16]byte // (9-24)
	Guid        [16]byte // (25-40)
	KeyExchange []byte   // (41-*) N 取决于RakpMessage1 AuthAlg
}

func (self *RakpMessage2) String() string {
	return fmt.Sprintf("RakpMessage2	MessageTag=%d, StatusCode=%d, SessionID=%d, Rand=%x, Guid=%x, KeyExchange=%x",
		self.MessageTag, self.StatusCode, self.SessionID, self.Rand, self.Guid, self.KeyExchange)
}

func (self *RakpMessage2) WriteBytes(w *Writer) {
	w.WriteUint8(self.MessageTag)
	w.WriteUint8(self.StatusCode)
	w.WriteUint16(0) // Reserved
	w.WriteUint32(self.SessionID)
	w.WriteBytes(self.Rand[:])
	w.WriteBytes(self.Guid[:])

	if nil != self.KeyExchange {
		w.WriteBytes(self.KeyExchange)
	}
}

func (self *RakpMessage2) ReadBytes(r *Reader, keyLength int) {
	if r.Len() < 40 {
		r.SetError(ErrInsufficientBytes)
		return
	}

	bs := r.ReadBytes(40)
	self.MessageTag = uint8(bs[0])
	self.StatusCode = uint8(bs[1])
	self.SessionID = binary.LittleEndian.Uint32(bs[4:])

	for i := 0; i < len(self.Rand); i++ {
		self.Rand[i] = bs[8+i]
	}
	for i := 0; i < len(self.Guid); i++ {
		self.Guid[i] = bs[24+i]
	}

	if keyLength > 0 {
		self.KeyExchange = make([]byte, keyLength)
		copy(self.KeyExchange, r.ReadBytes(keyLength))
	}
}

type RakpMessage3 struct {
	MessageTag uint8
	StatusCode uint8
	//Reserved  uint16
	SessionID uint32 // BMC session ID return by the previous open session response

	KeyExchange []byte // (9-*) N 取决于RakpMessage1 AuthAlg
}

func (self *RakpMessage3) String() string {
	return fmt.Sprintf("RakpMessage3:	MessageTag=%d, StatusCode=%d, SessionID=%d, KeyExchange=%x",
		self.MessageTag, self.StatusCode, self.SessionID, self.KeyExchange)
}

func (self *RakpMessage3) WriteBytes(w *Writer) {
	w.WriteUint8(self.MessageTag)
	w.WriteUint8(self.StatusCode)
	w.WriteUint16(0) // Reserved1
	w.WriteUint32(self.SessionID)

	if nil != self.KeyExchange {
		w.WriteBytes(self.KeyExchange)
	}
}

func (self *RakpMessage3) ReadBytes(r *Reader, keyLength int) {
	if r.Len() < 8 {
		r.SetError(ErrInsufficientBytes)
		return
	}

	bs := r.ReadBytes(8)
	self.MessageTag = uint8(bs[0])
	self.StatusCode = uint8(bs[1])
	self.SessionID = binary.LittleEndian.Uint32(bs[4:])

	if keyLength > 0 {
		self.KeyExchange = make([]byte, keyLength)
		copy(self.KeyExchange, r.ReadBytes(keyLength))
	}
}

type RakpMessage4 struct {
	MessageTag uint8
	StatusCode uint8
	//Reserved  uint16
	SessionID uint32 // BMC session ID return by the previous open session response

	IntegrityCheck []byte // (9-*) N 取决于RakpMessage1 AuthAlg
}

func (self *RakpMessage4) String() string {
	return fmt.Sprintf("RakpMessage4:	MessageTag=%d, StatusCode=%d, SessionID=%d, IntegrityCheck=%x",
		self.MessageTag, self.StatusCode, self.SessionID, self.IntegrityCheck)
}

func (self *RakpMessage4) WriteBytes(w *Writer) {
	w.WriteUint8(self.MessageTag)
	w.WriteUint8(self.StatusCode)
	w.WriteUint16(0) // Reserved
	w.WriteUint32(self.SessionID)

	if nil != self.IntegrityCheck {
		w.WriteBytes(self.IntegrityCheck)
	}
}

func (self *RakpMessage4) ReadBytes(r *Reader, keyLength int) {
	if r.Len() < 8 {
		r.SetError(ErrInsufficientBytes)
		return
	}

	bs := r.ReadBytes(8)
	self.MessageTag = uint8(bs[0])
	self.StatusCode = uint8(bs[1])
	self.SessionID = binary.LittleEndian.Uint32(bs[4:])

	if keyLength > 0 {
		self.IntegrityCheck = make([]byte, keyLength)
		copy(self.IntegrityCheck, r.ReadBytes(keyLength))
	}
}
