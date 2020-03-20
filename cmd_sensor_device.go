package goipmi

import (
	"bytes"
	"errors"
	"strings"

	"github.com/runner-mei/goipmi/protocol"
)

// section 35.2
type GetDeviceSDRInfoRequest struct {
	HasOperation bool
	Operation    uint8
}

func (self *GetDeviceSDRInfoRequest) SetSDRCountOperation() {
	self.HasOperation = true
	self.Operation = 1
}

func (self *GetDeviceSDRInfoRequest) SetSensorCountOperation() {
	self.HasOperation = true
	self.Operation = 0
}

func (self *GetDeviceSDRInfoRequest) WriteBytes(w *protocol.Writer) {
	if self.HasOperation {
		w.WriteUint8(self.Operation)
	}
}

type GetDeviceSDRInfoResponse struct {
	// CompletionCode
	Count uint8
	Flags uint8
	SPCI  uint32 // LS Byte first
}

// section 35.3
type GetDeviceSDRRequest struct {
	ReservationId uint16 // LS Byte first

	RecordId      uint16 // LS Byte first
	Offset        uint8
	WillReadBytes uint8 // FFh means read entire record.
}

type GetDeviceSDRResponse struct {
	// CompletionCode
	NextRecordId uint16 // LS Byte first
	Data         []byte
}

func (self *GetDeviceSDRResponse) ReadBytes(r *protocol.Reader) {
	self.NextRecordId = r.ReadUint16()
	self.Data = r.ReadBytes(r.Len())
}

// section 35.4
type ReserveDeviceSDRRequest struct {
}

type ReserveDeviceSDRResponse struct {
	// CompletionCode
	Id uint16 // LS Byte first
}

// section 35.5
type GetSensorReadingFactorsRequest struct {
	Number      uint8
	ReadingByte byte
}

type GetSensorReadingFactorsResponse struct {
	// CompletionCode
	NextReading uint8
	M           int16
	Tolerance   uint8
	B           int16
	Accuracy    int16
	AccuracyExp uint8
	Rexp        int8
	Bexp        int8
}

// section 35.6
type GetSensorHysteresisRequest struct {
	Number uint8
	mask   byte
}

type GetSensorHysteresisResponse struct {
	// CompletionCode
	Positive_ThresholdHysteresisValue uint8
	Negative_ThresholdHysteresisValue uint8
}

// section 35.9
type GetSensorThresholdsRequest struct {
	Number uint8
}

type GetSensorThresholdsResponse struct {
	// CompletionCode
	Flags                        uint8
	UpperNonrecoverableThreshold uint8
	UpperCriticalThreshold       uint8
	UpperNonCriticalThreshold    uint8
	LowerNonrecoverableThreshold uint8
	LowerCriticalThreshold       uint8
	LowerNonCriticalThreshold    uint8
}

func (self *GetSensorThresholdsResponse) HasUpperNonrecoverableThreshold() bool {
	return self.Flags&(1<<5) != 0
}
func (self *GetSensorThresholdsResponse) HasUpperCriticalThreshold() bool {
	return self.Flags&(1<<4) != 0
}
func (self *GetSensorThresholdsResponse) HasUpperNonCriticalThreshold() bool {
	return self.Flags&(1<<3) != 0
}
func (self *GetSensorThresholdsResponse) HasLowerNonrecoverableThreshold() bool {
	return self.Flags&(1<<2) != 0
}
func (self *GetSensorThresholdsResponse) HasLowerCriticalThreshold() bool {
	return self.Flags&(1<<1) != 0
}
func (self *GetSensorThresholdsResponse) HasLowerNonCriticalThreshold() bool {
	return self.Flags&1 != 0
}

// section 35.11
type GetSensorEventEnableRequest struct {
	Number uint8
}

type GetSensorEventEnableResponse struct {
	// CompletionCode
	Flags  uint8
	Values [5]byte
}

// section 35.14
type GetSensorReadingRequest struct {
	Number uint8
}

type GetSensorReadingResponse struct {
	// CompletionCode
	Reading uint8
	// DisabledAllEvent   bool
	// DisabledScaning    bool
	// ReadingUnavailable bool
	Flags []byte
}

func (self *GetSensorReadingResponse) IsOk() bool {
	for _, c := range self.Flags {
		if c != 0 {
			return false
		}
	}
	return true
}

func (self *GetSensorReadingResponse) GetAllEventDisabled() bool {
	return 0 == self.Flags[0]&0x80
}

func (self *GetSensorReadingResponse) GetScaningDisabled() bool {
	return 0 == self.Flags[0]&0x40
}

func (self *GetSensorReadingResponse) GetReadingUnavailable() bool {
	return 0 != self.Flags[0]&0x20
}

func (self *GetSensorReadingResponse) ToEventString(eventOrReadingTypeCode uint8) string {
	var buf bytes.Buffer
	if eventOrReadingTypeCode == 1 {
		for _, v := range []struct {
			t    ThresholdEventType
			name string
		}{{LOWER_NON_CRITICAL_GOING_LOW, "lower_non_critical_going_low"},
			{LOWER_NON_CRITICAL_GOING_HIGH, "lower_non_critical_going_high"},
			{LOWER_CRITICAL_GOING_LOW, "lower_critical_going_low"},
			{LOWER_CRITICAL_GOING_HIGH, "lower_critical_going_high"},
			{LOWER_NON_RECOVERABLE_GOING_LOW, "lower_non_recoverable_going_low"},
			{LOWER_NON_RECOVERABLE_GOING_HIGH, "lower_non_recoverable_going_high"},
			{UPPER_NON_CRITICAL_GOING_LOW, "upper_non_critical_going_low"},
			{UPPER_NON_CRITICAL_GOING_HIGH, "upper_non_critical_going_high"},
			{UPPER_CRITICAL_GOING_LOW, "upper_critical_going_low"},
			{UPPER_CRITICAL_GOING_HIGH, "upper_critical_going_high"},
			{UPPER_NON_RECOVERABLE_GOING_LOW, "upper_non_recoverable_going_low"},
			{UPPER_NON_RECOVERABLE_GOING_HIGH, "upper_non_recoverable_going_high"},
		} {
			if ok, err := self.GetAssertionThresholdEventOccurred(v.t); err != nil {
				//panic(err)
			} else if ok {
				buf.WriteString(v.name + "_assertion,")
			}

			if ok, err := self.GetDeassertionThresholdEventOccurred(v.t); err != nil {
				//panic(err)
			} else if ok {
				buf.WriteString(v.name + "_dessertion,")
			}
		}
		return buf.String()
	}

	var specs []ThresholdEventSpec
	if eventOrReadingTypeCode == 2 {
		specs = []ThresholdEventSpec{
			{typ: 0, name: "transition to idle"},
			{typ: 1, name: "transition to active"},
			{typ: 2, name: "transition to busy"},
		}
	} else if eventOrReadingTypeCode == 3 {
		specs = []ThresholdEventSpec{
			{typ: 0, name: "state deasserted"},
			{typ: 1, name: "state asserted"},
		}
	} else if eventOrReadingTypeCode == 4 {
		specs = []ThresholdEventSpec{
			{typ: 0, name: "predictive failure deasserted"},
			{typ: 1, name: "predictive failure asserted"},
		}
	} else if eventOrReadingTypeCode == 5 {
		specs = []ThresholdEventSpec{
			{typ: 0, name: "limit not exceeded"},
			{typ: 1, name: "limit exceeded"},
		}
	} else if eventOrReadingTypeCode == 6 {
		specs = []ThresholdEventSpec{
			{typ: 0, name: "performance met"},
			{typ: 1, name: "performance lags"},
		}
	} else if eventOrReadingTypeCode == 7 {
		specs = []ThresholdEventSpec{
			{typ: 0, name: "transition to OK"},
			{typ: 1, name: "transition to Non-critical from OK"},
			{typ: 2, name: "transition to Critical from less servere"},
			{typ: 3, name: "transition to Non-recoverable from less servere"},
			{typ: 4, name: "transition to Non-critical from more servere"},
			{typ: 5, name: "transition to Critical from Non-recoverable"},
			{typ: 6, name: "transition to Non-recoverable"},
			{typ: 7, name: "monitor"},
			{typ: 8, name: "informational"},
		}
	} else if eventOrReadingTypeCode == 8 {
		specs = []ThresholdEventSpec{
			{typ: 0, name: "device removed"},
			{typ: 0, name: "device inserted"},
		}
	} else if eventOrReadingTypeCode == 9 {
		specs = []ThresholdEventSpec{
			{typ: 0, name: "device disabled"},
			{typ: 0, name: "device enabled"},
		}
	} else if eventOrReadingTypeCode == 10 {
		specs = []ThresholdEventSpec{
			{typ: 0, name: "transition to Running"},
			{typ: 1, name: "transition to In Test"},
			{typ: 2, name: "transition to Power Off"},
			{typ: 3, name: "transition to On line"},
			{typ: 4, name: "transition to Off line"},
			{typ: 5, name: "transition to Off Duty"},
			{typ: 6, name: "transition to Degraded"},
			{typ: 7, name: "transition to Power Save"},
			{typ: 8, name: "install error"},
		}
	} else if eventOrReadingTypeCode == 11 {
		specs = []ThresholdEventSpec{
			{typ: 0, name: "full redundancy"},
			{typ: 1, name: "redundancy lost"},
			{typ: 2, name: "redundancy degraded"},
			{typ: 3, name: "non-redundant(sufficient resources from redundant)"},
			{typ: 4, name: "non-redundant(sufficient resources from insufficient resources)"},
			{typ: 5, name: "non-redundant(insufficient resources)"},
			{typ: 6, name: "redundancy Degraded from fully redundant"},
			{typ: 7, name: "redundancy Degraded from non-redundant"},
		}
	} else if eventOrReadingTypeCode == 12 {
		specs = []ThresholdEventSpec{
			{typ: 0, name: "D0 power state"},
			{typ: 1, name: "D1 power state"},
			{typ: 2, name: "D2 power state"},
			{typ: 3, name: "D3 power state"},
		}
	}

	for _, v := range specs {
		if ok, err := self.GetAssertionDiscreteEventOccurred(v.typ); err != nil {
			panic(err)
		} else if ok {
			if strings.HasSuffix(v.name, ")") {
				buf.WriteString(strings.TrimSuffix(v.name, ")") + " -- assertion),")
			} else {
				buf.WriteString(v.name + "(assertion),")
			}
		}
		if ok, err := self.GetDeassertionDiscreteEventOccurred(v.typ); err != nil {
			panic(err)
		} else if ok {
			if strings.HasSuffix(v.name, ")") {
				buf.WriteString(strings.TrimSuffix(v.name, ")") + " -- deassertion),")
			} else {
				buf.WriteString(v.name + "(deassertion),")
			}
		}
	}
	return buf.String()
}

type ThresholdEventSpec struct {
	typ  DiscreteEventType
	name string
}

type DiscreteEventType uint
type ThresholdEventType uint

const (
	LOWER_NON_CRITICAL_GOING_LOW ThresholdEventType = iota
	LOWER_NON_CRITICAL_GOING_HIGH
	LOWER_CRITICAL_GOING_LOW
	LOWER_CRITICAL_GOING_HIGH
	LOWER_NON_RECOVERABLE_GOING_LOW
	LOWER_NON_RECOVERABLE_GOING_HIGH
	UPPER_NON_CRITICAL_GOING_LOW
	UPPER_NON_CRITICAL_GOING_HIGH
	UPPER_CRITICAL_GOING_LOW
	UPPER_CRITICAL_GOING_HIGH
	UPPER_NON_RECOVERABLE_GOING_LOW
	UPPER_NON_RECOVERABLE_GOING_HIGH

	DISCRETE_STATE_0 DiscreteEventType = iota
	DISCRETE_STATE_1
	DISCRETE_STATE_2
	DISCRETE_STATE_3
	DISCRETE_STATE_4
	DISCRETE_STATE_5
	DISCRETE_STATE_6
	DISCRETE_STATE_7
	DISCRETE_STATE_8
	DISCRETE_STATE_9
	DISCRETE_STATE_10
	DISCRETE_STATE_11
	DISCRETE_STATE_12
	DISCRETE_STATE_13
	DISCRETE_STATE_14
)

var ErrNumOverflow = errors.New("num is overflow")

func (self *GetSensorReadingResponse) GetAssertionThresholdEventOccurred(num ThresholdEventType) (bool, error) {
	if num < 8 && len(self.Flags) >= 2 {
		return 0 != self.Flags[1]&(1<<num), nil
	}
	if num < 15 && len(self.Flags) >= 3 {
		num -= 8
		return 0 != self.Flags[2]&(1<<num), nil
	}
	return false, ErrNumOverflow
}

func (self *GetSensorReadingResponse) GetDeassertionThresholdEventOccurred(num ThresholdEventType) (bool, error) {
	if num < 8 {
		if len(self.Flags) < 4 {
			return false, nil //errors.New("Deassertion flags is missing.")
		}
		return 0 != self.Flags[3]&(1<<num), nil
	}
	if num < 15 {
		if len(self.Flags) < 5 {
			return false, nil //errors.New("Deassertion flags is missing.")
		}

		num -= 8
		return 0 != self.Flags[4]&(1<<num), nil
	}
	return false, ErrNumOverflow
}

func (self *GetSensorReadingResponse) GetAssertionDiscreteEventOccurred(num DiscreteEventType) (bool, error) {
	if num < 8 && len(self.Flags) >= 2 {
		return 0 != self.Flags[1]&(1<<num), nil
	}
	if num < 15 && len(self.Flags) >= 3 {
		num -= 8
		return 0 != self.Flags[2]&(1<<num), nil
	}
	return false, ErrNumOverflow
}

func (self *GetSensorReadingResponse) GetDeassertionDiscreteEventOccurred(num DiscreteEventType) (bool, error) {
	if num < 8 {
		if len(self.Flags) < 4 {
			return false, nil //errors.New("Deassertion flags is missing.")
		}
		return 0 != self.Flags[3]&(1<<num), nil
	}
	if num < 15 {
		if len(self.Flags) < 5 {
			return false, nil //errors.New("Deassertion flags is missing.")
		}

		num -= 8
		return 0 != self.Flags[4]&(1<<num), nil
	}
	return false, ErrNumOverflow
}

func (self *GetSensorReadingResponse) ReadBytes(r *protocol.Reader) {
	self.Reading = r.ReadUint8()
	self.Flags = r.ReadCopy(r.Len())
	if len(self.Flags) < 2 {
		r.SetError(ErrInsufficientBytes)
	}
}

// section 35.14
type GetSensorTypeRequest struct {
	Number uint8
}

type GetSensorTypeResponse struct {
	// CompletionCode
	SensorType              uint8
	EventAndReadingTypeCode uint8
}
