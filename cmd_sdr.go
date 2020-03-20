package goipmi

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/runner-mei/goipmi/protocol"
)

// section 33.9
type GetSDRInfoRequest struct {
}

type GetSDRInfoResponse struct {
	// CompletionCode
	Version     uint8  // minutes
	RecordCount uint16 // LS Byte first

	FreeSpace          uint16 // LS Byte first
	RecentAddTimestamp uint32 // LS Byte first
	RecentDelTimestamp uint32 // LS Byte first
	OperationSupport   uint8
}

func (self GetSDRInfoResponse) GetVersion() string {
	return fmt.Sprintf("%d.%d", int(self.Version)&0x0F, int(self.Version)>>4)
}

// section 33.10
type GetSDRAllocationInfoRequest struct {
}

type GetSDRAllocationInfoResponse struct {
	// CompletionCode
	AllocationUnitsLS          uint8 // LS Byte first
	AllocationUnitsMS          uint8
	AllocationSize             uint16 // LS Byte first
	FreeAllocationUnits        uint16 // LS Byte first
	LargestFreeAllocationUnits uint16 // LS Byte first
	MaxRecordSize              uint8
}

// section 33.11
type ReserveSDRRequest struct {
}

type ReserveSDRResponse struct {
	// CompletionCode
	Id uint16 // LS Byte first
}

// section 33.12
type GetSDRRequest struct {
	ReservationId uint16 // LS Byte first

	RecordId      uint16 // LS Byte first
	Offset        uint8
	WillReadBytes uint8 // FFh means read entire record.
}

type GetSDRResponse struct {
	// CompletionCode
	NextRecordId uint16 // LS Byte first
	Data         *RecordData
	DataLength   int
}

func (self *GetSDRResponse) ReadBytes(r *protocol.Reader) {
	self.NextRecordId = r.ReadUint16()
	self.DataLength = r.Len()
	if r.Len() > 0 {
		self.Data.Write(r.ReadBytes(r.Len()))
	}
}

// section 33.12
type GetSDRTimeRequest struct {
}

type GetSDRTimeResponse struct {
	// CompletionCode
	Time uint32 // LS Byte first
}

type RecordData struct {
	Data []byte
}

func (self *RecordData) SelIsOk() bool {
	return len(self.Data) == 16
}

func (self *RecordData) SdrIsOk() bool {
	if len(self.Data) < 5 {
		return false
	}
	return len(self.Data) >= int(uint8(self.Data[4]))
}

func (self *RecordData) SdrRecordLength() int {
	if len(self.Data) < 5 {
		return -1
	}
	return int(uint8(self.Data[4]))
}

func (self *RecordData) Write(bs []byte) {
	self.Data = append(self.Data, bs...)
}

func (self *RecordData) String() string {
	res, err := self.ToSdrRecord()
	if nil != err {
		return hex.EncodeToString(self.Data)
	}
	bs, _ := json.Marshal(res)
	return string(bs)
}

func (self *RecordData) ToSdrRecord() (Record, error) {
	var result Record

	switch self.Data[3] {
	case 0x01:
		result = &FullSensorRecord{}
	case 0x02:
		result = &CompactSensorRecord{}
	case 0x03:
		result = &EventOnlyRecord{}
	case 0x08:
		result = &EntityAssociationRecord{}
	case 0x09:
		result = &DeviceRelativeAssociationRecord{}
	case 0x10:
		result = &GenericDeviceLocatorRecord{}
	case 0x11:
		result = &FruDeviceLocatorRecord{}
	case 0x12:
		result = &McDeviceLocatorRecord{}
	case 0x13:
		result = &McDeviceConfirmationRecord{}
	case 0x14:
		result = &BMCMessageChannelInfoRecord{}
	case 0xC0:
		result = &OEMRecord{}
	default:
		return nil, errors.New("unknown record type - " + strconv.FormatUint(uint64(self.Data[3]), 10))
	}

	var reader = protocol.NewReader(self.Data)
	result.ReadBytes(reader)
	if reader.Err() != nil {
		return nil, reader.Err()
	}
	return result, nil
}
