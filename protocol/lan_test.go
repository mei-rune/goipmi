package protocol

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/runner-mei/goipmi/protocol/commands"
)

var ipmi_v1_auth_capabilities_request = []byte{0x06, 0x00, 0xff, 0x07, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x09, 0x20, 0x18, 0xc8, 0x81, 0x04, 0x38, 0x0e, 0x04, 0x31}
var ipmi_v1_auth_capabilities_response = []byte{06, 0x00, 0xff, 0x07, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x81, 0x1c, 0x63, 0x20, 0x04, 0x38, 0x00, 0x02, 0x80, 0x14, 0x02, 0x00, 0x00, 0x00, 0x00, 0x0c}

func assertEquals(t *testing.T, field string, actual, excepted interface{}) {
	if bs, ok := actual.([]byte); ok {
		if !bytes.Equal(bs, excepted.([]byte)) {
			t.Error("["+field+"] excepted is", len(excepted.([]byte)), excepted.([]byte))
			t.Error("["+field+"] actual   is", len(actual.([]byte)), actual.([]byte))
		}
	} else if !reflect.DeepEqual(actual, excepted) {
		t.Error("["+field+"] excepted is", fmt.Sprintf("%T", excepted), excepted)
		t.Error("["+field+"] actual   is", fmt.Sprintf("%T", actual), actual)
	}
}

func TestIPMIV1AuthCapabilitiesRequest(t *testing.T) {
	var rmcpHeader RMCPHeader
	var ipmiHeader IPMIV1Header
	var reqData AuthCapabilitiesRequest
	var req = Request{Data: &reqData}

	var r Reader
	r.Init(ipmi_v1_auth_capabilities_request)
	rmcpHeader.ReadBytes(&r)
	ipmiHeader.ReadBytes(&r)
	req.ReadBytes(&r)

	if e := r.Err(); nil != e {
		t.Error(e)
		return
	}

	assertEquals(t, "rmcp.Version", rmcpHeader.Version, uint8(rmcpVersion1))
	assertEquals(t, "rmcp.Reserved", rmcpHeader.Reserved, uint8(0))
	assertEquals(t, "rmcp.RMCPSequenceNumber", rmcpHeader.RMCPSequenceNumber, uint8(0xff))
	assertEquals(t, "rmcp.Class", rmcpHeader.Class, uint8(rmcpClassIPMI))

	assertEquals(t, "goipmi.AuthType", ipmiHeader.AuthType, uint8(AuthTypeNone))
	assertEquals(t, "goipmi.Sequence", ipmiHeader.Sequence, uint32(0))
	assertEquals(t, "goipmi.SessionID", ipmiHeader.SessionID, uint32(0))
	assertEquals(t, "goipmi.AuthCode", ipmiHeader.AuthCode, []byte(nil)) // make([]byte, 16))
	assertEquals(t, "goipmi.Length", ipmiHeader.Length, uint8(9))

	assertEquals(t, "goipmi.RsAddr", req.Body.RsAddr, uint8(0x20))
	assertEquals(t, "goipmi.NetFnRsLUN", req.Body.NetFnRsLUN, uint8(0x18))
	assertEquals(t, "goipmi.Checksum", req.Body.Checksum, uint8(0xc8))
	assertEquals(t, "goipmi.RqAddr", req.Body.RqAddr, uint8(0x81))
	assertEquals(t, "goipmi.RqSeq", req.Body.RqSeq, uint8(0x04))
	assertEquals(t, "goipmi.Cmd", req.Body.Cmd, commands.GetChannelAuthenticationCapabilities.Code)

	assertEquals(t, "req.ChannelNumber", reqData.ChannelNumber, uint8(0x0e))
	assertEquals(t, "req.PrivLevel", reqData.PrivLevel, uint8(0x04))

	// assertEquals(t, "req.ChannelNumber", reqData.ChannelNumber, uint8(0x0e))
	// assertEquals(t, "req.AuthTypeSupport", reqData.AuthTypeSupport, uint8(0x20))
	// assertEquals(t, "req.Status", reqData.Status, uint8(0x20))
	// assertEquals(t, "req.Reserved", reqData.Reserved, uint8(0x20))
	// assertEquals(t, "req.OEMID", reqData.OEMID, uint16(0x20))
	// assertEquals(t, "req.OEMAux", reqData.OEMAux, uint8(0x20))

	// if bs, e := ToBytes(&ping); nil != e {
	// 	t.Error(e)
	// 	return
	// } else if !bytes.Equal(bs, ping_request) {
	// 	t.Error("to bytes is failed.")
	// 	return
	// }

	ipmiHeader.Length = 0
	req.Body.Checksum = 0
	w := Writer{}
	w.Init(nil)
	rmcpHeader.WriteBytes(&w)
	ipmiHeader.WriteBytes(&w)

	if nil != w.Err() {
		t.Error(w.Err())
		return
	}
	old_length := w.Len()
	req.WriteBytes(&w)

	if nil != w.Err() {
		t.Error(w.Err())
		return
	}

	w.Bytes()[old_length-1] = uint8(w.Len() - old_length)

	if !bytes.Equal(w.Bytes(), ipmi_v1_auth_capabilities_request) {
		t.Error("excepted is", ipmi_v1_auth_capabilities_request)
		t.Error("actual   is", w.Bytes())
		return
	}
}

func TestIPMIV1AuthCapabilitiesResponse(t *testing.T) {
	var rmcpHeader RMCPHeader
	var ipmiHeader IPMIV1Header
	var respData AuthCapabilitiesResponse
	var resp = Response{Data: &respData}

	var r Reader
	r.Init(ipmi_v1_auth_capabilities_response)
	rmcpHeader.ReadBytes(&r)
	ipmiHeader.ReadBytes(&r)
	resp.ReadBytes(&r)

	if e := r.Err(); nil != e {
		t.Error(e)
		return
	}

	assertEquals(t, "rmcp.Version", rmcpHeader.Version, uint8(rmcpVersion1))
	assertEquals(t, "rmcp.Reserved", rmcpHeader.Reserved, uint8(0))
	assertEquals(t, "rmcp.RMCPSequenceNumber", rmcpHeader.RMCPSequenceNumber, uint8(0xff))
	assertEquals(t, "rmcp.Class", rmcpHeader.Class, uint8(rmcpClassIPMI))

	assertEquals(t, "goipmi.AuthType", ipmiHeader.AuthType, uint8(AuthTypeNone))
	assertEquals(t, "goipmi.Sequence", ipmiHeader.Sequence, uint32(0))
	assertEquals(t, "goipmi.SessionID", ipmiHeader.SessionID, uint32(0))
	assertEquals(t, "goipmi.AuthCode", ipmiHeader.AuthCode, []byte(nil)) // make([]byte, 16))
	assertEquals(t, "goipmi.Length", ipmiHeader.Length, uint8(16))

	assertEquals(t, "goipmi.RsAddr", resp.Body.RsAddr, uint8(0x81))
	assertEquals(t, "goipmi.NetFnRsLUN", resp.Body.NetFnRsLUN, uint8(0x1c))
	assertEquals(t, "goipmi.Checksum", resp.Body.Checksum, uint8(0x63))
	assertEquals(t, "goipmi.RqAddr", resp.Body.RqAddr, uint8(0x20))
	assertEquals(t, "goipmi.RqSeq", resp.Body.RqSeq, uint8(0x04))
	assertEquals(t, "goipmi.Cmd", resp.Body.Cmd, commands.GetChannelAuthenticationCapabilities.Code)
	assertEquals(t, "goipmi.CompletionCode", resp.CompletionCode, CompletionCode(CommandCompleted))

	assertEquals(t, "resp.ChannelNumber", respData.ChannelNumber, uint8(0x02))
	assertEquals(t, "resp.AuthTypeSupport", respData.AuthTypeSupport, uint8(0x80))
	assertEquals(t, "resp.Status", respData.Status, uint8(0x14))
	assertEquals(t, "resp.Reserved", respData.Reserved, uint8(0x02))
	assertEquals(t, "resp.OEMID", respData.OEMID, OemID(0))
	assertEquals(t, "resp.OEMAux", respData.OEMAux, uint8(0x00))

	// if bs, e := ToBytes(&ping); nil != e {
	//  t.Error(e)
	//  return
	// } else if !bytes.Equal(bs, ping_request) {
	//  t.Error("to bytes is failed.")
	//  return
	// }

	ipmiHeader.Length = 0
	resp.Body.Checksum = 0
	w := Writer{}
	w.Init(nil)
	rmcpHeader.WriteBytes(&w)
	ipmiHeader.WriteBytes(&w)

	if nil != w.Err() {
		t.Error(w.Err())
		return
	}
	old_length := w.Len()
	resp.WriteBytes(&w)

	if nil != w.Err() {
		t.Error(w.Err())
		return
	}

	w.Bytes()[old_length-1] = uint8(w.Len() - old_length)

	if !bytes.Equal(w.Bytes(), ipmi_v1_auth_capabilities_response) {
		t.Error("excepted is", ipmi_v1_auth_capabilities_response)
		t.Error("actual   is", w.Bytes())
		return
	}
}
