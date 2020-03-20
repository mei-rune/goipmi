package goipmi

import (
	"errors"

	"github.com/runner-mei/goipmi/protocol"
	"github.com/runner-mei/goipmi/protocol/commands"
)

var (
	NullRead = protocol.NullRead

	ErrInsufficientBytes  = protocol.ErrInsufficientBytes
	ErrPasswordNotMatch   = protocol.ErrPasswordNotMatch
	ErrReadingUnavailable = protocol.ErrReadingUnavailable
	ErrIgnoreSensor       = protocol.ErrIgnoreSensor
	ErrRequestData        = protocol.ErrRequestData

	ParsePrivLevel             = commands.ParsePrivLevel
	ParseAuthMethod            = protocol.ParseAuthMethod
	ParseIntegrityMethod       = protocol.ParseIntegrityMethod
	ParseConfidentialityMethod = protocol.ParseConfidentialityMethod

	SetFake = protocol.SetFake
)

type ConnectionOption = protocol.ConnectionOption

type ClientHandler interface {
	Open() error
	Close() error
	IsConnected() bool
	//Send(req *Request, resp *Response) error
	//DeviceID() (*DeviceIDResponse, error)

	Exec(cmd commands.CommandCode, req, resp interface{}) error
}

func NewClient(opt *protocol.ConnectionOption) (*Client, error) {
	handler, err := protocol.NewClient(opt)
	if err != nil {
		return nil, err
	}
	return &Client{ClientHandler: handler}, nil
}

type Client struct {
	ClientHandler
}

func (c *Client) IsConnected() bool {
	return c.ClientHandler.IsConnected()
}

// DeviceID get the Device ID of the BMC
func (c *Client) GetDeviceID() (*DeviceIDResponse, error) {
	req := &DeviceIDRequest{}
	resp := &DeviceIDResponse{}
	return resp, c.Exec(GetDeviceID, req, resp)
}

func (c *Client) GetACPIPowerState() (*GetACPIPowerStateResponse, error) {
	var getACPIPowerStateRequest GetACPIPowerStateRequest
	var getACPIPowerStateResponse GetACPIPowerStateResponse
	return &getACPIPowerStateResponse, c.Exec(GetACPIPowerState,
		&getACPIPowerStateRequest,
		&getACPIPowerStateResponse)
}

func (c *Client) GetChassisCapabilities() (*GetChassisCapabilitiesResponse, error) {
	var getChassisCapabilitiesRequest GetChassisCapabilitiesRequest
	var getChassisCapabilitiesResponse GetChassisCapabilitiesResponse
	return &getChassisCapabilitiesResponse,
		c.Exec(GetChassisCapabilities,
			&getChassisCapabilitiesRequest,
			&getChassisCapabilitiesResponse)
}

func (c *Client) GetChassisStatus() (*GetChassisStatusResponse, error) {
	var getChassisStatusRequest GetChassisStatusRequest
	var getChassisStatusResponse GetChassisStatusResponse
	return &getChassisStatusResponse,
		c.Exec(GetChassisStatus,
			&getChassisStatusRequest,
			&getChassisStatusResponse)
}

func (c *Client) GetSystemRestartCause() (*GetSystemRestartCauseResponse, error) {
	var getSystemRestartCauseRequest GetSystemRestartCauseRequest
	var getSystemRestartCauseResponse GetSystemRestartCauseResponse
	return &getSystemRestartCauseResponse,
		c.Exec(GetSystemRestartCause,
			&getSystemRestartCauseRequest,
			&getSystemRestartCauseResponse)
}

func (c *Client) GetSDRRepositoryInfo() (*GetSDRInfoResponse, error) {
	var getSDRInfoRequest GetSDRInfoRequest
	var getSDRInfoResponse GetSDRInfoResponse
	return &getSDRInfoResponse,
		c.Exec(GetSDRRepositoryInfo,
			&getSDRInfoRequest,
			&getSDRInfoResponse)
}

func (c *Client) GetReserveSDRRepository() (*ReserveSDRResponse, error) {
	var reserveSDRRequest ReserveSDRRequest
	var reserveSDRResponse ReserveSDRResponse
	return &reserveSDRResponse,
		c.Exec(ReserveSDRRepository,
			&reserveSDRRequest,
			&reserveSDRResponse)
}

func (c *Client) GetSensorReading(number uint8) (*GetSensorReadingResponse, error) {
	var getSensorReadingRequest GetSensorReadingRequest
	var getSensorReadingResponse GetSensorReadingResponse

	getSensorReadingRequest.Number = number
	return &getSensorReadingResponse,
		c.Exec(GetSensorReading,
			&getSensorReadingRequest,
			&getSensorReadingResponse)
}

func (c *Client) GetSensorHysteresis(number uint8) (*GetSensorHysteresisResponse, error) {
	var getSensorHysteresisRequest GetSensorHysteresisRequest
	var getSensorHysteresisResponse GetSensorHysteresisResponse

	getSensorHysteresisRequest.Number = number
	return &getSensorHysteresisResponse,
		c.Exec(GetSensorHysteresis,
			&getSensorHysteresisRequest,
			&getSensorHysteresisResponse)
}

func (c *Client) GetSensorThresholds(number uint8) (*GetSensorThresholdsResponse, error) {
	var getSensorThresholdsRequest GetSensorThresholdsRequest
	var getSensorThresholdsResponse GetSensorThresholdsResponse

	getSensorThresholdsRequest.Number = number
	return &getSensorThresholdsResponse,
		c.Exec(GetSensorThresholds,
			&getSensorThresholdsRequest,
			&getSensorThresholdsResponse)
}

func (c *Client) GetPOH() (*GetPOHCounterResponse, error) {
	var getPOHCounterRequest GetPOHCounterRequest
	var getPOHCounterResponse GetPOHCounterResponse
	return &getPOHCounterResponse,
		c.Exec(GetPOHCounter,
			&getPOHCounterRequest,
			&getPOHCounterResponse)
}

const BLOCK_LENGTH = 16

func (c *Client) ListSDR(reservationId uint16) ([]Record, error) {
	var results = make([]Record, 0, 32)
	record_id := uint16(0)
	for {
		offset := uint8(0)
		next_record_id := uint16(0)
		var data = RecordData{Data: make([]byte, 0, 64)}

		for j := 0; ; j++ {
			var blockLength uint8 = BLOCK_LENGTH
			if j > 0 {
				remainLen := data.SdrRecordLength() - (len(data.Data) - 4)
				if remainLen < BLOCK_LENGTH {
					blockLength = uint8(remainLen)
				}
			}

			var getSDRRequest = GetSDRRequest{
				ReservationId: reservationId,
				RecordId:      record_id,
				Offset:        offset,
				WillReadBytes: blockLength}
			var getSDRResponse GetSDRResponse
			getSDRResponse.Data = &data

			if e := c.Exec(GetSDR, &getSDRRequest, &getSDRResponse); e != nil {
				if e == ErrRequestData {
					next_record_id = getSDRResponse.NextRecordId
					break
				}
				return nil, errors.New("get Sdr info, " + e.Error())
			}

			if getSDRResponse.Data.SdrIsOk() || j > 200 {
				// if len(data.Data) == 0 {
				// 	for h := 0; ; h++ {
				// 		getSDRRequest.ReservationId = reservationId
				// 		getSDRRequest.RecordId = record_id
				// 		getSDRRequest.Offset = offset
				// 		getSDRRequest.WillReadBytes = 16
				// 		getSDRResponse.Data = &data
				// 		if e := c.Exec(GetSDR, &getSDRRequest, &getSDRResponse); e != nil {
				// 			return nil, errors.New("get Sdr info, " + e.Error())
				// 		}

				// 		if getSDRResponse.DataLength < 16 || h > 200 {
				// 			next_record_id = getSDRResponse.NextRecordId
				// 			break
				// 		}
				// 		offset += uint8(getSDRResponse.DataLength)
				// 	}
				// }

				next_record_id = getSDRResponse.NextRecordId
				break
			}

			offset += uint8(getSDRResponse.DataLength)
		}
		if len(data.Data) != 0 {
			record, e := data.ToSdrRecord()
			if nil != e {
				return nil, errors.New("toRecord:" + e.Error())
			}
			results = append(results, record)
		}

		record_id = next_record_id
		if record_id == 0xffff {
			break
		}
	}
	return results, nil
}

func (c *Client) ListSEL(reservationId uint16) ([]interface{}, error) {
	var results = make([]interface{}, 0, 32)
	record_id := uint16(0)
	for {
		offset := uint8(0)
		next_record_id := uint16(0)
		var data = RecordData{Data: make([]byte, 0, 64)}

		for j := 0; ; j++ {
			var getSELRequest = GetSELRequest{
				ReservationId: reservationId,
				RecordId:      record_id,
				Offset:        offset,
				WillReadBytes: BLOCK_LENGTH}
			var getSELResponse GetSELResponse
			getSELResponse.Data = &data

			if e := c.Exec(GetSELEntry, &getSELRequest, &getSELResponse); e != nil {
				return nil, errors.New("get SEL info, " + e.Error())
			}

			if getSELResponse.Data.SelIsOk() || j > 200 {
				// if len(data.Data) == 0 {
				// 	for h := 0; ; h++ {
				// 		getSELRequest.ReservationId = reservationId
				// 		getSELRequest.RecordId = record_id
				// 		getSELRequest.Offset = offset
				// 		getSELRequest.WillReadBytes = 16
				// 		getSELResponse.Data = &data
				// 		if e := c.Exec(GetSDR, &getSDRRequest, &getSELResponse); e != nil {
				// 			return nil, errors.New("get Sdr info, " + e.Error())
				// 		}

				// 		if getSELResponse.DataLength < 16 || h > 200 {
				// 			next_record_id = getSELResponse.NextRecordId
				// 			break
				// 		}
				// 		offset += uint8(getSELResponse.DataLength)
				// 	}
				// }

				next_record_id = getSELResponse.NextRecordId
				break
			}

			offset += uint8(getSELResponse.DataLength)
		}
		if len(data.Data) != 0 {
			// record, e := data.Data
			// if nil != e {
			// 	return nil, errors.New("toRecord:" + e.Error())
			// }
			results = append(results, data.Data)
		}

		record_id = next_record_id
		if record_id == 0xffff {
			break
		}
	}
	return results, nil
}

type SensorReadingResponse struct {
	Response *GetSensorReadingResponse
	Error    error
}

func (c *Client) ListFullSDRReading(sdr_list []Record) ([]*FullSensorRecord, []SensorReadingResponse, error) {
	records := make([]*FullSensorRecord, 0, len(sdr_list))
	results := make([]SensorReadingResponse, 0, len(sdr_list))

	for _, v := range sdr_list {
		if full, ok := v.(*FullSensorRecord); ok {
			records = append(records, full)
		}
	}

	for _, rec := range records {
		// if rec.CanIgnore() {
		// 	results = append(results, SensorReadingResponse{Error: ErrIgnoreSensor})
		// } else

		if res, err := c.GetSensorReading(rec.SensorNumber); err != nil {
			results = append(results, SensorReadingResponse{Error: err})
		} else if res.GetReadingUnavailable() {
			results = append(results, SensorReadingResponse{Error: ErrReadingUnavailable})
		} else {
			results = append(results, SensorReadingResponse{Response: res})
		}
	}
	return records, results, nil
}
