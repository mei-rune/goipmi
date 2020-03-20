package protocol

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/runner-mei/goipmi/protocol/commands"
)

// NetworkFunction identifies the functional class of an IPMI message
type IPMIVersion int

const (
	IPMIv1 = IPMIVersion(1)
	IPMIv2 = IPMIVersion(2)
)

var ErrTimeout = errors.New("timeout")
var ErrIPMIVersion = errors.New("ipmi version is unsupported.")

const IPMIBodySize = 6

type IPMIBody struct {
	RsAddr     uint8
	NetFnRsLUN uint8 //NetworkFunction
	Checksum   uint8
	RqAddr     uint8
	RqSeq      uint8
	Cmd        uint8
}

func (self *IPMIBody) String() string {
	return fmt.Sprintf("RsAddr=%d, NetFnRsLUN=%d,Checksum=%d, RqAddr=%d,RqSeq=%d, Cmd=%d",
		self.RsAddr, self.NetFnRsLUN, self.Checksum, self.RqAddr, self.RqSeq, self.Cmd)
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
	self.Cmd = uint8(bs[5])
}

// Request structure
type Request struct {
	Body IPMIBody

	Data interface{}
}

func (self *Request) String() string {
	return fmt.Sprintf("Request: Body=%s, Data=%s", &self.Body, self.Data)
}

func (self *Request) Init(cmd commands.CommandCode, data interface{}) *Request {
	self.Body.RsAddr = 0x20 // bmcSlaveAddr
	self.Body.NetFnRsLUN = uint8(cmd.NetworkFunction) << 2
	self.Body.RqAddr = 0x81 // remoteSWID
	self.Body.Cmd = cmd.Code
	self.Data = data
	return self
}

func (self *Request) WriteBytes(w *Writer) {
	self.Body.WriteBytes(w)
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

func (self *Request) ReadBytes(r *Reader) {
	self.Body.ReadBytes(r)
	if self.Data == nil {
		return
	}

	if rb, ok := self.Data.(Readable); ok {
		rb.ReadBytes(r.Fork(r.Len() - 1))
	} else {
		err := binary.Read(r.ToReader(),
			binary.LittleEndian, self.Data)
		r.SetError(err)
	}
	r.ReadByte() // read checKsum
}

// NetFn returns the NetworkFunction portion of the NetFn/RsLUN field
func (req *Request) NetFn() commands.NetworkFunction {
	return commands.NetworkFunction(req.Body.NetFnRsLUN >> 2)
}

func NewRequest(cmd commands.CommandCode, data interface{}) *Request {
	return &Request{Body: IPMIBody{RsAddr: 0x20, // bmcSlaveAddr
		NetFnRsLUN: uint8(cmd.NetworkFunction) << 2,
		RqAddr:     0x81, // remoteSWID
		Cmd:        cmd.Code},
		Data: data}
}

// Response to an IPMI request must include at least a CompletionCode
type Response struct {
	Body IPMIBody

	CompletionCode CompletionCode
	Data           interface{}
}

func (self *Response) String() string {
	return fmt.Sprintf("Response: CompletionCode=%s, Body=%s, Data=%s", self.CompletionCode, &self.Body, self.Data)
}

func (self *Response) Init(cmd commands.CommandCode, data interface{}) *Response {
	self.Body.RsAddr = 0x20 // bmcSlaveAddr
	self.Body.NetFnRsLUN = uint8(cmd.NetworkFunction) << 2
	self.Body.RqAddr = 0x81 // remoteSWID
	self.Body.Cmd = cmd.Code
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

	if self.Data == nil {
		return
	}

	if rb, ok := self.Data.(Readable); ok {
		rb.ReadBytes(r.Fork(r.Len() - 1))
	} else {
		err := binary.Read(r.ToReader(),
			binary.LittleEndian, self.Data)
		r.SetError(err)
	}
	r.ReadByte() // read checKsum
}

// NetFn returns the NetworkFunction portion of the NetFn/RsLUN field
func (resp *Response) NetFn() commands.NetworkFunction {
	return commands.NetworkFunction(resp.Body.NetFnRsLUN >> 2)
}

func (resp *Response) Code() CompletionCode {
	return resp.CompletionCode
}

func NewResponse(cmd commands.CommandCode, data interface{}) *Response {
	return &Response{Body: IPMIBody{RsAddr: 0x20, // bmcSlaveAddr
		RqAddr: 0x81}, // remoteSWID
		Data: data}
}
