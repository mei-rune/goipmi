package goipmi

// section 22.1
type SetBMCGlobalEnablesRequest struct {
	Flages uint8
}

func (self *SetBMCGlobalEnablesRequest) SetOEM2Enable(b bool) {
	if b {
		self.Flages = self.Flages | (uint8(1) << 7)
	} else {
		self.Flages = self.Flages & ^(uint8(1) << 7)
	}
}

func (self *SetBMCGlobalEnablesRequest) SetOEM1Enable(b bool) {
	if b {
		self.Flages = self.Flages | (uint8(1) << 6)
	} else {
		self.Flages = self.Flages & ^(uint8(1) << 6)
	}
}

func (self *SetBMCGlobalEnablesRequest) SetOEM0Enable(b bool) {
	if b {
		self.Flages = self.Flages | (uint8(1) << 5)
	} else {
		self.Flages = self.Flages & ^(uint8(1) << 5)
	}
}

func (self *SetBMCGlobalEnablesRequest) EnableSystemEventLogging() {
	self.Flages = self.Flages | (uint8(1) << 3)
}

func (self *SetBMCGlobalEnablesRequest) DisableSystemEventLogging() {
	self.Flages = self.Flages & ^(uint8(1) << 3)
}

func (self *SetBMCGlobalEnablesRequest) EnableEventMessageBuffer() {
	self.Flages = self.Flages | (uint8(1) << 2)
}

func (self *SetBMCGlobalEnablesRequest) DisableEventMessageBuffer() {
	self.Flages = self.Flages & ^(uint8(1) << 2)
}

func (self *SetBMCGlobalEnablesRequest) EnableEventMessageBufferFullInterrupt() {
	self.Flages = self.Flages | (uint8(1) << 1)
}

func (self *SetBMCGlobalEnablesRequest) DisableEventMessageBufferFullInterrupt() {
	self.Flages = self.Flages & ^(uint8(1) << 1)
}

func (self *SetBMCGlobalEnablesRequest) EnableEventMessageBufferQueueInterrupt() {
	self.Flages = self.Flages | (uint8(1) << 0)
}

func (self *SetBMCGlobalEnablesRequest) DisableEventMessageBufferQueueInterrupt() {
	self.Flages = self.Flages & ^(uint8(1) << 0)
}

// section 22.1
type SetBMCGlobalEnablesResponse struct {
	// CompletionCode
}

// section 22.2
type GetBMCGlobalEnablesRequest struct {
}

type GetBMCGlobalEnablesResponse struct {
	// CompletionCode
	Flages uint8
}

func (self *SetBMCGlobalEnablesRequest) OEM2Enabled() bool {
	return self.Flages&(uint8(1)<<7) == 0
}

func (self *SetBMCGlobalEnablesRequest) OEM1Enabled() bool {
	return self.Flages&(uint8(1)<<6) == 0
}

func (self *SetBMCGlobalEnablesRequest) OEM0Enabled() bool {
	return self.Flages&(uint8(1)<<5) == 0
}

func (self *SetBMCGlobalEnablesRequest) SystemEventLoggingEnabled() bool {
	return self.Flages&(uint8(1)<<3) == 0
}

func (self *SetBMCGlobalEnablesRequest) EventMessageBufferEnabled() bool {
	return self.Flages&(uint8(1)<<2) == 0
}

func (self *SetBMCGlobalEnablesRequest) EventMessageBufferFullInterruptEnabled() bool {
	return self.Flages&(uint8(1)<<1) == 0
}

func (self *SetBMCGlobalEnablesRequest) EventMessageBufferQueueInterruptEnabled() bool {
	return self.Flages&(uint8(1)<<0) == 0
}

// section 22.14b
type GetSystemInfoParametersRequest struct {
	Parameters        uint8
	ParameterSelector uint8
	SetSelector       uint8
	BlockSelector     uint8
}

func (self *GetSystemInfoParametersRequest) SetParameters(b bool) {
	if b {
		self.Parameters = self.Parameters | (uint8(1) << 7)
	} else {
		self.Parameters = self.Parameters & ^(uint8(1) << 7)
	}
}

// type GetSystemInfoParametersResponse struct {
// 	// CompletionCode
// 	ParameterRevision uint16
// 	ParameterData     ParameterData
// }

// func (self *GetSystemInfoParametersResponse) ReadBytes(r *ipmi.Reader) {
// 	self.ParameterRevision = r.ReadUint16()
// 	self.ParameterData.ReadBytes(r)
// }

// type ParameterData struct {
//   SetInProgress uint8
//   SystemFirmware uint8
// }

// func (self *GetSystemInfoParametersResponse) ReadBytes(r *ipmi.Reader) {
//   self.ParameterRevision = r.ReadUint16()
//   self.ParameterData.ReadBytes(r)
// }
