package goipmi

import (
	"errors"
	"math"

	"github.com/runner-mei/goipmi/protocol"
	"golang.org/x/text/encoding/charmap"
)

func init() {
	_ = []Calculable{
		&FullSensorRecord{},
	}
}

type Record interface {
	protocol.Readable

	GetHeader() SensorRecordHeader
}

type Calculable interface {
	Calc(value, length int32) (float64, error)
}

type SensorRecordHeader struct {
	// Header
	RecordId     uint16
	SdrVersion   uint8
	RecordType   uint8
	RecordLength uint8
}

type FullSensorRecord struct {
	// Header
	SensorRecordHeader

	// Key Fields
	SensorOwnerId  uint8
	SensorOwnerLUN uint8
	SensorNumber   uint8

	// Data
	EntityId               uint8
	IsPhysicalInstance     bool
	EntityInstance         uint8
	SensorInitialization   uint8
	SensorCapablilities    uint8
	SensorType             uint8
	EventOrReadingTypeCode uint8
	// Assertion Event mask / Lower threshold Reading Mask
	// Dessertion Event mask / Upper threshold Reading Mask
	// Discrate Reading mask / Settable threshold Mask, Readable threshold Mask
	Masks                             [6]uint8
	SensorUnits1                      uint8
	SensorUnits2                      uint8
	SensorUnits3                      uint8
	Linearization                     uint8
	M                                 int16
	Tolerance                         uint8
	B                                 int16
	Accuracy                          int16
	AccuracyExp                       uint8
	Direction                         uint8
	Rexp                              int8
	Bexp                              int8
	AnlogcharacteristicFlags          uint8
	NominalReading                    uint8
	NominalMaximum                    uint8
	NominalMinimum                    uint8
	MaximumReading                    uint8
	MinimumReading                    uint8
	UpperNonrecoverableThreshold      uint8
	UpperCriticalThreshold            uint8
	UpperNonCriticalThreshold         uint8
	LowerNonrecoverableThreshold      uint8
	LowerCriticalThreshold            uint8
	LowerNonCriticalThreshold         uint8
	Positive_ThresholdHysteresisValue uint8
	Negative_ThresholdHysteresisValue uint8
	Reserve1                          uint8
	Reserve2                          uint8
	Oem                               uint8
	IdTypeLength                      uint8
	IdString                          string
}

func (self *FullSensorRecord) GetHeader() SensorRecordHeader {
	return self.SensorRecordHeader
}

func (self *FullSensorRecord) CanIgnore() bool {
	// fmt.Printf("%x %x %x %x %x %x %x %x %x %x %x \r\n",
	// 	self.RecordId,
	// 	self.SdrVersion,
	// 	self.RecordType,
	// 	self.RecordLength,
	// 	self.SensorOwnerId,
	// 	self.SensorOwnerLUN,
	// 	self.SensorNumber,
	// 	self.EntityId,
	// 	self.EntityInstance,
	// 	self.SensorInitialization,
	// 	self.SensorCapablilities)
	return self.SensorCapablilities&0x80 != 0
}

func (self *FullSensorRecord) ReadBytes(r *protocol.Reader) {

	// Header
	self.RecordId = r.ReadUint16()    // 1-2
	self.SdrVersion = r.ReadUint8()   // 3
	self.RecordType = r.ReadUint8()   // 4
	self.RecordLength = r.ReadUint8() // 5

	// Key Fields
	self.SensorOwnerId = r.ReadUint8()  // 6
	self.SensorOwnerLUN = r.ReadUint8() // 7
	self.SensorNumber = r.ReadUint8()   // 8

	// Data
	self.EntityId = r.ReadUint8()       // 9
	self.EntityInstance = r.ReadUint8() // 10
	self.IsPhysicalInstance = 0 == int(self.EntityInstance>>7)

	self.SensorInitialization = r.ReadUint8()   // 11
	self.SensorCapablilities = r.ReadUint8()    // 12
	self.SensorType = r.ReadUint8()             // 13
	self.EventOrReadingTypeCode = r.ReadUint8() // 14

	self.Masks[0] = r.ReadUint8() // 15
	self.Masks[1] = r.ReadUint8() // 16
	self.Masks[2] = r.ReadUint8() // 17
	self.Masks[3] = r.ReadUint8() // 18
	self.Masks[4] = r.ReadUint8() // 19
	self.Masks[5] = r.ReadUint8() // 20

	self.SensorUnits1 = r.ReadUint8()  // 21
	self.SensorUnits2 = r.ReadUint8()  // 22
	self.SensorUnits3 = r.ReadUint8()  // 23
	self.Linearization = r.ReadUint8() // 24

	calcM := r.ReadUint8()     // 25
	tolerance := r.ReadUint8() // 26
	calcB := r.ReadUint8()     // 27
	accuracy := r.ReadUint8()  // 28
	direction := r.ReadUint8() // 29
	rexpBexp := r.ReadUint8()  // 30

	self.M = decode2sComplement((int16(tolerance&uint8(0xc0))<<2)|int16(calcM), 9)
	self.Tolerance = uint8(tolerance) & 0x3f

	self.B = decode2sComplement((int16(accuracy&uint8(0xc0))<<2)|int16(calcB), 9)
	self.Accuracy = int16(uint8(accuracy)&0x3f) | int16(uint8(direction)&0xf0)<<2
	self.AccuracyExp = (uint8(direction) & 0x0f) >> 2
	self.Direction = (uint8(direction) & 0x03)

	self.Rexp = int8(decode2sComplement(int16(rexpBexp&uint8(0xf0))>>4, 3))
	self.Bexp = int8(decode2sComplement(int16(rexpBexp&uint8(0x0f)), 3))

	self.AnlogcharacteristicFlags = r.ReadUint8()          // 31
	self.NominalReading = r.ReadUint8()                    // 32
	self.NominalMaximum = r.ReadUint8()                    // 33
	self.NominalMinimum = r.ReadUint8()                    // 34
	self.MaximumReading = r.ReadUint8()                    // 35
	self.MinimumReading = r.ReadUint8()                    // 36
	self.UpperNonrecoverableThreshold = r.ReadUint8()      // 37
	self.UpperCriticalThreshold = r.ReadUint8()            // 38
	self.UpperNonCriticalThreshold = r.ReadUint8()         // 39
	self.LowerNonrecoverableThreshold = r.ReadUint8()      // 40
	self.LowerCriticalThreshold = r.ReadUint8()            // 41
	self.LowerNonCriticalThreshold = r.ReadUint8()         // 42
	self.Positive_ThresholdHysteresisValue = r.ReadUint8() // 43
	self.Negative_ThresholdHysteresisValue = r.ReadUint8() // 44
	self.Reserve1 = r.ReadUint8()                          // 45
	self.Reserve2 = r.ReadUint8()                          // 46
	self.Oem = r.ReadUint8()                               // 47
	self.IdTypeLength = r.ReadUint8()                      // 48
	self.IdString = decodeName(self.IdTypeLength, r.ReadBytes(r.Len()))
}

func (self *FullSensorRecord) Calc(value, length int32) (float64, error) {
	return calcFormula(value, length, self.SensorUnits1, self.M, self.B, int16(self.Rexp), self.Linearization)
}

type CompactSensorRecord struct {
	// Header
	SensorRecordHeader

	// Key Fields
	SensorOwnerId  uint8
	SensorOwnerLUN uint8
	SensorNumber   uint8

	// Data
	EntityId               uint8
	IsPhysicalInstance     bool
	EntityInstance         uint8
	Initialization         uint8
	Capablilities          uint8
	Type                   uint8
	EventOrReadingTypeCode uint8
	// Assertion Event mask / Lower threshold Reading Mask
	// Dessertion Event mask / Upper threshold Reading Mask
	// Discrate Reading mask / Settable threshold Mask, Readable threshold Mask
	Masks                             [6]uint8
	SensorUnits1                      uint8
	SensorUnits2                      uint8
	SensorUnits3                      uint8
	SensorDirection                   uint8
	IdStringInstanceModifierType      uint8
	SensorRecordSharing               uint8
	EntityInstanceSharing             bool
	IdStringInstanceModifierOffset    uint8
	Positive_ThresholdHysteresisValue uint8
	Negative_ThresholdHysteresisValue uint8
	Reserve1                          uint8
	Reserve2                          uint8
	Reserve3                          uint8
	Oem                               uint8
	IdTypeLength                      uint8
	IdString                          string
}

func (self *CompactSensorRecord) GetHeader() SensorRecordHeader {
	return self.SensorRecordHeader
}

func (self *CompactSensorRecord) ReadBytes(r *protocol.Reader) {

	// Header
	self.RecordId = r.ReadUint16()    // 1-2
	self.SdrVersion = r.ReadUint8()   // 3
	self.RecordType = r.ReadUint8()   // 4
	self.RecordLength = r.ReadUint8() // 5

	// Key Fields
	self.SensorOwnerId = r.ReadUint8()  // 6
	self.SensorOwnerLUN = r.ReadUint8() // 7
	self.SensorNumber = r.ReadUint8()   // 8

	// Data
	self.EntityId = r.ReadUint8()       // 9
	self.EntityInstance = r.ReadUint8() // 10
	self.IsPhysicalInstance = 0 == int(self.EntityInstance>>7)

	self.Initialization = r.ReadUint8()         // 11
	self.Capablilities = r.ReadUint8()          // 12
	self.Type = r.ReadUint8()                   // 13
	self.EventOrReadingTypeCode = r.ReadUint8() // 14

	self.Masks[0] = r.ReadUint8() // 15
	self.Masks[1] = r.ReadUint8() // 16
	self.Masks[2] = r.ReadUint8() // 17
	self.Masks[3] = r.ReadUint8() // 18
	self.Masks[4] = r.ReadUint8() // 19
	self.Masks[5] = r.ReadUint8() // 20

	self.SensorUnits1 = r.ReadUint8() // 21
	self.SensorUnits2 = r.ReadUint8() // 22
	self.SensorUnits3 = r.ReadUint8() // 23
	recordSharing1 := r.ReadUint8()   // 24
	recordSharing2 := r.ReadUint8()   // 25

	self.SensorDirection = recordSharing1 >> 6
	self.IdStringInstanceModifierType = (recordSharing1 >> 4) & 0x03
	self.SensorRecordSharing = recordSharing1 & 0x0F

	self.EntityInstanceSharing = 0 == (recordSharing2 & 0x80)
	self.IdStringInstanceModifierOffset = (recordSharing2 & 0x7F)

	self.Positive_ThresholdHysteresisValue = r.ReadUint8() // 26
	self.Negative_ThresholdHysteresisValue = r.ReadUint8() // 27
	self.Reserve1 = r.ReadUint8()                          // 28
	self.Reserve2 = r.ReadUint8()                          // 29
	self.Reserve3 = r.ReadUint8()                          // 30
	self.Oem = r.ReadUint8()                               // 31
	self.IdTypeLength = r.ReadUint8()                      // 32
	self.IdString = decodeName(self.IdTypeLength, r.ReadBytes(r.Len()))
}

type EventOnlyRecord struct {
	// Header
	SensorRecordHeader

	// Key Fields
	SensorOwnerId  uint8
	SensorOwnerLUN uint8
	SensorNumber   uint8

	// Data
	EntityId               uint8
	IsPhysicalInstance     bool
	EntityInstance         uint8
	Type                   uint8
	EventOrReadingTypeCode uint8

	SensorDirection                uint8
	IdStringInstanceModifierType   uint8
	SensorRecordSharing            uint8
	EntityInstanceSharing          bool
	IdStringInstanceModifierOffset uint8

	Reserve1     uint8
	Oem          uint8
	IdTypeLength uint8
	IdString     string
}

func (self *EventOnlyRecord) GetHeader() SensorRecordHeader {
	return self.SensorRecordHeader
}

func (self *EventOnlyRecord) ReadBytes(r *protocol.Reader) {

	// Header
	self.RecordId = r.ReadUint16()    // 1-2
	self.SdrVersion = r.ReadUint8()   // 3
	self.RecordType = r.ReadUint8()   // 4
	self.RecordLength = r.ReadUint8() // 5

	// Key Fields
	self.SensorOwnerId = r.ReadUint8()  // 6
	self.SensorOwnerLUN = r.ReadUint8() // 7
	self.SensorNumber = r.ReadUint8()   // 8

	// Data
	self.EntityId = r.ReadUint8()       // 9
	self.EntityInstance = r.ReadUint8() // 10
	self.IsPhysicalInstance = 0 == int(self.EntityInstance>>7)

	self.Type = r.ReadUint8()                   // 11
	self.EventOrReadingTypeCode = r.ReadUint8() // 12
	recordSharing1 := r.ReadUint8()             // 13
	recordSharing2 := r.ReadUint8()             // 14

	self.SensorDirection = recordSharing1 >> 6
	self.IdStringInstanceModifierType = (recordSharing1 >> 4) & 0x03
	self.SensorRecordSharing = recordSharing1 & 0x0F

	self.EntityInstanceSharing = 0 == (recordSharing2 & 0x80)
	self.IdStringInstanceModifierOffset = recordSharing2 & 0x7F

	self.Reserve1 = r.ReadUint8()     // 15
	self.Oem = r.ReadUint8()          // 16
	self.IdTypeLength = r.ReadUint8() // 17
	self.IdString = decodeName(self.IdTypeLength, r.ReadBytes(r.Len()))
}

type ListOrRangeType uint8

const (
	AsList  ListOrRangeType = 0
	AsRange ListOrRangeType = 1
)

type EntityAssociationRecord struct {
	// Header
	SensorRecordHeader

	// Data
	ContainerEntityId       uint8
	ContainerEntityInstance uint8
	AsListOrRange           ListOrRangeType
	RecordLink              bool
	EntityAccessible        bool
	Flags                   uint8

	ContainedEntity1 uint8
	ContainedEntity2 uint8
	ContainedEntity3 uint8
	ContainedEntity4 uint8

	Instance1InEntity uint8
	Instance2InEntity uint8
	Instance3InEntity uint8
	Instance4InEntity uint8

	InstanceRange1Begin uint8
	InstanceRange1End   uint8
	InstanceRange2Begin uint8
	InstanceRange2End   uint8
}

func (self *EntityAssociationRecord) GetHeader() SensorRecordHeader {
	return self.SensorRecordHeader
}

func (self *EntityAssociationRecord) ReadBytes(r *protocol.Reader) {
	// Header
	self.RecordId = r.ReadUint16()    // 1-2
	self.SdrVersion = r.ReadUint8()   // 3
	self.RecordType = r.ReadUint8()   // 4
	self.RecordLength = r.ReadUint8() // 5

	// Data
	self.ContainerEntityId = r.ReadUint8()       // 6
	self.ContainerEntityInstance = r.ReadUint8() // 7
	self.Flags = r.ReadUint8()                   // 8

	self.AsListOrRange = ListOrRangeType((self.Flags & 0x80) >> 7)
	self.RecordLink = 0 == self.Flags&0x40
	self.EntityAccessible = 0 == self.Flags&0x20

	if self.AsListOrRange == 0 {
		self.ContainedEntity1 = r.ReadUint8()  // 9
		self.Instance1InEntity = r.ReadUint8() // 10
		self.ContainedEntity2 = r.ReadUint8()  // 11
		self.Instance2InEntity = r.ReadUint8() // 12
		self.ContainedEntity3 = r.ReadUint8()  // 13
		self.Instance3InEntity = r.ReadUint8() // 14
		self.ContainedEntity4 = r.ReadUint8()  // 15
		self.Instance4InEntity = r.ReadUint8() // 16
	} else {
		self.ContainedEntity1 = r.ReadUint8()    // 9
		self.InstanceRange1Begin = r.ReadUint8() // 10
		self.ContainedEntity2 = r.ReadUint8()    // 11
		self.InstanceRange1End = r.ReadUint8()   // 12
		self.ContainedEntity3 = r.ReadUint8()    // 13
		self.InstanceRange2Begin = r.ReadUint8() // 14
		self.ContainedEntity4 = r.ReadUint8()    // 15
		self.InstanceRange2End = r.ReadUint8()   // 16
	}
}

type DeviceRelativeAssociationRecord struct {
	// Header
	SensorRecordHeader

	// Data
	ContainerEntityId            uint8
	ContainerEntityInstance      uint8
	ContainerEntityDeviceAddress uint8
	ContainerEntityDeviceChannel uint8
	AsListOrRange                ListOrRangeType
	RecordLink                   bool
	EntityAccessible             bool
	Flags                        uint8

	ContainedEntity1DeviceAddress uint8
	ContainedEntity1DeviceChannel uint8
	ContainedEntity1              uint8
	ContainedEntity2DeviceAddress uint8
	ContainedEntity2DeviceChannel uint8
	ContainedEntity2              uint8
	ContainedEntity3DeviceAddress uint8
	ContainedEntity3DeviceChannel uint8
	ContainedEntity3              uint8
	ContainedEntity4DeviceAddress uint8
	ContainedEntity4DeviceChannel uint8
	ContainedEntity4              uint8

	Instance1InEntity uint8
	Instance2InEntity uint8
	Instance3InEntity uint8
	Instance4InEntity uint8

	InstanceRange1Begin uint8
	InstanceRange1End   uint8
	InstanceRange2Begin uint8
	InstanceRange2End   uint8
}

func (self *DeviceRelativeAssociationRecord) GetHeader() SensorRecordHeader {
	return self.SensorRecordHeader
}

func (self *DeviceRelativeAssociationRecord) ReadBytes(r *protocol.Reader) {
	// Header
	self.RecordId = r.ReadUint16()    // 1-2
	self.SdrVersion = r.ReadUint8()   // 3
	self.RecordType = r.ReadUint8()   // 4
	self.RecordLength = r.ReadUint8() // 5

	// Data
	self.ContainerEntityId = r.ReadUint8()                 // 6
	self.ContainerEntityInstance = r.ReadUint8()           // 7
	self.ContainerEntityDeviceAddress = r.ReadUint8() >> 1 // 8
	self.ContainerEntityDeviceChannel = r.ReadUint8() >> 4 // 9
	self.Flags = r.ReadUint8()                             // 10

	self.AsListOrRange = ListOrRangeType((self.Flags & 0x80) >> 7)
	self.RecordLink = 0 == self.Flags&0x40
	self.EntityAccessible = 0 == self.Flags&0x20

	if self.AsListOrRange == AsList {
		self.ContainedEntity1DeviceAddress = r.ReadUint8() >> 1 // 11
		self.ContainedEntity1DeviceChannel = r.ReadUint8() >> 4 // 12
		self.ContainedEntity1 = r.ReadUint8()                   // 13
		self.Instance1InEntity = r.ReadUint8()                  // 14

		self.ContainedEntity2DeviceAddress = r.ReadUint8() >> 1 // 15
		self.ContainedEntity2DeviceChannel = r.ReadUint8() >> 4 // 16
		self.ContainedEntity2 = r.ReadUint8()                   // 17
		self.Instance2InEntity = r.ReadUint8()                  // 18

		self.ContainedEntity3DeviceAddress = r.ReadUint8() >> 1 // 19
		self.ContainedEntity3DeviceChannel = r.ReadUint8() >> 4 // 20
		self.ContainedEntity3 = r.ReadUint8()                   // 21
		self.Instance3InEntity = r.ReadUint8()                  // 22

		self.ContainedEntity4DeviceAddress = r.ReadUint8() >> 1 // 23
		self.ContainedEntity4DeviceChannel = r.ReadUint8() >> 4 // 24
		self.ContainedEntity4 = r.ReadUint8()                   // 25
		self.Instance4InEntity = r.ReadUint8()                  // 26
	} else {
		self.ContainedEntity1DeviceAddress = r.ReadUint8() >> 1 // 11
		self.ContainedEntity1DeviceChannel = r.ReadUint8() >> 4 // 12
		self.ContainedEntity1 = r.ReadUint8()                   // 13
		self.InstanceRange1Begin = r.ReadUint8()                // 14

		self.ContainedEntity2DeviceAddress = r.ReadUint8() >> 1 // 15
		self.ContainedEntity2DeviceChannel = r.ReadUint8() >> 4 // 16
		self.ContainedEntity2 = r.ReadUint8()                   // 17
		self.InstanceRange1End = r.ReadUint8()                  // 18

		self.ContainedEntity3DeviceAddress = r.ReadUint8() >> 1 // 19
		self.ContainedEntity3DeviceChannel = r.ReadUint8() >> 4 // 20
		self.ContainedEntity3 = r.ReadUint8()                   // 21
		self.InstanceRange2Begin = r.ReadUint8()                // 22

		self.ContainedEntity4DeviceAddress = r.ReadUint8() >> 1 // 23
		self.ContainedEntity4DeviceChannel = r.ReadUint8() >> 4 // 24
		self.ContainedEntity4 = r.ReadUint8()                   // 25
		self.InstanceRange2End = r.ReadUint8()                  // 36
	}
}

type GenericDeviceLocatorRecord struct {
	// Header
	SensorRecordHeader
	// Key Fields
	DeviceAccessAddress uint8
	DeviceSlaveAddress  uint8
	AccessLUN           uint8
	AccessCommand       uint8
	AccessBus           uint8

	// Data
	AddressSpan        uint8
	Reserved           uint8
	DeviceType         uint8
	DeviceTypeModifier uint8
	EntityID           uint8
	EntityIntance      uint8
	Oem                uint8
	IdTypeLength       uint8
	IdString           string
}

func (self *GenericDeviceLocatorRecord) GetHeader() SensorRecordHeader {
	return self.SensorRecordHeader
}

func (self *GenericDeviceLocatorRecord) ReadBytes(r *protocol.Reader) {
	// Header
	self.RecordId = r.ReadUint16()    // 1-2
	self.SdrVersion = r.ReadUint8()   // 3
	self.RecordType = r.ReadUint8()   // 4
	self.RecordLength = r.ReadUint8() // 5

	// Key Fields
	self.DeviceAccessAddress = r.ReadUint8() >> 1 // 6
	self.DeviceSlaveAddress = r.ReadUint8() >> 1  // 7
	lunAndBus := r.ReadUint8()                    // 8

	self.AccessLUN = lunAndBus >> 5
	self.AccessCommand = (lunAndBus >> 3) & 0x03
	self.AccessBus = lunAndBus & 0x07

	// Data
	self.AddressSpan = r.ReadUint8() & 0x07 // 9
	self.Reserved = r.ReadUint8()           // 10
	self.DeviceType = r.ReadUint8()         // 11
	self.DeviceTypeModifier = r.ReadUint8() // 12
	self.EntityID = r.ReadUint8()           // 13
	self.EntityIntance = r.ReadUint8()      // 14
	self.Oem = r.ReadUint8()                // 15
	self.IdTypeLength = r.ReadUint8()       // 16
	self.IdString = decodeName(self.IdTypeLength, r.ReadBytes(r.Len()))
}

type FruDeviceLocatorRecord struct {
	// Header
	SensorRecordHeader
	// Key Fields
	DeviceAccessAddress             uint8
	FruDeviceIDOrDeviceSlaveAddress uint8
	IsLogical                       bool
	AccessLUN                       uint8
	AccessBus                       uint8
	ChannelNumber                   uint8

	// Data
	Reserved           uint8
	DeviceType         uint8
	DeviceTypeModifier uint8
	FruEntityID        uint8
	FruEntityIntance   uint8
	Oem                uint8
	IdTypeLength       uint8
	IdString           string
}

func (self *FruDeviceLocatorRecord) GetHeader() SensorRecordHeader {
	return self.SensorRecordHeader
}

func (self *FruDeviceLocatorRecord) ReadBytes(r *protocol.Reader) {
	// Header
	self.RecordId = r.ReadUint16()    // 1-2
	self.SdrVersion = r.ReadUint8()   // 3
	self.RecordType = r.ReadUint8()   // 4
	self.RecordLength = r.ReadUint8() // 5

	// Key Fields
	self.DeviceAccessAddress = r.ReadUint8() >> 1        // 6
	self.FruDeviceIDOrDeviceSlaveAddress = r.ReadUint8() // 7
	lunAndBus := r.ReadUint8()                           // 8
	self.ChannelNumber = r.ReadUint8() >> 4              // 9

	self.IsLogical = 1 == (lunAndBus & 0x80)
	self.AccessLUN = (lunAndBus & 0x18) >> 3
	self.AccessBus = lunAndBus & 0x07

	// Data
	self.Reserved = r.ReadUint8()           // 10
	self.DeviceType = r.ReadUint8()         // 11
	self.DeviceTypeModifier = r.ReadUint8() // 12
	self.FruEntityID = r.ReadUint8()        // 13
	self.FruEntityIntance = r.ReadUint8()   // 14
	self.Oem = r.ReadUint8()                // 15
	self.IdTypeLength = r.ReadUint8()       // 16
	self.IdString = decodeName(self.IdTypeLength, r.ReadBytes(r.Len()))
}

type McDeviceLocatorRecord struct {
	// Header
	SensorRecordHeader
	// Key Fields
	DeviceSlaveAddress uint8
	ChannelNumber      uint8

	// Data
	PowerStateNotificationAndGlobalInitialization uint8
	// ACPIDevicePowerStateNotificationRequired        bool
	// ACPISystemPowerStateNotificationRequired        bool
	// ControlerLogsInitializationAgentErrors          bool //  这个是 agent 出错
	// LogsInitializationAgentErrorsAccessingControler bool //  这个是 agent 访问 controler 出错

	DeviceCapablilities uint8
	Reserved1           uint8
	Reserved2           uint8
	Reserved3           uint8
	EntityID            uint8
	EntityIntance       uint8
	Oem                 uint8
	IdTypeLength        uint8
	IdString            string
}

func (self *McDeviceLocatorRecord) GetHeader() SensorRecordHeader {
	return self.SensorRecordHeader
}

func (self *McDeviceLocatorRecord) ReadBytes(r *protocol.Reader) {
	// Header
	self.RecordId = r.ReadUint16()    // 1-2
	self.SdrVersion = r.ReadUint8()   // 3
	self.RecordType = r.ReadUint8()   // 4
	self.RecordLength = r.ReadUint8() // 5

	// Key Fields
	self.DeviceSlaveAddress = r.ReadUint8() >> 1 // 6
	self.ChannelNumber = r.ReadUint8() & 0x07    // 7

	// Data
	self.PowerStateNotificationAndGlobalInitialization = r.ReadUint8() // 8
	self.DeviceCapablilities = r.ReadUint8()                           // 9

	self.Reserved1 = r.ReadUint8()     // 10
	self.Reserved2 = r.ReadUint8()     // 11
	self.Reserved3 = r.ReadUint8()     // 12
	self.EntityID = r.ReadUint8()      // 13
	self.EntityIntance = r.ReadUint8() // 14
	self.Oem = r.ReadUint8()           // 15
	self.IdTypeLength = r.ReadUint8()  // 16
	self.IdString = decodeName(self.IdTypeLength, r.ReadBytes(r.Len()))
}

type McDeviceConfirmationRecord struct {
	// Header
	SensorRecordHeader
	// Key Fields
	DeviceSlaveAddress uint8
	DeviceID           uint8
	ChannelNumber      uint8
	ChannelRevision    uint8

	// Data
	FirmwareMajorRevision uint8
	FirmwareMinorRevision uint8
	IPMIVersion           uint8
	ManufacturerID        [3]byte
	ProductID             uint16
	DeviceGUID            [16]byte
}

func (self *McDeviceConfirmationRecord) GetHeader() SensorRecordHeader {
	return self.SensorRecordHeader
}

func (self *McDeviceConfirmationRecord) ReadBytes(r *protocol.Reader) {
	// Header
	self.RecordId = r.ReadUint16()    // 1-2
	self.SdrVersion = r.ReadUint8()   // 3
	self.RecordType = r.ReadUint8()   // 4
	self.RecordLength = r.ReadUint8() // 5

	// Key Fields
	self.DeviceSlaveAddress = r.ReadUint8() >> 1 // 6
	self.DeviceID = r.ReadUint8()                // 7
	numberAndRevision := r.ReadUint8()           // 8
	self.ChannelNumber = numberAndRevision >> 4
	self.ChannelRevision = numberAndRevision & 0x07

	// Data
	self.FirmwareMajorRevision = r.ReadUint8() & 0x7F // 9
	self.FirmwareMinorRevision = r.ReadUint8()        // 10
	self.IPMIVersion = r.ReadUint8()                  // 11
	self.ManufacturerID[0] = r.ReadByte()             // 12
	self.ManufacturerID[1] = r.ReadByte()             // 13
	self.ManufacturerID[2] = r.ReadByte()             // 14
	self.ProductID = r.ReadUint16()                   // 15:16
	copy(self.DeviceGUID[:], r.ReadBytes(16))
}

type MessageChannelInfoType uint8
type BMCMessageChannelInfoRecord struct {
	// Header
	SensorRecordHeader

	// Data
	MessageChannel0Info             MessageChannelInfoType
	MessageChannel1Info             MessageChannelInfoType
	MessageChannel2Info             MessageChannelInfoType
	MessageChannel3Info             MessageChannelInfoType
	MessageChannel4Info             MessageChannelInfoType
	MessageChannel5Info             MessageChannelInfoType
	MessageChannel6Info             MessageChannelInfoType
	MessageChannel7Info             MessageChannelInfoType
	MessagingInterruptType          uint8
	EventMessageBufferInterruptType uint8
	Reserved                        uint8
}

func (self *BMCMessageChannelInfoRecord) GetHeader() SensorRecordHeader {
	return self.SensorRecordHeader
}

func (self *BMCMessageChannelInfoRecord) ReadBytes(r *protocol.Reader) {
	// Header
	self.RecordId = r.ReadUint16()    // 1-2
	self.SdrVersion = r.ReadUint8()   // 3
	self.RecordType = r.ReadUint8()   // 4
	self.RecordLength = r.ReadUint8() // 5

	self.MessageChannel0Info = MessageChannelInfoType(r.ReadUint8()) // 6
	self.MessageChannel1Info = MessageChannelInfoType(r.ReadUint8()) // 7
	self.MessageChannel2Info = MessageChannelInfoType(r.ReadUint8()) // 8
	self.MessageChannel3Info = MessageChannelInfoType(r.ReadUint8()) // 9
	self.MessageChannel4Info = MessageChannelInfoType(r.ReadUint8()) // 10
	self.MessageChannel5Info = MessageChannelInfoType(r.ReadUint8()) // 11
	self.MessageChannel6Info = MessageChannelInfoType(r.ReadUint8()) // 12
	self.MessageChannel7Info = MessageChannelInfoType(r.ReadUint8()) // 13
	self.MessagingInterruptType = r.ReadUint8()                      // 14
	self.EventMessageBufferInterruptType = r.ReadUint8()             // 15
	self.Reserved = r.ReadUint8()                                    // 16
}

type OEMRecord struct {
	// Header
	SensorRecordHeader
	ManufacturerID [3]byte
	OEMData        []byte
}

func (self *OEMRecord) GetHeader() SensorRecordHeader {
	return self.SensorRecordHeader
}

func (self *OEMRecord) ReadBytes(r *protocol.Reader) {
	// Header
	self.RecordId = r.ReadUint16()    // 1-2
	self.SdrVersion = r.ReadUint8()   // 3
	self.RecordType = r.ReadUint8()   // 4
	self.RecordLength = r.ReadUint8() // 5

	self.ManufacturerID[0] = r.ReadByte() // 6
	self.ManufacturerID[1] = r.ReadByte() // 7
	self.ManufacturerID[2] = r.ReadByte() // 8
	self.OEMData = r.ReadCopy(r.Len())    // 16
}

func calcFormula(value, length int32, sensorUnits1 uint8, M, B, Rexp int16, linearization uint8) (float64, error) {
	dataFormat := int32((uint8(sensorUnits1) & 0xc0) >> 6)
	base := int32(0)

	if length <= 1 {
		return 0, errors.New("length is < 1.")
	}

	switch dataFormat {
	case 0: // unsigned
		base = value
		break
	case 1: // 1's complement
		base = int32(decode1sComplement(int16(value), uint(length)-1))
		break
	case 2: // 2's complement
		base = int32(decode2sComplement(int16(value), uint(length)-1))
		break
	case 3: // no analog reading
		base = value
		break
	default:
		return 0, errors.New("Invalid data format in sensorUnits1")
	}

	result := (float64(M)*float64(base) + float64(B)) * math.Pow(10, float64(Rexp))

	switch linearization {
	case 0:
		return result, nil
	case 1:
		return math.Log(result), nil
	case 2:
		return math.Log10(result), nil
	case 3:
		return math.Log(result) / math.Log(2), nil
	case 5:
		return math.Pow(10, result), nil
	case 6:
		return math.Pow(2, result), nil
	case 7:
		return 1 / result, nil
	case 8:
		return math.Pow(result, 2), nil
	case 9:
		return math.Pow(result, 3), nil
	case 10:
		return math.Pow(result, 0.5), nil
	case 11:
		return math.Pow(result, 0.33), nil
	default:
		return 0, errors.New("Unsupported linearization type")
	}
}

// func littleEndianBcdByteToInt32(value byte) int32 {
//  lower := (byteToInt32(value) & 0xf0) >> 4
//  higher := byteToInt32(value) & 0x0f

//  return higher*10 + lower
// }

func byteToInt32(value int8) int32 {
	// if value is not less than 0 everything is fine
	// if value is lesser than 0, it means that we encoded there a number
	// greater than 127
	if value >= 0 {
		return int32(value)
	} else {
		return int32(value) + 256
	}
}

func decode2sComplement(value int16, msb uint) int16 {
	result := value
	base := false
	if (value & (0x1 << msb)) != 0 {
		base = true
	}

	for i := uint(15); i > msb; i-- {
		mask := int16(0x1) << i
		if !base {
			result &= ^mask
		} else {
			result |= mask
		}
	}
	return result
}

func decode1sComplement(value int16, msb uint) int16 {
	result := value
	base := false
	if (value & (0x1 << msb)) != 0 {
		base = true
	}
	if base {
		for i := uint(15); i > msb; i-- {
			mask := int16(0x1) << i
			result |= mask
		}
		result = -(^result)
	}
	return result
}

func decodeName(codingTypeAndLen uint8, name []byte) string {
	length := codingTypeAndLen & 0x3F
	switch (codingTypeAndLen & 0xc0) >> 6 {
	case 0: // unicode
		return string(name[:length])
	case 1: // BCD plus
		return decodeBcdPlus(name[:length])
	case 2: // 6-bit packed ASCII
		return decode6bitAscii(name[:length])
	case 3: // 8-bit ASCII + Latin 1
		bs, err := charmap.ISO8859_1.NewDecoder().Bytes(name[:length])
		if nil != err {
			return string(name[:length])
		}
		return string(bs[:length])
	default:
		panic("Invalid coding type.")
	}
}

func decodeBcdPlus(text []byte) string {
	result := make([]rune, len(text)*2)
	for i, c := range text {
		result[2*i] = decodeBcdPlusChar((c & 0xf0) >> 4)
		result[2*i+1] = decodeBcdPlusChar(c & 0xf)
	}

	return string(result)
}

func decodeBcdPlusChar(ch byte) rune {
	switch ch {
	case 0x0:
		return '0'
	case 0x1:
		return '1'
	case 0x2:
		return '2'
	case 0x3:
		return '3'
	case 0x4:
		return '4'
	case 0x5:
		return '5'
	case 0x6:
		return '6'
	case 0x7:
		return '7'
	case 0x8:
		return '8'
	case 0x9:
		return '9'
	case 0xa:
		return ' '
	case 0xb:
		return '-'
	case 0xc:
		return '.'
	case 0xd:
		return ':'
	case 0xe:
		return ','
	case 0xf:
		return '_'
	default:
		panic("Invalid ch value")
	}
}

func decode6bitAscii(text []byte) string {
	cnt := len(text)
	if cnt%3 != 0 {
		cnt += 3 - cnt%3
	}

	newText := make([]byte, cnt/3*4)

	index := 0
	for i, c := range text {
		switch i % 3 {
		case 0:
			newText[index] = c & 0x3f
			index++
			newText[index] = (c & 0xc0) >> 6
		case 1:
			newText[index] |= (c & 0xf) << 2
			index++
			newText[index] = (c & 0xf0) >> 4
		case 2:
			newText[index] |= (c & 0x3) << 4
			index++
			newText[index] = (c & 0xfc) >> 2
			index++
		}
	}

	for i, c := range newText {
		newText[i] = c + 0x20
	}
	return string(newText)
}
