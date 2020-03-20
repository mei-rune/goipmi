package goipmi

import "github.com/runner-mei/goipmi/protocol/commands"

var (
	// IPM Device “Global” Commands
	// reserved= commands.CommandCode{Name: "reserved", NetworkFunction: commands.NetworkFunctionApp, Code: 0x00}
	GetDeviceID          = commands.CommandCode{Name: "Get Device ID", NetworkFunction: commands.NetworkFunctionApp, Code: 0x01, PrivilegeLevel: commands.PrivLevelUser}
	BroadcastGetDeviceID = commands.CommandCode{Name: "Broadcast 'Get Device ID'", NetworkFunction: commands.NetworkFunctionApp, Code: 0x01, PrivilegeLevel: commands.PrivLevelNone} // Local only
	ColdReset            = commands.CommandCode{Name: "Cold Reset", NetworkFunction: commands.NetworkFunctionApp, Code: 0x02, PrivilegeLevel: commands.PrivLevelAdmin}
	WarmReset            = commands.CommandCode{Name: "Warm Reset", NetworkFunction: commands.NetworkFunctionApp, Code: 0x03, PrivilegeLevel: commands.PrivLevelAdmin}
	GetSelfTestResults   = commands.CommandCode{Name: "Get Self Test Results", NetworkFunction: commands.NetworkFunctionApp, Code: 0x04, PrivilegeLevel: commands.PrivLevelUser}
	ManufacturingTestOn  = commands.CommandCode{Name: "Manufacturing Test On", NetworkFunction: commands.NetworkFunctionApp, Code: 0x05, PrivilegeLevel: commands.PrivLevelAdmin}
	SetACPIPowerState    = commands.CommandCode{Name: "Set ACPI Power State", NetworkFunction: commands.NetworkFunctionApp, Code: 0x06, PrivilegeLevel: commands.PrivLevelAdmin}
	GetACPIPowerState    = commands.CommandCode{Name: "Get ACPI Power State", NetworkFunction: commands.NetworkFunctionApp, Code: 0x07, PrivilegeLevel: commands.PrivLevelUser}
	GetDeviceGUID        = commands.CommandCode{Name: "Get Device GUID", NetworkFunction: commands.NetworkFunctionApp, Code: 0x08, PrivilegeLevel: commands.PrivLevelUser}

	// GetNetFnSupport= commands.CommandCode{Name: "Get NetFn Support", NetworkFunction: commands.NetworkFunctionApp, Code: 0x09, PrivilegeLevel: commands.PrivLevelUser}
	// GetCommandSupport= commands.CommandCode{Name: "Get Command Support", NetworkFunction: commands.NetworkFunctionApp, Code: 0x0A, PrivilegeLevel: commands.PrivLevelUser}
	// GetCommandSubFunctionSupport= commands.CommandCode{Name: "Get Command Sub-function Support", NetworkFunction: commands.NetworkFunctionApp, Code: 0x0B, PrivilegeLevel: commands.PrivLevelUser}

	//reserved                                           App  09h-0Fh
	SetCommandEnables            = commands.CommandCode{Name: "Set Command Enables", NetworkFunction: commands.NetworkFunctionApp, Code: 0x60, PrivilegeLevel: commands.PrivLevelAdmin}
	GetCommandEnables            = commands.CommandCode{Name: "Get Command Enables", NetworkFunction: commands.NetworkFunctionApp, Code: 0x61, PrivilegeLevel: commands.PrivLevelUser}
	SetCommandSubFunctionEnables = commands.CommandCode{Name: "Set Command Sub-function Enables", NetworkFunction: commands.NetworkFunctionApp, Code: 0x62, PrivilegeLevel: commands.PrivLevelAdmin}
	GetCommandSubFunctionEnables = commands.CommandCode{Name: "Get Command Sub-function Enables", NetworkFunction: commands.NetworkFunctionApp, Code: 0x63, PrivilegeLevel: commands.PrivLevelUser}
	GetOEMNetFnIANASupport       = commands.CommandCode{Name: "Get OEM NetFn IANA Support", NetworkFunction: commands.NetworkFunctionApp, Code: 0x64, PrivilegeLevel: commands.PrivLevelUser}

	// // MC Watchdog Timer  Commands
	// Reset Watchdog Timer                               App  22h
	// Set Watchdog Timer                                 App  24h
	// Set Timer Use Field                                          0
	// Set Timer Actions                                            1
	// Clear Timer Use Expiration Flags                             2
	// Set Countdown value                                          3
	// Get Watchdog Timer                                 App  25h
	ResetWatchdogTimer = commands.CommandCode{Name: "Reset Watchdog Timer", NetworkFunction: commands.NetworkFunctionApp, Code: 0x22, PrivilegeLevel: commands.PrivLevelOperator}
	SetWatchdogTimer   = commands.CommandCode{Name: "Set Watchdog Timer", NetworkFunction: commands.NetworkFunctionApp, Code: 0x24, PrivilegeLevel: commands.PrivLevelOperator}
	GetWatchdogTimer   = commands.CommandCode{Name: "Get Watchdog Timer", NetworkFunction: commands.NetworkFunctionApp, Code: 0x25, PrivilegeLevel: commands.PrivLevelUser}

	// //BMC Device and Messaging Commands
	// Set BMC Global Enables                             App  2Eh
	// Change message queue interrupt enable                   0
	// Change event message buffer full
	// interrupt enable                                         1
	// Change event message buffer enable                   2
	// Change System Event Logging enable                    3
	// reserved / unspecified                                -
	// Change OEM 0 enable                                   5
	// Change OEM 1 enable                                      6
	// Change OEM 2 enable                                      7
	SetBMCGlobalEnables = commands.CommandCode{Name: "Set BMC Global Enables", NetworkFunction: commands.NetworkFunctionApp, Code: 0x2E, PrivilegeLevel: commands.PrivLevelNone}
	GetBMCGlobalEnables = commands.CommandCode{Name: "Get BMC Global Enables", NetworkFunction: commands.NetworkFunctionApp, Code: 0x2F, PrivilegeLevel: commands.PrivLevelUser}
	ClearMessageFlags   = commands.CommandCode{Name: "Clear Message Flags", NetworkFunction: commands.NetworkFunctionApp, Code: 0x30, PrivilegeLevel: commands.PrivLevelNone}
	// Receive Message Queue clear                   0
	// Event Message Buffer clear                   1
	// reserved / unspecified                   2
	// Watchdog pre                              -
	// timeout interrupt clear                   3
	// reserved / unspecified                   4
	// OEM 0 clear                              5
	// OEM 1 clear                              6
	// OEM 2 clear                              7
	GetMessageFlags             = commands.CommandCode{Name: "Get Message Flags", NetworkFunction: commands.NetworkFunctionApp, Code: 0x31, PrivilegeLevel: commands.PrivLevelNone}
	EnableMessageChannelReceive = commands.CommandCode{Name: "Enable Message Channel Receive", NetworkFunction: commands.NetworkFunctionApp, Code: 0x32, PrivilegeLevel: commands.PrivLevelNone}

	// reserved / unspecified                      0
	// Channel 1 enable/disable                    1
	// Channel 2 enable/disable                    2
	// Channel 3 enable/disable                    3
	// Channel 4 enable/disable                    4
	// Channel 5 enable/disable                    5
	// Channel 6 enable/disable                    6
	// Channel 7 enable/disable                    7
	// Channel 8 enable/disable                    8
	// Channel 9 enable/disable                    9
	// Channel Ah enable/disable                   10
	// Channel Bh enable/disable                   11
	GetMessage  = commands.CommandCode{Name: "Get Message", NetworkFunction: commands.NetworkFunctionApp, Code: 0x33, PrivilegeLevel: commands.PrivLevelNone}
	SendMessage = commands.CommandCode{Name: "Send Message", NetworkFunction: commands.NetworkFunctionApp, Code: 0x34, PrivilegeLevel: commands.PrivLevelUser} // Administator for some channels
	// Send to channel 0
	// Send to channel 1
	// Send to channel 2
	// Send to channel 3
	// Send to channel 4
	// Send to channel 5
	// Send to channel 6
	// Send to channel 7
	// Send to channel 8
	// Send to channel 9
	// Send to channel Ah
	// Send to channel Bh
	ReadEventMessageBuffer     = commands.CommandCode{Name: "Read Event Message Buffer", NetworkFunction: commands.NetworkFunctionApp, Code: 0x35, PrivilegeLevel: commands.PrivLevelNone} // System interface only
	GetBTInterfaceCapabilities = commands.CommandCode{Name: "Get BT Interface Capabilities", NetworkFunction: commands.NetworkFunctionApp, Code: 0x36, PrivilegeLevel: commands.PrivLevelUser}
	GetSystemGUID              = commands.CommandCode{Name: "Get System GUID", NetworkFunction: commands.NetworkFunctionApp, Code: 0x37, PrivilegeLevel: commands.PrivLevelUnprotected}
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
	GetChannelAuthenticationCapabilities = commands.GetChannelAuthenticationCapabilities
	GetSessionChallenge                  = commands.GetSessionChallenge
	ActivateSession                      = commands.ActivateSession
	SetSessionPrivilegeLevel             = commands.SetSessionPrivilegeLevel
	CloseSession                         = commands.CloseSession
	GetSessionInfo                       = commands.GetSessionInfo

	// unassigned= commands.CommandCode{Name: "unassigned", NetworkFunction: commands.NetworkFunctionApp, Code: 0x3E}
	GetAuthCode = commands.CommandCode{Name: "Get AuthCode", NetworkFunction: commands.NetworkFunctionApp, Code: 0x3F, PrivilegeLevel: commands.PrivLevelOperator}
	// Change configuration for channel 0
	// Change configuration for channel 1
	// Change configuration for channel 2
	// Change configuration for channel 3
	// Change configuration for channel 4
	// Change configuration for channel 5
	// Change configuration for channel 6
	// Change configuration for channel 7
	// Change configuration for channel 8
	// Change configuration for channel 9
	// Change configuration for channel Ah
	// Change configuration for channel Bh
	// Get Channel Access                      App  41h
	SetChannelAccess      = commands.CommandCode{Name: "Set Channel Access", NetworkFunction: commands.NetworkFunctionApp, Code: 0x40, PrivilegeLevel: commands.PrivLevelAdmin}
	GetChannelAccess      = commands.CommandCode{Name: "Get Channel Access", NetworkFunction: commands.NetworkFunctionApp, Code: 0x41, PrivilegeLevel: commands.PrivLevelUser}
	GetChannelInfoCommand = commands.CommandCode{Name: "Get Channel Info Command", NetworkFunction: commands.NetworkFunctionApp, Code: 0x42, PrivilegeLevel: commands.PrivLevelUser}

	SetUserAccessCommand   = commands.CommandCode{Name: "Set User Access Command", NetworkFunction: commands.NetworkFunctionApp, Code: 0x43, PrivilegeLevel: commands.PrivLevelAdmin}
	GetUserAccessCommand   = commands.CommandCode{Name: "Get User Access Command", NetworkFunction: commands.NetworkFunctionApp, Code: 0x44, PrivilegeLevel: commands.PrivLevelOperator}
	SetUserName            = commands.CommandCode{Name: "Set User Name", NetworkFunction: commands.NetworkFunctionApp, Code: 0x45, PrivilegeLevel: commands.PrivLevelAdmin}
	GetUserNameCommand     = commands.CommandCode{Name: "Get User Name Command", NetworkFunction: commands.NetworkFunctionApp, Code: 0x46, PrivilegeLevel: commands.PrivLevelOperator}
	SetUserPasswordCommand = commands.CommandCode{Name: "Set User Password Command", NetworkFunction: commands.NetworkFunctionApp, Code: 0x47, PrivilegeLevel: commands.PrivLevelAdmin}
	ActivatePayload        = commands.CommandCode{Name: "Activate Payload", NetworkFunction: commands.NetworkFunctionApp, Code: 0x48}   // Depends on payload type.
	DeactivatePayload      = commands.CommandCode{Name: "Deactivate Payload", NetworkFunction: commands.NetworkFunctionApp, Code: 0x49} // Depends on payload type.

	SetUserPayloadAccess     = commands.CommandCode{Name: "Set User Payload Access", NetworkFunction: commands.NetworkFunctionApp, Code: 0x4C, PrivilegeLevel: commands.PrivLevelAdmin}
	GetUserPayloadAccess     = commands.CommandCode{Name: "Get User Payload Access", NetworkFunction: commands.NetworkFunctionApp, Code: 0x4D, PrivilegeLevel: commands.PrivLevelOperator}
	GetChannelPayloadSupport = commands.CommandCode{Name: "Get Channel Payload Support", NetworkFunction: commands.NetworkFunctionApp, Code: 0x4E, PrivilegeLevel: commands.PrivLevelUser}
	GetChannelPayloadVersion = commands.CommandCode{Name: "Get Channel Payload Version", NetworkFunction: commands.NetworkFunctionApp, Code: 0x4F, PrivilegeLevel: commands.PrivLevelUser}
	GetChannelOEMPayloadInfo = commands.CommandCode{Name: "Get Channel OEM Payload Info", NetworkFunction: commands.NetworkFunctionApp, Code: 0x50, PrivilegeLevel: commands.PrivLevelUser}

	// unassigned= commands.CommandCode{Name: "unassigned", NetworkFunction: commands.NetworkFunctionApp, Code: 0x51}
	MasterWriteRead = commands.CommandCode{Name: "Master Write-Read", NetworkFunction: commands.NetworkFunctionApp, Code: 0x52, PrivilegeLevel: commands.PrivLevelOperator}
	// reserved / unspecified
	// 0
	// Access to public bus, channel 1
	// 1
	// Access to public bus, channel 2
	// 2
	// Access to public bus, channel 3
	// 3
	// Access to public bus, channel 4
	// 4
	// Access to public bus, cha
	// nnel 5
	// 5
	// Access to public bus, channel 6
	// 6
	// Access to public bus, channel 7
	// 7
	// Access to private bus 0
	// 8
	// Access to private bus 1
	// 9
	// Access to private bus 2
	// 10
	// Access to private bus 3
	// 11
	// Access to private bus 4
	// 12
	// Access to private bus
	// 5
	// 13
	// Access to private bus 6
	// 14
	// Access to private bus 7
	// 15
	// Access to public bus, channel 8
	// 16
	// Access to public bus, channel 9
	// 17
	// Access to public bus, channel Ah
	// 18
	// Access to public bus, channel Bh
	// 19

	// unassigned= commands.CommandCode{Name: "unassigned", NetworkFunction: commands.NetworkFunctionApp, Code: 0x53}

	GetChannelCipherSuites         = commands.CommandCode{Name: "Get Channel Cipher Suites", NetworkFunction: commands.NetworkFunctionApp, Code: 0x54, PrivilegeLevel: commands.PrivLevelUnprotected}
	SuspendResumePayloadEncryption = commands.CommandCode{Name: "Suspend/Resume Payload Encryption", NetworkFunction: commands.NetworkFunctionApp, Code: 0x55, PrivilegeLevel: commands.PrivLevelUser}
	SetChannelSecurityKeys         = commands.CommandCode{Name: "Set Channel Security Keys", NetworkFunction: commands.NetworkFunctionApp, Code: 0x56, PrivilegeLevel: commands.PrivLevelAdmin}
	GetSystemInterfaceCapabilities = commands.CommandCode{Name: "Get System Interface Capabilities", NetworkFunction: commands.NetworkFunctionApp, Code: 0x57, PrivilegeLevel: commands.PrivLevelUser}
	SetSystemInfoParameters        = commands.CommandCode{Name: "Set System Info Parameters", NetworkFunction: commands.NetworkFunctionApp, Code: 0x58, PrivilegeLevel: commands.PrivLevelAdmin}
	GetSystemInfoParameters        = commands.CommandCode{Name: "Get System Info Parameters", NetworkFunction: commands.NetworkFunctionApp, Code: 0x59, PrivilegeLevel: commands.PrivLevelUser}

	// Chassis Device Commands
	GetChassisCapabilities = commands.CommandCode{Name: "Get Chassis Capabilities", NetworkFunction: commands.NetworkFunctionChassis, Code: 0x00, PrivilegeLevel: commands.PrivLevelUser}
	GetChassisStatus       = commands.CommandCode{Name: "Get Chassis Status", NetworkFunction: commands.NetworkFunctionChassis, Code: 0x01, PrivilegeLevel: commands.PrivLevelUser}

	// reserved unspecified  0
	// power up  1
	// power cycle 2
	// hard reset 3
	// pulse diagnostic interrupt  4
	// initiate soft shutdown via overtemp  5
	// Chassis Reset
	ChassisControl  = commands.CommandCode{Name: "Chassis Control", NetworkFunction: commands.NetworkFunctionChassis, Code: 0x02, PrivilegeLevel: commands.PrivLevelOperator}
	ChassisReset    = commands.CommandCode{Name: "Chassis Reset", NetworkFunction: commands.NetworkFunctionChassis, Code: 0x03, PrivilegeLevel: commands.PrivLevelOperator}
	ChassisIdentify = commands.CommandCode{Name: "Chassis Identify", NetworkFunction: commands.NetworkFunctionChassis, Code: 0x04, PrivilegeLevel: commands.PrivLevelOperator}
	// Force On Indefinitely  0

	SetFrontPanelButtonEnables = commands.CommandCode{Name: "Set Front Panel Button Enables", NetworkFunction: commands.NetworkFunctionChassis, Code: 0x0A, PrivilegeLevel: commands.PrivLevelAdmin}

	// Power off via front panel                      0
	// Reset via front panel                    1
	// Diagnostic Interrupt via front panel     2
	// Standby (sleep) via front panel          3
	SetChassisCapabilities = commands.CommandCode{Name: "Set Chassis Capabilities", NetworkFunction: commands.NetworkFunctionChassis, Code: 0x05, PrivilegeLevel: commands.PrivLevelAdmin}
	SetPowerRestorePolicy  = commands.CommandCode{Name: "Set Power Restore Policy", NetworkFunction: commands.NetworkFunctionChassis, Code: 0x06, PrivilegeLevel: commands.PrivLevelOperator}
	SetPowerCycleInterval  = commands.CommandCode{Name: "Set Power Cycle Interval", NetworkFunction: commands.NetworkFunctionChassis, Code: 0x0B, PrivilegeLevel: commands.PrivLevelAdmin}
	GetSystemRestartCause  = commands.CommandCode{Name: "Get System Restart Cause", NetworkFunction: commands.NetworkFunctionChassis, Code: 0x07, PrivilegeLevel: commands.PrivLevelUser}
	SetSystemBootOptions   = commands.CommandCode{Name: "Set System Boot Options", NetworkFunction: commands.NetworkFunctionChassis, Code: 0x08, PrivilegeLevel: commands.PrivLevelOperator}
	// reserved / unspecified                            0
	// Write parameter 1 (service partition selector)    1
	// Write parameter 2 (service partition scan)        2
	// Write parameter 3 (‘valid bit’ clearing)          3
	// Write parameter 4 (boot info acknowledge) [also see sub functions
	// 8 through 12 for add’l modifiers ]                4
	// Write parameter 5 (boot flags)                    5
	// Write parameter 6 (initiator info)                6
	// Write parameter 7 (initiator mailbox)             7
	// Write “OEM has handled boot info” bit             8
	// Write “SMS has handled boot info.” bit            9
	// Write “OS / service partition has handled boot info.” bit.
	//                                                   10
	// Write “OS Loader has handled boot info.” bit.     11
	// Write “BIOS/POST has handled boot info.” bit.     12
	GetSystemBootOptions = commands.CommandCode{Name: "Get System Boot Options", NetworkFunction: commands.NetworkFunctionChassis, Code: 0x09, PrivilegeLevel: commands.PrivLevelOperator}
	// unassigned= commands.CommandCode{Name: "unassigned", NetworkFunction: commands.NetworkFunctionChassis, Code:  0x0Ch-0E}
	GetPOHCounter = commands.CommandCode{Name: "Get POH Counter", NetworkFunction: commands.NetworkFunctionChassis, Code: 0x0F, PrivilegeLevel: commands.PrivLevelUser}

	// Event Commands
	SetEventReceiver = commands.CommandCode{Name: "Set Event Receiver", NetworkFunction: commands.NetworkFunctionSensor, Code: 0x00, PrivilegeLevel: commands.PrivLevelAdmin}
	GetEventReceiver = commands.CommandCode{Name: "Get Event Receiver", NetworkFunction: commands.NetworkFunctionSensor, Code: 0x01, PrivilegeLevel: commands.PrivLevelUser}
	PlatformEvent    = commands.CommandCode{Name: "Platform Event (a.k.a. 'Event Message')", NetworkFunction: commands.NetworkFunctionSensor, Code: 0x02, PrivilegeLevel: commands.PrivLevelOperator}
	// unassigned= commands.CommandCode{Name: "unassigned", NetworkFunction: commands.NetworkFunctionSensor, Code:  0x03h-0F}
	// PEF and Alerting Commands

	// PEF and Alerting Commands

	GetPEFCapabilities  = commands.CommandCode{Name: "Get PEF Capabilities", NetworkFunction: commands.NetworkFunctionSensor, Code: 0x10, PrivilegeLevel: commands.PrivLevelUser}
	ArmPEFPostponeTimer = commands.CommandCode{Name: "Arm PEF Postpone Timer", NetworkFunction: commands.NetworkFunctionSensor, Code: 0x11, PrivilegeLevel: commands.PrivLevelAdmin}

	// Disable Postpone Timer     0
	// Arm Timer                  1
	// Temporary PEF disable      2

	SetPEFConfigurationParameters = commands.CommandCode{Name: "Set PEF Configuration Parameters", NetworkFunction: commands.NetworkFunctionSensor, Code: 0x12, PrivilegeLevel: commands.PrivLevelAdmin}
	// Write parameter 1 (PEF control)                       0
	// Write parameter 2 (PEF Action global control)         1
	// Write parameter 3 (PEF Startup Delay)                 2
	// Write parameter 4 (PEF Alert Startup Delay)           3
	// Write parameter 6 (Event Filter Table)                4
	// Write parameter 7 (Event Filter Table Data 1)         5
	// Write parameter 9 (Alert Policy Table)                6
	// Write parameter 10 (System GUID)                      7
	// Write parameter 12 (Alert String Keys) volatile       8
	// Write parameter 12 (Alert String Keys)-non-volatile   9
	// Write parameter 13 (Alert Strings)-volatile           10
	// Write parameter 13 (Alert Strings)-non-volatile       11
	// Write parameter 15 (Group Control Table)-non-volatile 12
	// Write OEM parameters                                  13
	GetPEFConfigurationParameters = commands.CommandCode{Name: "Get PEF Configuration Parameters", NetworkFunction: commands.NetworkFunctionSensor, Code: 0x13, PrivilegeLevel: commands.PrivLevelOperator}
	SetLastProcessedEventID       = commands.CommandCode{Name: "Set Last Processed Event ID", NetworkFunction: commands.NetworkFunctionSensor, Code: 0x14, PrivilegeLevel: commands.PrivLevelAdmin}
	GetLastProcessedEventID       = commands.CommandCode{Name: "Get Last Processed Event ID", NetworkFunction: commands.NetworkFunctionSensor, Code: 0x15, PrivilegeLevel: commands.PrivLevelAdmin}
	AlertImmediate                = commands.CommandCode{Name: "Alert Immediate", NetworkFunction: commands.NetworkFunctionSensor, Code: 0x16, PrivilegeLevel: commands.PrivLevelAdmin}
	// reserved / unspecified 0
	// Alert to Channel 1
	// Alert to Channel 2
	// Alert to Channel 3
	// Alert to Channel 4
	// Alert to Channel 5
	// Alert to Channel 6
	// Alert to Channel 7
	// Platform Event Parameters  8
	// Alert to Channel 8         9
	// Alert to Channel 9         10
	// Alert to Channel Ah        11
	// Alert to Channel Bh        12
	PETAcknowledge = commands.CommandCode{Name: "PET Acknowledge", NetworkFunction: commands.NetworkFunctionSensor, Code: 0x17, PrivilegeLevel: commands.PrivLevelUnprotected}

	// Sensor Device Commands

	GetDeviceSDRInfo               = commands.CommandCode{Name: "Get Device SDR Info", NetworkFunction: commands.NetworkFunctionSensor, Code: 0x20}           // Local only
	GetDeviceSDR                   = commands.CommandCode{Name: "Get Device SDR", NetworkFunction: commands.NetworkFunctionSensor, Code: 0x21}                // Local only
	ReserveDeviceSDRRepository     = commands.CommandCode{Name: "Reserve Device SDR Repository", NetworkFunction: commands.NetworkFunctionSensor, Code: 0x22} // Local only
	GetSensorReadingFactors        = commands.CommandCode{Name: "Get Sensor Reading Factors", NetworkFunction: commands.NetworkFunctionSensor, Code: 0x23, PrivilegeLevel: commands.PrivLevelUser}
	SetSensorHysteresis            = commands.CommandCode{Name: "Set Sensor Hysteresis", NetworkFunction: commands.NetworkFunctionSensor, Code: 0x24, PrivilegeLevel: commands.PrivLevelOperator}
	GetSensorHysteresis            = commands.CommandCode{Name: "Get Sensor Hysteresis", NetworkFunction: commands.NetworkFunctionSensor, Code: 0x25, PrivilegeLevel: commands.PrivLevelUser}
	SetSensorThresholds            = commands.CommandCode{Name: "Set Sensor Threshold", NetworkFunction: commands.NetworkFunctionSensor, Code: 0x26, PrivilegeLevel: commands.PrivLevelOperator}
	GetSensorThresholds            = commands.CommandCode{Name: "Get Sensor Threshold", NetworkFunction: commands.NetworkFunctionSensor, Code: 0x27, PrivilegeLevel: commands.PrivLevelUser}
	SetSensorEventEnable           = commands.CommandCode{Name: "Set Sensor Event Enable", NetworkFunction: commands.NetworkFunctionSensor, Code: 0x28, PrivilegeLevel: commands.PrivLevelOperator}
	GetSensorEventEnable           = commands.CommandCode{Name: "Get Sensor Event Enable", NetworkFunction: commands.NetworkFunctionSensor, Code: 0x29, PrivilegeLevel: commands.PrivLevelUser}
	ReArmSensorEvents              = commands.CommandCode{Name: "Re-arm Sensor Events", NetworkFunction: commands.NetworkFunctionSensor, Code: 0x2A, PrivilegeLevel: commands.PrivLevelOperator}
	GetSensorEventStatus           = commands.CommandCode{Name: "Get Sensor Event Status", NetworkFunction: commands.NetworkFunctionSensor, Code: 0x2B, PrivilegeLevel: commands.PrivLevelUser}
	GetSensorReading               = commands.CommandCode{Name: "Get Sensor Reading", NetworkFunction: commands.NetworkFunctionSensor, Code: 0x2D, PrivilegeLevel: commands.PrivLevelUser}
	SetSensorType                  = commands.CommandCode{Name: "Set Sensor Type", NetworkFunction: commands.NetworkFunctionSensor, Code: 0x2E, PrivilegeLevel: commands.PrivLevelOperator}
	GetSensorType                  = commands.CommandCode{Name: "Get Sensor Type", NetworkFunction: commands.NetworkFunctionSensor, Code: 0x2F, PrivilegeLevel: commands.PrivLevelUser}
	SetSensorReadingAndEventStatus = commands.CommandCode{Name: "Set Sensor Reading And Event Status", NetworkFunction: commands.NetworkFunctionSensor, Code: 0x30, PrivilegeLevel: commands.PrivLevelOperator}

	// FRU Device Commands
	GetFRUInventoryAreaInfo = commands.CommandCode{Name: "Get FRU Inventory Area Info", NetworkFunction: commands.NetworkFunctionStorage, Code: 0x10, PrivilegeLevel: commands.PrivLevelUser}
	ReadFRUData             = commands.CommandCode{Name: "Read FRU Data", NetworkFunction: commands.NetworkFunctionStorage, Code: 0x11, PrivilegeLevel: commands.PrivLevelUser}
	WriteFRUData            = commands.CommandCode{Name: "Write FRU Data", NetworkFunction: commands.NetworkFunctionStorage, Code: 0x12, PrivilegeLevel: commands.PrivLevelOperator}

	// SDR Device Commands
	GetSDRRepositoryInfo           = commands.CommandCode{Name: "Get SDR Repository Info", NetworkFunction: commands.NetworkFunctionStorage, Code: 0x20, PrivilegeLevel: commands.PrivLevelUser}
	GetSDRRepositoryAllocationInfo = commands.CommandCode{Name: "Get SDR Repository Allocation Info", NetworkFunction: commands.NetworkFunctionStorage, Code: 0x21, PrivilegeLevel: commands.PrivLevelUser}
	ReserveSDRRepository           = commands.CommandCode{Name: "Reserve SDR Repository", NetworkFunction: commands.NetworkFunctionStorage, Code: 0x22, PrivilegeLevel: commands.PrivLevelUser}
	GetSDR                         = commands.CommandCode{Name: "Get SDR", NetworkFunction: commands.NetworkFunctionStorage, Code: 0x23, PrivilegeLevel: commands.PrivLevelUser}
	AddSDR                         = commands.CommandCode{Name: "Add SDR", NetworkFunction: commands.NetworkFunctionStorage, Code: 0x24, PrivilegeLevel: commands.PrivLevelOperator}
	PartialAddSDR                  = commands.CommandCode{Name: "Partial Add SDR", NetworkFunction: commands.NetworkFunctionStorage, Code: 0x25, PrivilegeLevel: commands.PrivLevelOperator}
	DeleteSDR                      = commands.CommandCode{Name: "Delete SDR", NetworkFunction: commands.NetworkFunctionStorage, Code: 0x26, PrivilegeLevel: commands.PrivLevelOperator}
	ClearSDRRepository             = commands.CommandCode{Name: "Clear SDR Repository", NetworkFunction: commands.NetworkFunctionStorage, Code: 0x27, PrivilegeLevel: commands.PrivLevelOperator}
	GetSDRRepositoryTime           = commands.CommandCode{Name: "Get SDR Repository Time", NetworkFunction: commands.NetworkFunctionStorage, Code: 0x28, PrivilegeLevel: commands.PrivLevelUser}
	SetSDRRepositoryTime           = commands.CommandCode{Name: "Set SDR Repository Time", NetworkFunction: commands.NetworkFunctionStorage, Code: 0x29, PrivilegeLevel: commands.PrivLevelOperator}
	EnterSDRRepositoryUpdateMode   = commands.CommandCode{Name: "Enter SDR Repository Update Mode", NetworkFunction: commands.NetworkFunctionStorage, Code: 0x2A, PrivilegeLevel: commands.PrivLevelOperator}
	ExitSDRRepositoryUpdateMode    = commands.CommandCode{Name: "Exit SDR Repository Update Mode", NetworkFunction: commands.NetworkFunctionStorage, Code: 0x2B, PrivilegeLevel: commands.PrivLevelOperator}
	RunInitializationAgent         = commands.CommandCode{Name: "Run Initialization Agent", NetworkFunction: commands.NetworkFunctionStorage, Code: 0x2C, PrivilegeLevel: commands.PrivLevelOperator}

	// SEL Device Commands
	GetSELInfo            = commands.CommandCode{Name: "Get SEL Info", NetworkFunction: commands.NetworkFunctionStorage, Code: 0x40, PrivilegeLevel: commands.PrivLevelUser}
	GetSELAllocationInfo  = commands.CommandCode{Name: "Get SEL Allocation Info", NetworkFunction: commands.NetworkFunctionStorage, Code: 0x41, PrivilegeLevel: commands.PrivLevelUser}
	ReserveSEL            = commands.CommandCode{Name: "Reserve SEL", NetworkFunction: commands.NetworkFunctionStorage, Code: 0x42, PrivilegeLevel: commands.PrivLevelUser}
	GetSELEntry           = commands.CommandCode{Name: "Get SEL Entry", NetworkFunction: commands.NetworkFunctionStorage, Code: 0x43, PrivilegeLevel: commands.PrivLevelUser}
	AddSELEntry           = commands.CommandCode{Name: "Add SEL Entry", NetworkFunction: commands.NetworkFunctionStorage, Code: 0x44, PrivilegeLevel: commands.PrivLevelOperator}
	PartialAddSELEntry    = commands.CommandCode{Name: "Partial Add SEL Entry", NetworkFunction: commands.NetworkFunctionStorage, Code: 0x45, PrivilegeLevel: commands.PrivLevelOperator}
	DeleteSELEntry        = commands.CommandCode{Name: "Delete SEL Entry", NetworkFunction: commands.NetworkFunctionStorage, Code: 0x46, PrivilegeLevel: commands.PrivLevelOperator}
	ClearSEL              = commands.CommandCode{Name: "Clear SEL", NetworkFunction: commands.NetworkFunctionStorage, Code: 0x47, PrivilegeLevel: commands.PrivLevelOperator}
	GetSELTime            = commands.CommandCode{Name: "Get SEL Time", NetworkFunction: commands.NetworkFunctionStorage, Code: 0x48, PrivilegeLevel: commands.PrivLevelUser}
	SetSELTime            = commands.CommandCode{Name: "Set SEL Time", NetworkFunction: commands.NetworkFunctionStorage, Code: 0x49, PrivilegeLevel: commands.PrivLevelOperator}
	GetAuxiliaryLogStatus = commands.CommandCode{Name: "Get Auxiliary Log Status", NetworkFunction: commands.NetworkFunctionStorage, Code: 0x5A, PrivilegeLevel: commands.PrivLevelUser}
	SetAuxiliaryLogStatus = commands.CommandCode{Name: "Set Auxiliary Log Status", NetworkFunction: commands.NetworkFunctionStorage, Code: 0x5B, PrivilegeLevel: commands.PrivLevelAdmin}
	// Set MCA                              0
	// Set OEM1                              1
	// Set OEM2                              2
	GetSELTimeUTCOffset = commands.CommandCode{Name: "Get SEL Time UTC Offset", NetworkFunction: commands.NetworkFunctionStorage, Code: 0x5C, PrivilegeLevel: commands.PrivLevelUser}
	SetSELTimeUTCOffset = commands.CommandCode{Name: "Set SEL Time UTC Offset", NetworkFunction: commands.NetworkFunctionStorage, Code: 0x5D, PrivilegeLevel: commands.PrivLevelOperator}
	// // LAN Device Commands
	SetLANConfigurationParameters = commands.CommandCode{Name: "Set LAN Configuration Parameters", NetworkFunction: commands.NetworkFunctionTransport, Code: 0x01, PrivilegeLevel: commands.PrivLevelAdmin}
	// reserved / unspecified         0
	// Set for channel 1              1
	// Set for channel 2              2
	// Set for channel 3              3
	// Set for channel 4              4
	// Set for channel 5              5
	// Set for channel 6              6
	// Set for channel 7              7
	// The following sub-function enables apply across each channel for which ‘Set’ has been enabled: Write parameters 3, 4, 6, 7, 12 - 15 (IP Address, IP Address Source, Subnet Mask, IPv4 Header Parameters, Default Gateway Address, Default Gateway MAC Address, Backup Gateway Address, Backup Gateway MAC Address)
	//                                                                             8
	// Write parameter 5 (MAC Address)                                             9
	// Write parameters 8, 9 (Primary & Secondary RMCP Port)                       10
	// Write parameter 10, 11 (Gratuitous ARP control, Gratuitous ARP interval)    11
	// Write Parameter 16 (Community String)                                       12
	// Write Parameter 18 (Destination Type) -volatile                             13
	// Write Parameter 18 (Destination Type) -non-volatile                         14
	// Write Parameter 19 (Destination Addresses) -volatile                        15
	// Write Parameter 19 (Destination Addresses) -non-volatile                    16
	// Write Parameter 20 (802.1q VLAN ID)                                         17
	// Write Parameter 21 (802.1q Priority)                                        18
	// Write Parameter 24 (RMCP+ Messaging Cipher Suite Privilege Levels)          19
	// Write OEM Parameters                                                        20
	// Set for channel                                         8                   21
	// Set for channel                                         9                   22
	// Set for channel Ah                                                          23
	// Set for channel Bh                                                          24
	// Write Parameter 25 (Destination Address VLAN TAGs)                          25
	// Write Parameter 26 (Bad Password Threshold)                                 26
	GetLANConfigurationParameters = commands.CommandCode{Name: "Get LAN Configuration Parameters", NetworkFunction: commands.NetworkFunctionTransport, Code: 0x02, PrivilegeLevel: commands.PrivLevelOperator}
	SuspendBMCARPs                = commands.CommandCode{Name: "Suspend BMC ARPs", NetworkFunction: commands.NetworkFunctionTransport, Code: 0x03, PrivilegeLevel: commands.PrivLevelAdmin}
	// reserved / unspecified   0
	// ARP Response for channel 1 1
	// ARP Response for channel 2 2
	// ARP Response for channel 3 3
	// ARP Response for channel 4 4
	// ARP Response for channel 5 5
	// ARP Response for channel 6 6
	// ARP Response for channel 7 7
	// reserved / unspecified   8
	// Gratuitous ARP for channel 1 9
	// Gratuitous ARP for channel 2 10
	// Gratuitous ARP for channel 3 11
	// Gratuitous ARP for channel 4 12
	// Gratuitous ARP for channel 5 13
	// Gratuitous ARP for channel 6 14
	// Gratuitous ARP for channel 7 15
	// ARP Response for channel 8 16
	// ARP Response for channel 9 17
	// ARP Response for channel Ah 18
	// ARP Response for channel Bh 19
	// Gratuitous ARP for channel 8h 20
	// Gratuitous ARP for channel 9h 21
	// Gratuitous ARP for channel Ah 22
	// Gratuitous ARP for channel Bh 23

	GetIP_UDP_RMCPStatistics = commands.CommandCode{Name: "Get IP/UDP/RMCP Statistics", NetworkFunction: commands.NetworkFunctionTransport, Code: 0x04, PrivilegeLevel: commands.PrivLevelUser}
	// Serial/Modem Device Commands
	SetSerialModemConfiguration = commands.CommandCode{Name: "Set Serial/Modem Configuration", NetworkFunction: commands.NetworkFunctionTransport, Code: 0x10, PrivilegeLevel: commands.PrivLevelAdmin}

	// reserved / unspecified 0
	// Set for channel 1 1
	// Set for channel 2 2
	// Set for channel 3 3
	// Set for channel 4 4
	// Set for channel 5 5
	// Set for channel 6 6
	// Set for channel 7 7
	// The following sub-function enables apply across each channel for which ‘Set’ has been enabled: Write Parameter 2 (Authentication Type Enables)
	//                                                                                  8
	// Write Parameter 3 (Connection Mode)                                              9
	// Write Parameters 4 & 6 (Session Inactivity Timeout, Session Termination)        10
	// Write Parameter 5 (Channel Callback Control)                                    11
	// Write Parameter 7 (IPMI Messaging Comm Settings)                                12
	// Write Parameter 8 (Mux Switch Control)                                          13
	// Write Parameters 9, 10, 11, 12, & 13 (Modem Ring Time, Modem Init String, Modem Escape Sequence, Modem Hang-up Sequence, Modem Dial Command)
	//                                                                                 14
	// Write Parameters 14 & 18 (Page Blackout Interval, Call Retry Interval)          15
	// Write Parameter Community String 15                                             16
	// Write Parameters 17, 19, 21, 23 [Destination Info (volatile} Destination Comm Settings (volatile} Destination Dial Strings (volatile} Destination IP Addresses (volatile)]
	//                                                                                 17
	// Write Parameters 17, 19, 21, 23 [Destination Info (non-volatile} Destination Comm Settings (non-volatile} Destination Dial Strings (non-volatile} Destination IP Addresses (non-volatile)]
	//                                                                                 18
	// Write Parameters 25, 26, 27, & 28 (TAP Account, TAP Passwords, TAP Pager ID Strings, TAP Service Settings)  19
	// Write Parameter 29 (Terminal Mode Configuration)                                                            20
	// Write Parameters 30, 33, 35, 36, & 48 (PPP Protocol Options, PPP Link Authentication, PPP ACCM, PPP Snoop ACCM, PPP Remote Console IP Address)                       21
	// Write Parameters 31 & 32 (PPP Primary RMCP Port Number, PPP Secondary RMCP Port Number)                                                                              22
	// Write Parameter 34 (CHAP Name)                                                                                                                                       23
	// Write Parameter 45 (PPP UDP ProxyIP Header)                                                                                                                          24
	// Write Parameters 38-44 -volatile (Account 0) (PPP Account Dial String Selector, PPP Account IP Addresses / BMC IP Address, PPP Account User Names, PPP Account User Domains, PPP Account User Passwords, PPP Account Authentication Settings,PPP Account Connection Hold Times)
	//                                                                                                                                                                      25
	// Write Parameters 38-44 -non-volatile (Account 1) (PPP Account Dial String Selector, PPP Account IP Addresses / BMC IP Address, PPP Account User Names, PPP Account User Domains, PPP Account User Passwords, PPP Account Authentication Settings,PPP Account Connection Hold Times)
	//                                                                                                                                                                      26
	// Write Parameters 38-44 -non-volatile (Accounts 2-n) (PPP Account Dial String Selector, PPP Account IP Addresses / BMC IP Address, PPP Account User Names, PPP Account User Domains, PPP Account User Passwords, PPP Account Authentication Settings,PPP Account Connection Hold Times)
	//                                                                                                                                                                      27
	// Write Parameter 49 (System Phone Number)                                                                                                                             28
	// Write OEM Parameters                                                                                                                                                 29
	// Set for channel 8                                                                                                                                                    30
	// Set for channel 9                                                                                                                                                    31
	// Set for channel Ah                                                                                                                                                   32
	// Set for channel Bh                                                                                                                                                   33
	// Write Parameter 54 (Bad Password Threshold)                                                                                                                          34

	GetSerialModemConfiguration = commands.CommandCode{Name: "Get Serial/Modem Configuration", NetworkFunction: commands.NetworkFunctionTransport, Code: 0x11, PrivilegeLevel: commands.PrivLevelOperator}
	SetSerialModemMux           = commands.CommandCode{Name: "Set Serial/Modem Mux", NetworkFunction: commands.NetworkFunctionTransport, Code: 0x12, PrivilegeLevel: commands.PrivLevelOperator}
	// reserved / unspecified        0
	// Function 1h (request switch of mux to system)                             1
	// Function 2h (request switch of mux to BMC)                                2
	// Function 3h (force switch of mux to system)                               3
	// Function 4h (force switch of mux to BMC)                                  4
	// Function 5h (block requests to switch mux to system)                      5
	// Function 6h (allow requests to switch mux to system)                      6
	// Function 7h (block requests to switch mux to BMC)                         7
	// Function 8h (allow requests to switch mux to BMC)                         8
	GetTAPResponseCodes        = commands.CommandCode{Name: "Get TAP Response Codes", NetworkFunction: commands.NetworkFunctionTransport, Code: 0x13, PrivilegeLevel: commands.PrivLevelUser}
	SetPPPUDPProxyTransmitData = commands.CommandCode{Name: "Set PPP UDP Proxy Transmit Data", NetworkFunction: commands.NetworkFunctionTransport, Code: 0x14} // System only

	// reserved / unspecified 0
	// Set for channel 1
	// Set for channel 2
	// Set for channel 3
	// Set for channel 4
	// Set for channel 5
	// Set for channel 6
	// Set for channel 7
	// Set for channel 8
	// Set for channel 9
	// Set for channel Ah
	// Set for channel Bh
	GetPPPUDPProxyTransmitData = commands.CommandCode{Name: "Get PPP UDP Proxy Transmit Data", NetworkFunction: commands.NetworkFunctionTransport, Code: 0x15} // System only
	SendPPPUDPProxyPacket      = commands.CommandCode{Name: "Send PPP UDP Proxy Packet", NetworkFunction: commands.NetworkFunctionTransport, Code: 0x16}       // System only
	// reserved / unspecified 0
	// Send for channel 1
	// Send for channel 2
	// Send for channel 3
	// Send for channel 4
	// Send for channel 5
	// Send for channel 6
	// Send for channel 7
	// Send for channel 8
	// Send for channel 9
	// Send for channel Ah
	// Send for channel Bh
	GetPPPUDPProxyReceiveData   = commands.CommandCode{Name: "Get PPP UDP Proxy Receive Data", NetworkFunction: commands.NetworkFunctionTransport, Code: 0x17}                                                // System only
	SerialModemConnectionActive = commands.CommandCode{Name: "Serial/Modem Connection Active", NetworkFunction: commands.NetworkFunctionTransport, Code: 0x18, PrivilegeLevel: commands.PrivLevelUnprotected} // Before session
	Callback                    = commands.CommandCode{Name: "Callback", NetworkFunction: commands.NetworkFunctionTransport, Code: 0x19, PrivilegeLevel: commands.PrivLevelAdmin}
	// reserved / unspecified 0
	// Callback using channel 1 parameters 1
	// Callback using channel 2 parameters 2
	// Callback using channel 3 parameters 3
	// Callback using channel 4 parameters 4
	// Callback using channel 5 parameters 5
	// Callback using channel 6 parameters 6
	// Callback using channel 7 parameters 7
	// Callback using channel 8 parameters 8
	// Callback using channel 9 parameters 9
	// Callback using channel Ah parameters 10
	// Callback using channel Bh parameters 11
	SetUserCallbackOptions = commands.CommandCode{Name: "Set User Callback Options", NetworkFunction: commands.NetworkFunctionTransport, Code: 0x1A, PrivilegeLevel: commands.PrivLevelAdmin}
	// reserved / unspecified
	// 0
	// Set for channel 1
	// Set for channel 2
	// Set for channel 3
	// Set for channel 4
	// Set for channel 5
	// Set for channel 6
	// Set for channel 7
	// Set for channel 8
	// Set for channel 9
	// Set for channel Ah
	// Set for channel Bh

	GetUserCallbackOptions        = commands.CommandCode{Name: "Get User Callback Options", NetworkFunction: commands.NetworkFunctionTransport, Code: 0x1B, PrivilegeLevel: commands.PrivLevelUser}
	SetSerialRoutingMux           = commands.CommandCode{Name: "Set Serial Routing Mux", NetworkFunction: commands.NetworkFunctionTransport, Code: 0x1C, PrivilegeLevel: commands.PrivLevelAdmin}
	SOLActivating                 = commands.CommandCode{Name: "SOL Activating", NetworkFunction: commands.NetworkFunctionTransport, Code: 0x20} // Weird.
	SetSOLConfigurationParameters = commands.CommandCode{Name: "Set SOL Configuration Parameters", NetworkFunction: commands.NetworkFunctionTransport, Code: 0x21, PrivilegeLevel: commands.PrivLevelAdmin}

	// reserved / unspecified  0
	// Set for channel 1
	// Set for channel 2
	// Set for channel 3
	// Set for channel 4
	// Set for channel 5
	// Set for channel 6
	// Set for channel 7
	// The following sub-function enables apply across each channel for which ‘Set’ has been enabled:
	// Write Parameter 1 (SOL Enable)                               8
	// Write Parameter 2 (SOL Authentication)                       9
	// Write Parameters 3 & 4 (Character Accumulate Interval & Character Send Threshold, SOL Retry)
	//                                                              10
	// Write Parameter 5 (SOL non-volatile bit rate-non-volatile)   11
	// Write Parameter 6 (SOL volatile bit rate-volatile)           12
	// Write Parameter 8 (SOL Payload Port Number)                  13
	// Set for channel 8                                            14
	// Set for channel 9                                            15
	// Set for channel Ah                                           16
	// Set for channel Bh                                           17
	GetSOLConfigurationParameters = commands.CommandCode{Name: "Get SOL Configuration Parameters", NetworkFunction: commands.NetworkFunctionTransport, Code: 0x22, PrivilegeLevel: commands.PrivLevelUser}

	// Forwarded Commands
	// Forwarded Command (NOTE: This command is a byproduct  of the
	// Command Forwarding capability being enabled on one or more
	// channels and cannot be directly enabled/disabled via Firmware
	// Firmwall)

	ForwardedCommand        = commands.CommandCode{Name: "Forwarded Command", NetworkFunction: commands.NetworkFunctionTransport, Code: 0x30} // Weird
	SetForwardedCommands    = commands.CommandCode{Name: "Set Forwarded Commands", NetworkFunction: commands.NetworkFunctionTransport, Code: 0x31, PrivilegeLevel: commands.PrivLevelAdmin}
	GetForwardedCommands    = commands.CommandCode{Name: "Get Forwarded Commands", NetworkFunction: commands.NetworkFunctionTransport, Code: 0x32, PrivilegeLevel: commands.PrivLevelUser}
	EnableForwardedCommands = commands.CommandCode{Name: "Enable Forwarded Commands", NetworkFunction: commands.NetworkFunctionTransport, Code: 0x33, PrivilegeLevel: commands.PrivLevelAdmin}
	GetBridgeState          = commands.CommandCode{Name: "Get Bridge State", NetworkFunction: commands.NetworkFunctionBridge, Code: 0x00, PrivilegeLevel: commands.PrivLevelUser}
	SetBridgeState          = commands.CommandCode{Name: "Set Bridge State", NetworkFunction: commands.NetworkFunctionBridge, Code: 0x01, PrivilegeLevel: commands.PrivLevelOperator}
	GetICMBAddress          = commands.CommandCode{Name: "Get ICMB Address", NetworkFunction: commands.NetworkFunctionBridge, Code: 0x02, PrivilegeLevel: commands.PrivLevelUser}
	SetICMBAddress          = commands.CommandCode{Name: "Set ICMB Address", NetworkFunction: commands.NetworkFunctionBridge, Code: 0x03, PrivilegeLevel: commands.PrivLevelOperator}
	SetBridgeProxyAddress   = commands.CommandCode{Name: "Set Bridge ProxyAddress", NetworkFunction: commands.NetworkFunctionBridge, Code: 0x04, PrivilegeLevel: commands.PrivLevelOperator}
	GetBridgeStatistics     = commands.CommandCode{Name: "Get Bridge Statistics", NetworkFunction: commands.NetworkFunctionBridge, Code: 0x05, PrivilegeLevel: commands.PrivLevelUser}
	GetICMBCapabilities     = commands.CommandCode{Name: "Get ICMB Capabilities", NetworkFunction: commands.NetworkFunctionBridge, Code: 0x06, PrivilegeLevel: commands.PrivLevelUser}
	ClearBridgeStatistics   = commands.CommandCode{Name: "Clear Bridge Statistics", NetworkFunction: commands.NetworkFunctionBridge, Code: 0x08, PrivilegeLevel: commands.PrivLevelOperator}
	GetBridgeProxyAddress   = commands.CommandCode{Name: "Get Bridge Proxy Address", NetworkFunction: commands.NetworkFunctionBridge, Code: 0x09, PrivilegeLevel: commands.PrivLevelUser}
	GetICMBConnectorInfo    = commands.CommandCode{Name: "Get ICMB Connector Info", NetworkFunction: commands.NetworkFunctionBridge, Code: 0x0A, PrivilegeLevel: commands.PrivLevelUser}
	GetICMBConnectionID     = commands.CommandCode{Name: "Get ICMB Connection ID", NetworkFunction: commands.NetworkFunctionBridge, Code: 0x0B, PrivilegeLevel: commands.PrivLevelUser}
	SendICMBConnectionID    = commands.CommandCode{Name: "Send ICMB Connection ID", NetworkFunction: commands.NetworkFunctionBridge, Code: 0x0C, PrivilegeLevel: commands.PrivLevelUser}
	// Discovery Commands (ICMB)
	PrepareForDiscovery = commands.CommandCode{Name: "PrepareForDiscovery", NetworkFunction: commands.NetworkFunctionBridge, Code: 0x10, PrivilegeLevel: commands.PrivLevelOperator}
	GetAddresses        = commands.CommandCode{Name: "GetAddresses", NetworkFunction: commands.NetworkFunctionBridge, Code: 0x11, PrivilegeLevel: commands.PrivLevelUser}
	SetDiscovered       = commands.CommandCode{Name: "SetDiscovered", NetworkFunction: commands.NetworkFunctionBridge, Code: 0x12, PrivilegeLevel: commands.PrivLevelOperator}
	GetChassisDeviceId  = commands.CommandCode{Name: "GetChassisDeviceId", NetworkFunction: commands.NetworkFunctionBridge, Code: 0x13, PrivilegeLevel: commands.PrivLevelUser}
	SetChassisDeviceId  = commands.CommandCode{Name: "SetChassisDeviceId", NetworkFunction: commands.NetworkFunctionBridge, Code: 0x14, PrivilegeLevel: commands.PrivLevelOperator}
	// Bridging Commands (ICMB)
	BridgeRequest = commands.CommandCode{Name: "BridgeRequest", NetworkFunction: commands.NetworkFunctionBridge, Code: 0x20, PrivilegeLevel: commands.PrivLevelOperator}
	BridgeMessage = commands.CommandCode{Name: "BridgeMessage", NetworkFunction: commands.NetworkFunctionBridge, Code: 0x21, PrivilegeLevel: commands.PrivLevelOperator}
	// Event Commands (ICMB)
	GetEventCount          = commands.CommandCode{Name: "GetEventCount", NetworkFunction: commands.NetworkFunctionBridge, Code: 0x30, PrivilegeLevel: commands.PrivLevelUser}
	SetEventDestination    = commands.CommandCode{Name: "SetEventDestination", NetworkFunction: commands.NetworkFunctionBridge, Code: 0x31, PrivilegeLevel: commands.PrivLevelOperator}
	SetEventReceptionState = commands.CommandCode{Name: "SetEventReceptionState", NetworkFunction: commands.NetworkFunctionBridge, Code: 0x32, PrivilegeLevel: commands.PrivLevelOperator}
	SendICMBEventMessage   = commands.CommandCode{Name: "SendICMBEventMessage", NetworkFunction: commands.NetworkFunctionBridge, Code: 0x33, PrivilegeLevel: commands.PrivLevelOperator}
	GetEventDestination    = commands.CommandCode{Name: "GetEventDestination (optional)", NetworkFunction: commands.NetworkFunctionBridge, Code: 0x34, PrivilegeLevel: commands.PrivLevelUser}
	GetEventReceptionState = commands.CommandCode{Name: "GetEventReceptionState (optional)", NetworkFunction: commands.NetworkFunctionBridge, Code: 0x35, PrivilegeLevel: commands.PrivLevelUser}
	// OEM commands for PICMG
	/** [IPMI2] Section 5.1, table 5-1, page 41, "Group Extension" */
	PICMGExtension         = commands.CommandCode{Name: "PICMG Non-IPMI Command", NetworkFunction: commands.NetworkFunctionGroupExtension, Code: 0x00, PrivilegeLevel: commands.PrivLevelUnprotected}
	DMTFExtension          = commands.CommandCode{Name: "DMTF Non-IPMI Command", NetworkFunction: commands.NetworkFunctionGroupExtension, Code: 0x01, PrivilegeLevel: commands.PrivLevelUnprotected}
	SSIForumExtension      = commands.CommandCode{Name: "SSI Forum Non-IPMI Command", NetworkFunction: commands.NetworkFunctionGroupExtension, Code: 0x02, PrivilegeLevel: commands.PrivLevelUnprotected}
	VITAStandardsExtension = commands.CommandCode{Name: "VITA Standards Organization Non-IPMI Command", NetworkFunction: commands.NetworkFunctionGroupExtension, Code: 0x03, PrivilegeLevel: commands.PrivLevelUnprotected}
	DCMIExtension          = commands.CommandCode{Name: "DCMI Specifications Non-IPMI Command", NetworkFunction: commands.NetworkFunctionGroupExtension, Code: 0xDC, PrivilegeLevel: commands.PrivLevelUnprotected}
)
