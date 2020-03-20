package protocol

import (
	"bytes"
	"reflect"

	"testing"
)

var ping_request = []byte{0x06, 0x00, 0xff, 0x06,
	0x00, 0x00, 0x11, 0xbe,
	0x80, 0x00, 0x00, 0x00}

var pong_response = []byte{0x06, 0x00, 0xff, 0x06,
	0x00, 0x00, 0x11, 0xbe,
	0x40, 0x11, 0x00, 0x10,
	0x00, 0x00, 0x11, 0xbe,
	0x00, 0x00, 0x00, 0x00,
	0x81, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00}

func TestPingRequest(t *testing.T) {
	var ping asfMessage
	if e := FromBytes(&ping, ping_request); nil != e {
		t.Error(e)
		return
	}

	assertEquals := func(t *testing.T, field string, actual, excepted interface{}) {
		if !reflect.DeepEqual(actual, excepted) {
			t.Error("["+field+"] excepted is", excepted)
			t.Error("["+field+"] actual   is", actual)
		}
	}

	assertEquals(t, "rmcp.Version", ping.RMCP.Version, uint8(6))
	assertEquals(t, "rmcp.Reserved", ping.RMCP.Reserved, uint8(0))
	assertEquals(t, "rmcp.RMCPSequenceNumber", ping.RMCP.RMCPSequenceNumber, uint8(0xff))
	assertEquals(t, "rmcp.Class", ping.RMCP.Class, uint8(6))

	assertEquals(t, "asf.IANAEnterpriseNumber", ping.ASF.IANAEnterpriseNumber, uint32(4542))
	assertEquals(t, "asf.MessageType", ping.ASF.MessageType, uint8(asfMessageTypePing))
	assertEquals(t, "asf.MessageTag", ping.ASF.MessageTag, uint8(0))
	assertEquals(t, "asf.Reserved", ping.ASF.Reserved, uint8(0))
	assertEquals(t, "asf.DataLength", ping.ASF.DataLength, uint8(0))

	if bs, e := ToBytes(&ping); nil != e {
		t.Error(e)
		return
	} else if !bytes.Equal(bs, ping_request) {
		t.Error("to bytes is failed.")
		return
	}
}

func TestPongResponse(t *testing.T) {
	var pong asfPong
	var msg = asfMessage{Data: &pong}

	if e := FromBytes(&msg, pong_response); nil != e {
		t.Error(e)
		return
	}

	assertEquals := func(t *testing.T, field string, actual, excepted interface{}) {
		if !reflect.DeepEqual(actual, excepted) {
			t.Error("["+field+"] excepted is", excepted)
			t.Error("["+field+"] actual   is", actual)
		}
	}

	assertEquals(t, "rmcp.Version", msg.RMCP.Version, uint8(6))
	assertEquals(t, "rmcp.Reserved", msg.RMCP.Reserved, uint8(0))
	assertEquals(t, "rmcp.RMCPSequenceNumber", msg.RMCP.RMCPSequenceNumber, uint8(0xff))
	assertEquals(t, "rmcp.Class", msg.RMCP.Class, uint8(6))

	assertEquals(t, "asf.IANAEnterpriseNumber", msg.ASF.IANAEnterpriseNumber, uint32(4542))
	assertEquals(t, "asf.MessageType", msg.ASF.MessageType, uint8(asfMessageTypePong))
	assertEquals(t, "asf.MessageTag", msg.ASF.MessageTag, uint8(0x11))
	assertEquals(t, "asf.Reserved", msg.ASF.Reserved, uint8(0))
	assertEquals(t, "asf.DataLength", msg.ASF.DataLength, uint8(16))

	assertEquals(t, "pong.IANAEnterpriseNumber", pong.IANAEnterpriseNumber, uint32(4542))
	assertEquals(t, "pong.OEM", pong.OEM, uint32(0))
	assertEquals(t, "pong.SupportedEntities", pong.SupportedEntities, uint8(129))
	assertEquals(t, "pong.SupportedInteractions", pong.SupportedInteractions, uint8(0))
	assertEquals(t, "pong.Reserved", pong.Reserved, [6]uint8{0, 0, 0, 0, 0, 0})

	if bs, e := ToBytes(&msg); nil != e {
		t.Error(e)
		return
	} else if !bytes.Equal(bs, pong_response) {
		t.Error("to bytes is failed.")
		return
	}
}
