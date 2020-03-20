package goipmi

/*
import (
	"encoding/binary"
	"errors"
)

// NetworkFunction identifies the functional class of an IPMI message
type IPMIVersion int

const (
	IPMIv1 = IPMIVersion(1)
	IPMIv2 = IPMIVersion(2)
)

var ErrIPMIVersion = errors.New("ipmi version is unsupported.")

// NetworkFunction identifies the functional class of an IPMI message
type NetworkFunction uint8

// Network Function Codes per section 5.1
const (
	NetworkFunctionChassis = NetworkFunction(0x00)
	NetworkFunctionApp     = NetworkFunction(0x06)
)

// Command fields on an IPMI message
type Command uint8

// Command Number Assignments (table G-1)
const (
	CommandGetDeviceID              = Command(0x01)
	CommandGetAuthCapabilities      = Command(0x38)
	CommandGetSessionChallenge      = Command(0x39)
	CommandActivateSession          = Command(0x3a)
	CommandSetSessionPrivilegeLevel = Command(0x3b)
	CommandCloseSession             = Command(0x3c)
	CommandChassisControl           = Command(0x02)
	CommandChassisStatus            = Command(0x01)
	CommandSetSystemBootOptions     = Command(0x08)
	CommandGetSystemBootOptions     = Command(0x09)
)

const IPMIBodySize = 6

type IPMIBody struct {
	RsAddr     uint8
	NetFnRsLUN uint8 //NetworkFunction
	Checksum   uint8
	RqAddr     uint8
	RqSeq      uint8
	Cmd        Command
}

func (self *IPMIBody) WriteBytes(w *Writer) {
	w.WriteUint8(self.RsAddr)
	w.WriteUint8(self.NetFnRsLUN)
	self.Checksum = checksum(self.RsAddr, self.NetFnRsLUN)
	w.WriteUint8(self.Checksum)
	w.WriteUint8(self.RqAddr)
	w.WriteUint8(self.RqSeq)
	w.WriteUint8(uint8(self.Cmd))
}

func (self *IPMIBody) ReadBytes(r *Reader) {
	if r.Len() < IPMIBodySize {
		r.SetError(ErrInsufficientBytes)
		return
	}

	bs := r.ReadBytes(IPMIBodySize)
	self.RsAddr = uint8(bs[0])
	self.NetFnRsLUN = uint8(bs[1]) //NetworkFunction
	self.Checksum = uint8(bs[2])
	self.RqAddr = uint8(bs[3])
	self.RqSeq = uint8(bs[4])
	self.Cmd = Command(bs[5])
}

// Request structure
type Request struct {
	Body IPMIBody

	Data interface{}
}

func (self *Request) Init(fn NetworkFunction, cmd Command, data interface{}) *Request {
	self.Body.RsAddr = 0x20 // bmcSlaveAddr
	self.Body.NetFnRsLUN = uint8(fn) << 2
	self.Body.RqAddr = 0x81 // remoteSWID
	self.Body.Cmd = cmd
	self.Data = data
	return self
}

func (self *Request) WriteBytes(w *Writer) {
	self.Body.WriteBytes(w)

	old_length := w.Len()
	err := binary.Write(w.ToWriter(),
		binary.LittleEndian, self.Data)
	if nil != err {
		w.SetError(err)
		return
	}

	sum := checksum(self.Body.RqAddr, self.Body.RqSeq, uint8(self.Body.Cmd)) +
		checksumBytes(w.Bytes()[old_length:])
	w.WriteUint8(sum)
}

func (self *Request) ReadBytes(r *Reader) {
	self.Body.ReadBytes(r)
	err := binary.Read(r.ToReader(),
		binary.LittleEndian, self.Data)
	r.SetError(err)
}

// NetFn returns the NetworkFunction portion of the NetFn/RsLUN field
func (req *Request) NetFn() NetworkFunction {
	return NetworkFunction(req.Body.NetFnRsLUN >> 2)
}

func NewRequest(fn NetworkFunction, cmd Command, data interface{}) *Request {
	return &Request{Body: IPMIBody{RsAddr: 0x20, // bmcSlaveAddr
		NetFnRsLUN: uint8(fn) << 2,
		RqAddr:     0x81, // remoteSWID
		Cmd:        cmd},
		Data: data}
}

// Response to an IPMI request must include at least a CompletionCode
	type Response struct {
	Body IPMIBody

	CompletionCode CompletionCode
	Data           interface{}
}

func (self *Response) Init(fn NetworkFunction, cmd Command, data interface{}) *Response {
	self.Body.RsAddr = 0x20 // bmcSlaveAddr
	self.Body.NetFnRsLUN = uint8(fn)
	self.Body.RqAddr = 0x81 // remoteSWID
	self.Body.Cmd = cmd
	self.Data = data
	return self
}

func (self *Response) WriteBytes(w *Writer) {
	self.Body.WriteBytes(w)
	w.WriteUint8(uint8(self.CompletionCode))

	// err := binary.Write(w.ToWriter(),
	// 	binary.LittleEndian, self.Data)
	// w.SetError(err)

	old_length := w.Len()

	if wr, ok := self.Data.(Writable); ok {
		wr.WriteBytes(w)
	} else {
		err := binary.Write(w.ToWriter(),
			binary.LittleEndian, self.Data)
		if nil != err {
			w.SetError(err)
			return
		}
	}

	sum := checksum(self.Body.RqAddr, self.Body.RqSeq, uint8(self.Body.Cmd)) +
		checksumBytes(w.Bytes()[old_length:])
	w.WriteUint8(sum)
}

func (self *Response) ReadBytes(r *Reader) {
	self.Body.ReadBytes(r)
	self.CompletionCode = CompletionCode(r.ReadUint8())

	if rb, ok := self.Data.(Readable); ok {
		rb.ReadBytes(r)
	} else {
		err := binary.Read(r.ToReader(),
			binary.LittleEndian, self.Data)
		r.SetError(err)
	}
}

// NetFn returns the NetworkFunction portion of the NetFn/RsLUN field
func (resp *Response) NetFn() NetworkFunction {
	return NetworkFunction(resp.Body.NetFnRsLUN >> 2)
}

func (resp *Response) Code() CompletionCode {
	return resp.CompletionCode
}

func NewResponse(data interface{}) *Response {
	return &Response{Body: IPMIBody{RsAddr: 0x20, // bmcSlaveAddr
		RqAddr: 0x81}, // remoteSWID
		Data: data}
}

// DeviceIDRequest per section 20.1
type DeviceIDRequest struct{}

// DeviceIDResponse per section 20.1
type DeviceIDResponse struct {
	DeviceID                uint8
	DeviceRevision          uint8
	FirmwareRevision1       uint8
	FirmwareRevision2       uint8
	IPMIVersion             uint8
	AdditionalDeviceSupport uint8
	ManufacturerIDReserved  uint8
	ManufacturerID          OemID
	ProductID               uint16
}

// // DeviceIDRequest per section 20.1
// type DeviceIDRequest struct{}

// // DeviceIDResponse per section 20.1
// type DeviceIDResponse struct {
// 	CompletionCode
// 	DeviceID                uint8
// 	DeviceRevision          uint8
// 	FirmwareRevision1       uint8
// 	FirmwareRevision2       uint8
// 	IPMIVersion             uint8
// 	AdditionalDeviceSupport uint8
// 	ManufacturerID          OemID
// 	ProductID               uint16
// }

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

// PrivLevel
const (
	PrivLevelNone = iota
	PrivLevelCallback
	PrivLevelUser
	PrivLevelOperator
	PrivLevelAdmin
	PrivLevelOEM
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
*/
