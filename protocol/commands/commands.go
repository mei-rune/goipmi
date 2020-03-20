package commands

import (
	"errors"
	"strings"
)

type PrivLevelType uint8

// PrivLevel
const (
	PrivLevelNone PrivLevelType = iota
	PrivLevelCallback
	PrivLevelUser
	PrivLevelOperator
	PrivLevelAdmin
	PrivLevelOEM
	PrivLevelUnprotected
)

func ParsePrivLevel(s string) (PrivLevelType, error) {
	switch strings.ToLower(s) {
	case "", "none":
		return PrivLevelNone, nil
	case "callback":
		return PrivLevelCallback, nil
	case "user":
		return PrivLevelUser, nil
	case "operator":
		return PrivLevelOperator, nil
	case "administrator":
		return PrivLevelAdmin, nil
	case "oem":
		return PrivLevelOEM, nil
	default:
		return PrivLevelNone, errors.New(s + " is unsupported priv level")
	}
}

// NetworkFunction identifies the functional class of an IPMI message
type NetworkFunction uint8

// Network Function Codes per section 5.1
const (
	NetworkFunctionChassis        = NetworkFunction(0x00)
	NetworkFunctionBridge         = NetworkFunction(0x02)
	NetworkFunctionSensor         = NetworkFunction(0x04)
	NetworkFunctionApp            = NetworkFunction(0x06)
	NetworkFunctionFirmware       = NetworkFunction(0x08)
	NetworkFunctionStorage        = NetworkFunction(0x0A)
	NetworkFunctionTransport      = NetworkFunction(0x0C)
	NetworkFunctionGroupExtension = NetworkFunction(0x2C)
)

// // Command fields on an IPMI message
// type Command uint8

// // Command Number Assignments (table G-1)
// const (
//   CommandGetDeviceID              = Command(0x01)
//   CommandGetAuthCapabilities      = Command(0x38)
//   CommandGetSessionChallenge      = Command(0x39)
//   CommandActivateSession          = Command(0x3a)
//   CommandSetSessionPrivilegeLevel = Command(0x3b)
//   CommandCloseSession             = Command(0x3c)
//   CommandChassisControl           = Command(0x02)
//   CommandChassisStatus            = Command(0x01)
//   CommandSetSystemBootOptions     = Command(0x08)
//   CommandGetSystemBootOptions     = Command(0x09)
// )

type CommandCode struct {
	Name            string
	NetworkFunction NetworkFunction
	Code            uint8
	PrivilegeLevel  PrivLevelType
}

// session commands
var (
	// Get Channel Authentication Capabilities  App  38h
	// Get Session Challenge                    App  39h
	// Activate Session                         App  3Ah
	// Set Session Privilege Level              App  3Bh
	// Close Session                            App  3Ch
	// reserved / unspecified  0
	// Close Channel 1
	// Close Channel 2
	// Close Channel 3
	// Close Channel 4
	// Close Channel 5
	// Close Channel 6
	// Close Channel 7
	// Close Channel 8
	// Close Channel 9
	// Close Channel Ah
	// Close Channel Bh
	GetChannelAuthenticationCapabilities = CommandCode{Name: "Get Channel Authentication Capabilities", NetworkFunction: NetworkFunctionApp, Code: 0x38, PrivilegeLevel: PrivLevelUnprotected}
	GetSessionChallenge                  = CommandCode{Name: "Get Session Challenge", NetworkFunction: NetworkFunctionApp, Code: 0x39, PrivilegeLevel: PrivLevelUnprotected}
	ActivateSession                      = CommandCode{Name: "Activate Session", NetworkFunction: NetworkFunctionApp, Code: 0x3A, PrivilegeLevel: PrivLevelUnprotected}
	SetSessionPrivilegeLevel             = CommandCode{Name: "Set Session Privilege Level", NetworkFunction: NetworkFunctionApp, Code: 0x3B, PrivilegeLevel: PrivLevelUser}
	CloseSession                         = CommandCode{Name: "Close Session", NetworkFunction: NetworkFunctionApp, Code: 0x3C, PrivilegeLevel: PrivLevelCallback}
	GetSessionInfo                       = CommandCode{Name: "Get Session Info", NetworkFunction: NetworkFunctionApp, Code: 0x3D, PrivilegeLevel: PrivLevelUser}
)
