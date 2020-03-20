package goipmi

import (
	"strconv"

	"github.com/runner-mei/goipmi/protocol"
)

// DeviceIDRequest per section 20.1
type DeviceIDRequest struct{}

// DeviceIDResponse per section 20.1
type DeviceIDResponse struct {
	// CompletionCode
	DeviceID                uint8
	DeviceRevision          uint8
	FirmwareRevision1       uint8
	FirmwareRevision2       uint8
	IPMIVersion             uint8
	AdditionalDeviceSupport uint8
	ManufacturerIDReserved  uint8
	ManufacturerID          protocol.OemID
	ProductID               uint16
}

// ColdReset per section 20.1
type ColdResetRequest struct{}

type ColdResetResponse struct {
	// CompletionCode
}

// WarmReset per section 20.1
type WarmResetRequest struct{}

type WarmResetResponse struct {
	// CompletionCode
}

// section 20.1
type SelfTestResultsRequest struct{}

type SelfTestResultsResponse struct {
	// CompletionCode
	Result uint8
	Detail uint8
}

func (self *SelfTestResultsResponse) IsOK() bool {
	return self.Result == 0x55
}

func (self *SelfTestResultsResponse) Message() string {
	switch self.Result {
	case 0x55:
		return "OK"
	case 0x56:
		return "NOTIMPLEMENTED"
	case 0x57:
		return "Corrupted or inaccessible data or device"
	case 0x58:
		return "Fatal hardware error"
	case 0xFF:
		return "reserved"
	default:
		return "device-specific failure"
	}
}

// section 20.1
type ManufacturingTestOnRequest struct {
	Text []byte
}

type ManufacturingTestOnResponse struct {
	// CompletionCode
}

// section 20.1
type SetACPIPowerStateRequest struct {
	SystemState uint8
	DeviceState uint8
}

type SetACPIPowerStateResponse struct {
	// CompletionCode
}

// section 20.1
type GetACPIPowerStateRequest struct {
}

type GetACPIPowerStateResponse struct {
	// CompletionCode
	SystemState uint8
	DeviceState uint8
}

func (self GetACPIPowerStateResponse) GetACPISystemPowerState() string {
	state := int(self.SystemState) & 0x8F
	switch state {
	case 0:
		return "S0/G0 working"
	case 1:
		return "S1"
	case 2:
		return "S2"
	case 3:
		return "S3"
	case 4:
		return "S4"
	case 5:
		return "S5/G2"
	case 6:
		return "S4/S5"
	case 7:
		return "G3"
	case 8:
		return "sleeping"
	case 9:
		return "G1 sleeping"
	case 0x0A:
		return "override"
	case 0x20:
		return "Legacy On"
	case 0x21:
		return "Legacy Off"
	case 0x2A:
		return "unknown"
	default:
		return "unknown(" + strconv.Itoa(state) + ")"
	}
}

func (self GetACPIPowerStateResponse) GetACPIDevicePowerState() string {
	state := int(self.DeviceState) & 0x8F
	switch state {
	case 0:
		return "D0"
	case 1:
		return "D1"
	case 2:
		return "D2"
	case 3:
		return "D3"
	case 0x2A:
		return "unknown"
	default:
		return "unknown(" + strconv.Itoa(state) + ")"
	}
}

// section 20.1
type GetDeviceGuidRequest struct {
}

type GetDeviceGuidResponse struct {
	// CompletionCode
	Guid [16]byte
}
