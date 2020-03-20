package goipmi

import (
	"fmt"

	"github.com/runner-mei/goipmi/protocol"
)

// section 31.2
type GetSELInfoRequest struct {
}

type GetSELInfoResponse struct {
	// CompletionCode
	Version   uint8  // minutes
	Enitities uint16 // LS Byte first

	FreeSpace          uint16 // LS Byte first
	RecentAddTimestamp uint32 // LS Byte first
	RecentDelTimestamp uint32 // LS Byte first
	OperationSupport   uint8
}

func (self GetSELInfoResponse) GetVersion() string {
	return fmt.Sprintf("%d.%d", int(self.Version)&0x0F, int(self.Version)>>4)
}

// section 31.3
type GetSELAllocationInfoRequest struct {
}

type GetSELAllocationInfoResponse struct {
	// CompletionCode
	AllocationUnits            uint16 // LS Byte first
	AllocationSize             uint16 // LS Byte first
	FreeAllocationUnits        uint16 // LS Byte first
	LargestFreeAllocationUnits uint16 // LS Byte first
	MaxRecordSize              uint8
}

// section 31.4
type ReserveSELRequest struct {
}

type ReserveSELResponse struct {
	// CompletionCode
	Id uint16 // LS Byte first
}

// section 31.5
type GetSELRequest struct {
	ReservationId uint16 // LS Byte first

	RecordId      uint16 // LS Byte first
	Offset        uint8
	WillReadBytes uint8 // FFh means read entire record.
}

type GetSELResponse struct {
	// CompletionCode
	NextRecordId uint16 // LS Byte first
	Data         *RecordData
	DataLength   int
}

func (self *GetSELResponse) ReadBytes(r *protocol.Reader) {
	self.NextRecordId = r.ReadUint16()
	self.DataLength = r.Len()
	self.Data.Write(r.ReadBytes(r.Len()))
}

// section 31.10
type GetSELTimeRequest struct {
}

type GetSELTimeResponse struct {
	// CompletionCode
	Time uint32 // LS Byte first
}

// section 31.12
type GetAuxiliaryLogStatusRequest struct {
	Type uint8
}

type GetAuxiliaryLogStatusResponse struct {
	// CompletionCode
	Data []byte
}

func (self *GetAuxiliaryLogStatusResponse) ReadBytes(r *protocol.Reader) {
	self.Data = r.ReadCopy(r.Len())
}
