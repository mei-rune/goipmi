package goipmi

import (
	"strconv"

	"github.com/runner-mei/goipmi/protocol"
)

// section 28
type GetChassisCapabilitiesRequest struct {
}

type GetChassisCapabilitiesResponse struct {
	// CompletionCode
	CapabilitiesFlags                    uint8
	ChassisFruInfoDeviceAddress          uint8
	ChassisSDRDeviceAddress              uint8
	ChassisSELDeviceAddress              uint8
	ChassisSystemManagementDeviceAddress uint8
	ChassisBridgeDeviceAddress           uint8 // option
}

func (self *GetChassisCapabilitiesResponse) PowerInterlock() bool {
	return self.CapabilitiesFlags&(uint8(1)<<3) != 0
}

func (self *GetChassisCapabilitiesResponse) DiagnosticInterrupt() bool {
	return self.CapabilitiesFlags&(uint8(1)<<2) != 0
}

func (self *GetChassisCapabilitiesResponse) FrontPanelLockout() bool {
	return self.CapabilitiesFlags&(uint8(1)<<1) != 0
}

func (self *GetChassisCapabilitiesResponse) HasInstrusionSensor() bool {
	return self.CapabilitiesFlags&(uint8(1)<<0) != 0
}

func (self *GetChassisCapabilitiesResponse) ReadBytes(r *protocol.Reader) {
	if r.Len() < 5 {
		r.SetError(ErrInsufficientBytes)
		return
	}

	bs := r.ReadBytes(5)
	self.CapabilitiesFlags = uint8(bs[0])
	self.ChassisFruInfoDeviceAddress = uint8(bs[1])
	self.ChassisSDRDeviceAddress = uint8(bs[2])
	self.ChassisSELDeviceAddress = uint8(bs[3])
	self.ChassisSystemManagementDeviceAddress = uint8(bs[4])

	if r.Len() >= 6 {
		self.ChassisBridgeDeviceAddress = r.ReadUint8() // option
	}
}

// section 28
type SetChassisCapabilitiesRequest struct {
	CapabilitiesFlags                    uint8
	ChassisFruInfoDeviceAddress          uint8
	ChassisSDRDeviceAddress              uint8
	ChassisSELDeviceAddress              uint8
	ChassisSystemManagementDeviceAddress uint8
	HasChassisBridgeDeviceAddress        bool
	ChassisBridgeDeviceAddress           uint8 // option
}

func (self *SetChassisCapabilitiesRequest) WriteBytes(w *protocol.Writer) {
	w.WriteUint8(self.CapabilitiesFlags)
	w.WriteUint8(self.ChassisFruInfoDeviceAddress)
	w.WriteUint8(self.ChassisSDRDeviceAddress)
	w.WriteUint8(self.ChassisSELDeviceAddress)
	w.WriteUint8(self.ChassisSystemManagementDeviceAddress)

	if self.HasChassisBridgeDeviceAddress {
		w.WriteUint8(self.ChassisBridgeDeviceAddress)
	}
}

type SetChassisCapabilitiesResponse struct {
	// CompletionCode
}

// section 28
type GetChassisStatusRequest struct {
}

type GetChassisStatusResponse struct {
	// CompletionCode
	CurrentPowerState            uint8
	LastPowerEvent               uint8
	ChassisState                 uint8
	FrontPanelCapAndEnableStatus uint8 // option
}

func (self *GetChassisStatusResponse) PowerRestorePolicy() int {
	return int(self.CurrentPowerState) >> 5
}

func (self *GetChassisStatusResponse) PowerRestorePolicyString() string {
	state := self.PowerRestorePolicy()
	switch state {
	case 0:
		return "power off"
	case 1:
		return "power restore"
	case 2:
		return "power up"
	case 3:
		return "unknown"
	default:
		return "unknown(" + strconv.Itoa(state) + ")"
	}
}

func (self *GetChassisStatusResponse) PowerControlFault() bool {
	return self.CurrentPowerState&(uint8(1)<<4) != 0
}

func (self *GetChassisStatusResponse) PowerFault() bool {
	return self.CurrentPowerState&(uint8(1)<<3) != 0
}

func (self *GetChassisStatusResponse) Interlock() bool {
	return self.CurrentPowerState&(uint8(1)<<2) != 0
}

func (self *GetChassisStatusResponse) PowerOverload() bool {
	return self.CurrentPowerState&(uint8(1)<<1) != 0
}

func (self *GetChassisStatusResponse) PowerOn() bool {
	return self.CurrentPowerState&(uint8(1)<<1) != 0
}

func (self *GetChassisStatusResponse) LastPowerEventPowerOn() bool {
	return self.LastPowerEvent&(uint8(1)<<4) != 0
}

func (self *GetChassisStatusResponse) LastPowerEventPowerDownByFault() bool {
	return self.LastPowerEvent&(uint8(1)<<3) != 0
}

func (self *GetChassisStatusResponse) LastPowerEventPowerDownByInterlock() bool {
	return self.LastPowerEvent&(uint8(1)<<2) != 0
}

func (self *GetChassisStatusResponse) LastPowerEventPowerDownByOverload() bool {
	return self.LastPowerEvent&(uint8(1)<<1) != 0
}

func (self *GetChassisStatusResponse) LastPowerEventACFailed() bool {
	return self.LastPowerEvent&(uint8(1)<<0) != 0
}

func (self *GetChassisStatusResponse) IdentityCommandSupported() bool {
	return self.ChassisState&(uint8(1)<<6) != 0
}

func (self *GetChassisStatusResponse) FanFault() bool {
	return self.ChassisState&(uint8(1)<<3) != 0
}

func (self *GetChassisStatusResponse) DriverFault() bool {
	return self.ChassisState&(uint8(1)<<2) != 0
}

func (self *GetChassisStatusResponse) FrontPanelLockoutActived() bool {
	return self.ChassisState&(uint8(1)<<1) != 0
}

func (self *GetChassisStatusResponse) ChassisInstrusionActived() bool {
	return self.ChassisState&(uint8(1)<<0) != 0
}

func (self *GetChassisStatusResponse) FrontPanelStandbyButtonDisableAllowed() bool {
	return self.FrontPanelCapAndEnableStatus&(uint8(1)<<7) != 0
}

func (self *GetChassisStatusResponse) FrontPanelDiagnosticInterruptButtonDisableAllowed() bool {
	return self.FrontPanelCapAndEnableStatus&(uint8(1)<<6) != 0
}

func (self *GetChassisStatusResponse) FrontPanelResetButtonDisableAllowed() bool {
	return self.FrontPanelCapAndEnableStatus&(uint8(1)<<5) != 0
}

func (self *GetChassisStatusResponse) FrontPanelPowerOffButtonDisableAllowed() bool {
	return self.FrontPanelCapAndEnableStatus&(uint8(1)<<4) != 0
}

func (self *GetChassisStatusResponse) FrontPanelStandbyButtonDisabled() bool {
	return self.FrontPanelCapAndEnableStatus&(uint8(1)<<3) != 0
}

func (self *GetChassisStatusResponse) FrontPanelDiagnosticInterruptButtonDisabled() bool {
	return self.FrontPanelCapAndEnableStatus&(uint8(1)<<2) != 0
}

func (self *GetChassisStatusResponse) FrontPanelResetButtonDisabled() bool {
	return self.FrontPanelCapAndEnableStatus&(uint8(1)<<1) != 0
}

func (self *GetChassisStatusResponse) FrontPanelPowerOffButtonDisabled() bool {
	return self.FrontPanelCapAndEnableStatus&(uint8(1)<<0) != 0
}

func (self *GetChassisStatusResponse) ReadBytes(r *protocol.Reader) {
	if r.Len() < 3 {
		r.SetError(ErrInsufficientBytes)
		return
	}

	bs := r.ReadBytes(3)
	self.CurrentPowerState = uint8(bs[0])
	self.LastPowerEvent = uint8(bs[1])
	self.ChassisState = uint8(bs[2])

	if r.Len() >= 4 {
		self.FrontPanelCapAndEnableStatus = r.ReadUint8() // option
	}
}

// section 28
type ChassisControlRequest struct {
	Code uint8
}

type ChassisControlCode uint8

const (
	ChassisControl_PowerDown               ChassisControlCode = 0
	ChassisControl_PowerUp                 ChassisControlCode = 1
	ChassisControl_PowerCycle              ChassisControlCode = 2
	ChassisControl_PowerHardReset          ChassisControlCode = 3
	ChassisControl_PlusDiagnosticInterrupt ChassisControlCode = 4
	ChassisControl_SoftShutdown            ChassisControlCode = 4
)

type ChassisControlResponse struct {
	// CompletionCode
}

// section 28
type ChassisResetRequest struct {
}

type ChassisResetResponse struct {
	// CompletionCode
}

// section 28
type ChassisIdentifyRequest struct {
	IdentifyInterval uint8 // option
	ForceIdentifyOn  uint8 // option
}

type ChassisIdentifyResponse struct {
	// CompletionCode
}

// section 28
type SetPanelEnableRequest struct {
	Enable uint8
}

type SetPanelEnableResponse struct {
	// CompletionCode
}

// section 28
type SetPowerRestorePolicyRequest struct {
	Policy uint8
}

type PowerRestorePolicyType uint8

const (
	PowerRestorePolicy_NoChange   PowerRestorePolicyType = 3
	PowerRestorePolicy_AlwaysUp   PowerRestorePolicyType = 2
	PowerRestorePolicy_Restore    PowerRestorePolicyType = 1
	PowerRestorePolicy_AlwaysDown PowerRestorePolicyType = 0
)

type SetPowerRestorePolicyResponse struct {
	// CompletionCode
	Policy uint8
}

func (self *SetPowerRestorePolicyResponse) AlwaysUp() bool {
	return self.Policy&(1<<2) != 0
}

func (self *SetPowerRestorePolicyResponse) AlwaysRestore() bool {
	return self.Policy&(1<<1) != 0
}

func (self *SetPowerRestorePolicyResponse) AlwaysDown() bool {
	return self.Policy&1 != 0
}

// section 28
type SetPowerCycleIntervalRequest struct {
	Interval uint8 // seconds
}

type SetPowerCycleIntervalResponse struct {
	// CompletionCode
}

// section 28
type GetSystemRestartCauseRequest struct {
}

type GetSystemRestartCauseResponse struct {
	// CompletionCode
	ResetCause    uint8
	ChannelNumber uint8
}

func (self *GetSystemRestartCauseResponse) GetResetCause() uint8 {
	return self.ResetCause & uint8(0x0F)
}

func (self *GetSystemRestartCauseResponse) GetResetCauseString() string {
	state := self.GetResetCause()
	switch state {
	case 0:
		return "unknown"
	case 1:
		return "chassis control command"
	case 2:
		return "reset by pushbutton"
	case 3:
		return "power up via power pushbutton"
	case 4:
		return "watchdog expiration"
	case 5:
		return "OEM"
	case 6:
		return "auto power up on AC being applied 'always restore'"
	case 7:
		return "auto power up on AC being applied 'restore previous power state'"
	case 8:
		return "reset by PEF"
	case 9:
		return "power cycle via PEF"
	case 10:
		return "soft reset"
	case 11:
		return "power up via RTC"
	default:
		return "unknown(" + strconv.Itoa(int(state)) + ")"
	}
}

// section 28
type SetSystemBootOptionsRequest struct {
	ParameterValid uint8
	ParameterData  []byte
}

func (self *SetSystemBootOptionsRequest) WriteBytes(w *protocol.Writer) {
	w.WriteUint8(self.ParameterValid)
	w.WriteBytes(self.ParameterData)
}

type SetSystemBootOptionsResponse struct {
	// CompletionCode
}

// section 28
type GetSystemBootOptionsRequest struct {
	ParameterSelector uint8
	SetSelector       uint8
	BlockSelector     uint8
}

type GetSystemBootOptionsResponse struct {
	// CompletionCode
	ParameterVersion uint8
	ParameterValid   uint8
	ParameterData    SystemBootParameterData
}

func (self *GetSystemBootOptionsResponse) ParameterIsValid() bool {
	return self.ParameterValid&(1<<7) != 0
}

// func (self *GetSystemBootOptionsResponse) ReadBytes(r *protocol.Reader) {
// 	if r.Len() < 2 {
// 		r.SetError(ipmi.ErrInsufficientBytes)
// 		return
// 	}

// 	bs := r.ReadBytes(2)
// 	self.ParameterVersion = uint8(bs[0])
// 	self.ParameterValid = uint8(bs[1])
// 	self.ParameterData.ReadBytes(r)
// }

type SystemBootParameterData struct {
	SetInProgress               uint8
	ServicePartitionSelector    uint8
	ServicePartitionScan        uint8
	BMCBootFlagValidBitClearing uint8
	BootInfoAcknowledge         [2]uint8
	BootFlags                   [5]uint8
	BootInitialorInfo           [9]uint8
	BootInitialorMailbox        [17]uint8
	//ssfsfd
}

// func (self *SystemBootParameterData) ReadBytes(r *protocol.Reader) {
// }

// section 28
type GetPOHCounterRequest struct {
}

type GetPOHCounterResponse struct {
	// CompletionCode
	Count          uint8  // n minutes
	CounterReading uint32 // LS Byte first
}
