package protocol

import (
	"fmt"

	"github.com/runner-mei/goipmi/protocol/commands"
)

type transport interface {
	open() error
	close() error
	isConnected() bool
	send(interface{}, interface{}) error
}

// Client provides common high level functionality around the underlying transport
type Client struct {
	*ConnectionOption
	transport
}

// NewClient creates a new Client with the given Connection properties
func NewClient(c *ConnectionOption) (*Client, error) {
	var t transport
	switch c.Interface {
	case "lan":
		if c.IntegrityAlgorithm != RAKPAlgorithmIntegrity_None {
			return nil, fmt.Errorf("unsupported integrity algorithm for lan: %s", c.IntegrityAlgorithm.String())
		}

		if c.ConfidentialityAlgorithm != RAKPAlgorithmEncryto_None {
			return nil, fmt.Errorf("unsupported confidentiality algorithm for lan: %s", c.ConfidentialityAlgorithm.String())
		}

		t = newLan(c)
	case "lanplus":
		if c.ConfidentialityAlgorithm != RAKPAlgorithmEncryto_None &&
			c.ConfidentialityAlgorithm != RAKPAlgorithmEncryto_AES_CBC_128 {
			return nil, fmt.Errorf("unsupported confidentiality algorithm: %s", c.ConfidentialityAlgorithm.String())
		}
		t = newLanPlus(c)
	default:
		return nil, fmt.Errorf("unsupported interface: %s", c.Interface)
	}

	return &Client{
		ConnectionOption: c,
		transport:        t,
	}, nil
}

// Open a new IPMI session
func (c *Client) IsConnected() bool {
	// TODO: auto-select transport based on BMC capabilities
	return c.isConnected()
}

// Open a new IPMI session
func (c *Client) Open() error {
	// TODO: auto-select transport based on BMC capabilities
	return c.open()
}

// Close the IPMI session
func (c *Client) Close() error {
	return c.close()
}

// Send a Request and unmarshal to given Response type
func (c *Client) Send(req *Request, resp *Response) error {
	// TODO: handle retry, timeouts, etc.
	return c.send(req, resp)
}

func (c *Client) Exec(cmd commands.CommandCode, req, resp interface{}) error {
	request := NewRequest(cmd, req)
	response := NewResponse(cmd, resp)
	if err := c.Send(request, response); err != nil {
		return err
	}
	if response.Code() != CommandCompleted {
		return response.Code()
	}
	return nil
}

// // DeviceID get the Device ID of the BMC
// func (c *Client) DeviceID() (*goipmi.DeviceIDResponse, error) {
// 	req := NewRequest(goipmi.GetDeviceID, &goipmi.DeviceIDRequest{})
// 	respData := &goipmi.DeviceIDResponse{}
// 	resp := NewResponse(goipmi.GetDeviceID, respData)
// 	return respData, c.Send(req, resp)
// }
