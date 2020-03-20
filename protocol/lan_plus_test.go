package protocol

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
)

var ipmi_v2_open_session_request = []byte{0x06, 0x00, 0xff, 0x07, 0x06, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0xa4, 0xa3, 0xa2, 0xa0, 0x00, 0x00, 0x00, 0x08, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x08, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x08, 0x01, 0x00, 0x00, 0x00}

func TestIPMIV2OpenSessionRequest(t *testing.T) {
	var l lanPlus
	var payload OpenSessionRequest

	if err := l.FromBytes(&payload, ipmi_v2_open_session_request); nil != err {
		t.Error(err)
		return
	}

	assertEquals(t, "req.MessageTag", payload.MessageTag, uint8(0))
	assertEquals(t, "req.PrivLevel", payload.PrivLevel, uint8(0))
	assertEquals(t, "req.SessionID", payload.SessionID, uint32(2695013284)) //0xA0A2A3A4

	assertEquals(t, "req.Authentication", payload.Authentication.Algorithm, RAKPAlgorithmAuth_HMAC_SHA1)
	assertEquals(t, "req.Integrity", payload.Integrity.Algorithm, RAKPAlgorithmIntegrity_HMAC_SHA1_96)
	assertEquals(t, "req.Confidentiality", payload.Confidentiality.Algorithm, RAKPAlgorithmEncryto_AES_CBC_128)

	l.sequence--
	bs, err := l.ToBytes(&payload, make([]byte, 0, 1024))
	if nil != err {
		t.Error(err)
		return
	}

	if !bytes.Equal(bs, ipmi_v2_open_session_request) {
		t.Error("excepted is", ipmi_v2_open_session_request)
		t.Error("actual   is", bs)
		return
	}

}

func TestIPMIV2OpenSessionRequest_22(t *testing.T) {
	var l lanPlus
	var payload OpenSessionRequest
	excepted_bs := []byte{0x06, 0x00, 0xff, 0x07, 0x06, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0xa4, 0xa3, 0xa2, 0xa0, 0x00, 0x00, 0x00, 0x08, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x00}

	//excepted_bs := []byte{0x06, 0x00, 0xff, 0x07, 0x06, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0xa4, 0xa3, 0xa2, 0xa0, 0x00, 0x00, 0x00, 0x08, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x00}

	if err := l.FromBytes(&payload, excepted_bs); nil != err {
		t.Error(err)
		return
	}

	assertEquals(t, "req.MessageTag", payload.MessageTag, uint8(0))
	assertEquals(t, "req.PrivLevel", payload.PrivLevel, uint8(0))
	assertEquals(t, "req.SessionID", payload.SessionID, uint32(2695013284)) //0xA0A2A3A4

	assertEquals(t, "req.Authentication", payload.Authentication.Algorithm, RAKPAlgorithmAuth_HMAC_SHA1)
	assertEquals(t, "req.Integrity", payload.Integrity.Algorithm, RAKPAlgorithmIntegrity_None)
	assertEquals(t, "req.Confidentiality", payload.Confidentiality.Algorithm, RAKPAlgorithmEncryto_None)

	l.sequence--
	bs, err := l.ToBytes(&payload, make([]byte, 0, 1024))
	if nil != err {
		t.Error(err)
		return
	}

	if !bytes.Equal(bs, excepted_bs) {
		t.Error("excepted is", excepted_bs)
		t.Error("actual   is", bs)
		return
	}

}

func TestIPMIV2OpenSessionResponse_22(t *testing.T) {
	var l lanPlus
	var payload OpenSessionResponse

	excepted_bs := []byte{0x06, 0x00, 0xff, 0x07, 0x06, 0x11, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x00, 0x00, 0x00, 0x04, 0x00, 0xa4, 0xa3, 0xa2, 0xa0, 0x00, 0x0a, 0x00, 0x03, 0x00, 0x00, 0x00, 0x08, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x00}
	if err := l.FromBytes(&payload, excepted_bs); nil != err {
		t.Error(err)
		return
	}

	assertEquals(t, "req.MessageTag", payload.MessageTag, uint8(0))
	assertEquals(t, "req.StatusCode", payload.StatusCode, uint8(0))
	assertEquals(t, "req.PrivLevel", payload.PrivLevel, uint8(4))
	assertEquals(t, "req.SessionID", payload.SessionID, uint32(2695013284)) //0xA0A2A3A4
	assertEquals(t, "req.MSessionID", payload.MSessionID, uint32(50334208)) //0x00993554

	assertEquals(t, "req.Authentication", payload.Authentication.Algorithm, RAKPAlgorithmAuth_HMAC_SHA1)
	assertEquals(t, "req.Integrity", payload.Integrity.Algorithm, RAKPAlgorithmIntegrity_None)
	assertEquals(t, "req.Confidentiality", payload.Confidentiality.Algorithm, RAKPAlgorithmEncryto_None)

	l.sequence--
	bs, err := l.ToBytes(&payload, make([]byte, 0, 1024))
	if nil != err {
		t.Error(err)
		return
	}

	if !bytes.Equal(excepted_bs, bs) {
		t.Error("excepted is", excepted_bs)
		t.Error("actual   is", bs)
		return
	}
}

func TestIPMIV2OpenSessionResponse_33(t *testing.T) {
	var l lanPlus
	var payload OpenSessionResponse

	excepted_bs := []byte{0x06, 0x00, 0xff, 0x07, 0x06, 0x11, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x00, 0x00, 0x00, 0x04, 0x00, 0x65, 0x00, 0x00, 0x00, 0x00, 0x07, 0x00, 0x03, 0x00, 0x00, 0x00, 0x08, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x08, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x00}
	if err := l.FromBytes(&payload, excepted_bs); nil != err {
		t.Error(err)
		return
	}

	assertEquals(t, "req.MessageTag", payload.MessageTag, uint8(0))
	assertEquals(t, "req.StatusCode", payload.StatusCode, uint8(0))
	assertEquals(t, "req.PrivLevel", payload.PrivLevel, uint8(4))
	assertEquals(t, "req.SessionID", payload.SessionID, uint32(101))        //0xA0A2A3A4
	assertEquals(t, "req.MSessionID", payload.MSessionID, uint32(50333440)) //0x00993554

	assertEquals(t, "req.Authentication", payload.Authentication.Algorithm, RAKPAlgorithmAuth_HMAC_SHA1)
	assertEquals(t, "req.Integrity", payload.Integrity.Algorithm, RAKPAlgorithmIntegrity_HMAC_SHA1_96)
	assertEquals(t, "req.Confidentiality", payload.Confidentiality.Algorithm, RAKPAlgorithmEncryto_None)

	l.sequence--
	bs, err := l.ToBytes(&payload, make([]byte, 0, 1024))
	if nil != err {
		t.Error(err)
		return
	}

	if !bytes.Equal(excepted_bs, bs) {
		t.Error("excepted is", excepted_bs)
		t.Error("actual   is", bs)
		return
	}
}

var ipmi_v2_open_session_response = []byte{0x06, 0x00, 0xff, 0x07, 0x06, 0x11, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x00, 0x00, 0x00, 0x04, 0x00, 0xa4, 0xa3, 0xa2, 0xa0, 0x54, 0x35, 0x99, 0x00, 0x00, 0x00, 0x00, 0x08, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x08, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x08, 0x01, 0x00, 0x00, 0x00}

func TestIPMIV2OpenSessionResponse(t *testing.T) {
	var l lanPlus
	var payload OpenSessionResponse

	if err := l.FromBytes(&payload, ipmi_v2_open_session_response); nil != err {
		t.Error(err)
		return
	}

	assertEquals(t, "req.MessageTag", payload.MessageTag, uint8(0))
	assertEquals(t, "req.StatusCode", payload.StatusCode, uint8(0))
	assertEquals(t, "req.PrivLevel", payload.PrivLevel, uint8(4))
	assertEquals(t, "req.SessionID", payload.SessionID, uint32(2695013284)) //0xA0A2A3A4
	assertEquals(t, "req.MSessionID", payload.MSessionID, uint32(10040660)) //0x00993554

	assertEquals(t, "req.Authentication", payload.Authentication.Algorithm, RAKPAlgorithmAuth_HMAC_SHA1)
	assertEquals(t, "req.Integrity", payload.Integrity.Algorithm, RAKPAlgorithmIntegrity_HMAC_SHA1_96)
	assertEquals(t, "req.Confidentiality", payload.Confidentiality.Algorithm, RAKPAlgorithmEncryto_AES_CBC_128)

	l.sequence--
	bs, err := l.ToBytes(&payload, make([]byte, 0, 1024))
	if nil != err {
		t.Error(err)
		return
	}

	if !bytes.Equal(bs, ipmi_v2_open_session_response) {
		t.Error("excepted is", ipmi_v2_open_session_response)
		t.Error("actual   is", bs)
		return
	}
}

func TestIPMIV2OpenSessionResponse_2(t *testing.T) {
	var excepted = []byte{0x06, 0x00, 0xff, 0x07, 0x06, 0x11, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x04, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	var l lanPlus
	var payload OpenSessionResponse

	if err := l.FromBytes(&payload, excepted); nil != err {
		t.Error(err)
		return
	}

	assertEquals(t, "req.MessageTag", payload.MessageTag, uint8(0))
	assertEquals(t, "req.StatusCode", payload.StatusCode, uint8(4))
	assertEquals(t, "req.PrivLevel", payload.PrivLevel, uint8(0))
	assertEquals(t, "req.SessionID", payload.SessionID, uint32(0))   //0xA0A2A3A4
	assertEquals(t, "req.MSessionID", payload.MSessionID, uint32(0)) //0x00993554

	assertEquals(t, "req.Authentication", payload.Authentication.Algorithm, RAKPAlgorithmAuth_None)
	assertEquals(t, "req.Integrity", payload.Integrity.Algorithm, RAKPAlgorithmIntegrity_None)
	assertEquals(t, "req.Confidentiality", payload.Confidentiality.Algorithm, RAKPAlgorithmEncryto_None)
}

var ipmi_v2_rakp_message_1 = []byte{0x06, 0x00, 0xff, 0x07, 0x06, 0x12, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x29, 0x00, 0x00, 0x00, 0x00, 0x00, 0x54, 0x35, 0x99, 0x00, 0x45, 0x7b, 0x81, 0xbb, 0x81, 0xb8, 0x6c, 0x94, 0x68, 0x97, 0x11, 0x3a, 0xb5, 0xff, 0x3b, 0x30, 0x14, 0x00, 0x00, 0x0d, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72}

func TestIPMIV2RakpMessage1(t *testing.T) {
	var l lanPlus
	var payload RakpMessage1

	if err := l.FromBytes(&payload, ipmi_v2_rakp_message_1); nil != err {
		t.Error(err)
		return
	}

	assertEquals(t, "req.MessageTag", payload.MessageTag, uint8(0))
	assertEquals(t, "req.SessionID", payload.SessionID, uint32(10040660)) // MSessionID 0x00993554

	assertEquals(t, "req.Rand", payload.Rand, [16]byte{69, 123, 129, 187, 129, 184, 108, 148, 104, 151, 17, 58, 181, 255, 59, 48})
	assertEquals(t, "req.PrivLevel", payload.PrivLevel, uint8(20))

	assertEquals(t, "req.UserNameLength", payload.UserNameLength, uint8(13))
	assertEquals(t, "req.UserName", string(payload.UserName[:payload.UserNameLength]), "Administrator")

	l.sequence--
	bs, err := l.ToBytes(&payload, make([]byte, 0, 1024))
	if nil != err {
		t.Error(err)
		return
	}

	if !bytes.Equal(bs, ipmi_v2_rakp_message_1) {
		t.Error("excepted is", ipmi_v2_rakp_message_1)
		t.Error("actual   is", bs)
		return
	}
}

func TestIPMIV2RakpMessage1_22(t *testing.T) {
	var l lanPlus
	var payload RakpMessage1
	excepted_bs := []byte{0x06, 0x00, 0xff, 0x07, 0x06, 0x12, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x22, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0a, 0x00, 0x03, 0xc9, 0x40, 0xa0, 0xed, 0xbe, 0x81, 0x8d, 0xcc, 0xbf, 0x8b, 0xb9, 0x2c, 0x48, 0xff, 0xbd, 0x40, 0x14, 0x00, 0x00, 0x06, 0x55, 0x53, 0x45, 0x52, 0x49, 0x44}

	//excepted_bs := []byte{0x06, 0x00, 0xff, 0x07, 0x06, 0x12, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x22, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0a, 0x00, 0x03, 0xc9, 0x40, 0xa0, 0xed, 0xbe, 0x81, 0x8d, 0xcc, 0xbf, 0x8b, 0xb9, 0x2c, 0x48, 0xff, 0xbd, 0x40, 0x14, 0x00, 0x00, 0x06, 0x55, 0x53, 0x45, 0x52, 0x49, 0x44}
	if err := l.FromBytes(&payload, excepted_bs); nil != err {
		t.Error(err)
		return
	}

	assertEquals(t, "req.MessageTag", payload.MessageTag, uint8(0))
	assertEquals(t, "req.SessionID", payload.SessionID, uint32(50334208)) // MSessionID 0x00993554

	assertEquals(t, "req.Rand", payload.Rand, [16]byte{201, 64, 160, 237, 190, 129, 141, 204, 191, 139, 185, 44, 72, 255, 189, 64})
	assertEquals(t, "req.PrivLevel", payload.PrivLevel, uint8(20))

	assertEquals(t, "req.UserNameLength", payload.UserNameLength, uint8(6))
	assertEquals(t, "req.UserName", string(payload.UserName[:payload.UserNameLength]), "USERID")

	l.sequence--
	bs, err := l.ToBytes(&payload, make([]byte, 0, 1024))
	if nil != err {
		t.Error(err)
		return
	}

	if !bytes.Equal(bs, excepted_bs) {
		t.Error("excepted is", excepted_bs)
		t.Error("actual   is", bs)
		return
	}
}

var ipmi_v2_rakp_message_2 = []byte{0x06, 0x00, 0xff, 0x07, 0x06, 0x13, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3c, 0x00, 0x00, 0x00, 0x00, 0x00, 0xa4, 0xa3, 0xa2, 0xa0, 0x52, 0x65, 0x75, 0x19, 0x88, 0xb7, 0x3c, 0x5e, 0x42, 0xe8, 0x94, 0x7c, 0x25, 0xb1, 0x38, 0xb5, 0x34, 0x39, 0x34, 0x33, 0x32, 0x39, 0x43, 0x4e, 0x47, 0x30, 0x31, 0x33, 0x53, 0x36, 0x33, 0x34, 0xb8, 0x41, 0x46, 0x7a, 0xa6, 0x1f, 0x7e, 0xf4, 0xe1, 0x60, 0x0c, 0x85, 0x76, 0x1f, 0x07, 0xb2, 0x74, 0x54, 0x33, 0xf6}

func TestIPMIV2RakpMessage2(t *testing.T) {
	var l lanPlus
	var payload RakpMessage2

	l.authenticationAlgorithm = RAKPAlgorithmAuth_HMAC_SHA1
	l.integrityAlgorithm = RAKPAlgorithmIntegrity_HMAC_SHA1_96
	l.confidentialityAlgorithm = RAKPAlgorithmEncryto_AES_CBC_128

	if err := l.FromBytes(&payload, ipmi_v2_rakp_message_2); nil != err {
		t.Error(err)
		return
	}

	assertEquals(t, "req.MessageTag", payload.MessageTag, uint8(0))
	assertEquals(t, "req.StatusCode", payload.StatusCode, uint8(0))
	assertEquals(t, "req.SessionID", payload.SessionID, uint32(2695013284)) //Lsession_id 0xA0A2A3A4

	assertEquals(t, "req.Rand", payload.Rand, [16]byte{82, 101, 117, 25, 136, 183, 60, 94, 66, 232, 148, 124, 37, 177, 56, 181})
	assertEquals(t, "req.Guid", payload.Guid, [16]byte{52, 57, 52, 51, 50, 57, 67, 78, 71, 48, 49, 51, 83, 54, 51, 52})
	assertEquals(t, "req.KeyExchange", payload.KeyExchange, []byte{184, 65, 70, 122, 166, 31, 126, 244, 225, 96, 12, 133, 118, 31, 7, 178, 116, 84, 51, 246})

	l.sequence--
	bs, err := l.ToBytes(&payload, make([]byte, 0, 1024))
	if nil != err {
		t.Error(err)
		return
	}

	if !bytes.Equal(bs, ipmi_v2_rakp_message_2) {
		t.Error("excepted is", ipmi_v2_rakp_message_2)
		t.Error("actual   is", bs)
		return
	}
}

func TestIPMIV2RakpMessage2_22(t *testing.T) {
	var l lanPlus
	var payload RakpMessage2

	l.authenticationAlgorithm = RAKPAlgorithmAuth_HMAC_SHA1
	l.integrityAlgorithm = RAKPAlgorithmIntegrity_None
	l.confidentialityAlgorithm = RAKPAlgorithmEncryto_None

	excepted_bs := []byte{0x06, 0x00, 0xff, 0x07, 0x06, 0x13, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3c, 0x00, 0x00, 0x00, 0x00, 0x00, 0xa4, 0xa3, 0xa2, 0xa0, 0x11, 0xf1, 0x35, 0x6c, 0xa0, 0x04, 0xb3, 0xf3, 0xcc, 0x78, 0x4b, 0xf6, 0xb1, 0xc2, 0x49, 0xb8, 0xe1, 0xe9, 0xe1, 0x24, 0x9c, 0x08, 0x11, 0xe2, 0xa5, 0x00, 0x40, 0xf2, 0xe9, 0x21, 0x4f, 0x42, 0x84, 0x9f, 0xbe, 0x4c, 0xa5, 0x5d, 0x41, 0x30, 0xdf, 0xd9, 0x1a, 0x80, 0x40, 0x58, 0xa8, 0x5c, 0x4c, 0xb6, 0x1f, 0x3a}

	if err := l.FromBytes(&payload, excepted_bs); nil != err {
		t.Error(err)
		return
	}

	assertEquals(t, "req.MessageTag", payload.MessageTag, uint8(0))
	assertEquals(t, "req.StatusCode", payload.StatusCode, uint8(0))
	assertEquals(t, "req.SessionID", payload.SessionID, uint32(2695013284)) //Lsession_id 0xA0A2A3A4

	assertEquals(t, "req.Rand", payload.Rand, [16]byte{17, 241, 53, 108, 160, 4, 179, 243, 204, 120, 75, 246, 177, 194, 73, 184})
	assertEquals(t, "req.Guid", payload.Guid, [16]byte{225, 233, 225, 36, 156, 8, 17, 226, 165, 0, 64, 242, 233, 33, 79, 66})
	assertEquals(t, "req.KeyExchange", payload.KeyExchange, []byte{132, 159, 190, 76, 165, 93, 65, 48, 223, 217, 26, 128, 64, 88, 168, 92, 76, 182, 31, 58})

	l.sequence--
	bs, err := l.ToBytes(&payload, make([]byte, 0, 1024))
	if nil != err {
		t.Error(err)
		return
	}

	if !bytes.Equal(bs, excepted_bs) {
		t.Error("excepted is", excepted_bs)
		t.Error("actual   is", bs)
		return
	}
}

// >> rakp2 mac input buffer (71 bytes)
//  a4 a3 a2 a0 54 35 99 00 5e b1 18 c2 29 0a ca f0
//  b7 da f5 95 11 33 eb 9f a9 9e 46 6a f2 ff ea ab
//  d3 8b 6c 37 c4 ed a2 ae 34 39 34 33 32 39 43 4e
//  47 30 31 33 53 36 33 34 14 0d 41 64 6d 69 6e 69
//  73 74 72 61 74 6f 72
// >> rakp2 authcode (20 bytes)
//  31 32 33 34 35 36 61 62 63 00 00 00 00 00 00 00
//  00 00 00 00
// >> rakp2 mac as computed by the remote console (20 bytes)
//  ed b3 3e 4c 06 86 42 cf f9 c7 61 8a 1b 2e 6b f0
//  5d 76 0f 97
func TestRakp2HmacSha1(t *testing.T) {
	var w Writer
	lconsole_rand := [16]byte{94, 177, 24, 194, 41, 10, 202, 240, 183, 218, 245, 149, 17, 51, 235, 159}
	mconsole_rand := [16]byte{169, 158, 70, 106, 242, 255, 234, 171, 211, 139, 108, 55, 196, 237, 162, 174}
	mguid := [16]byte{52, 57, 52, 51, 50, 57, 67, 78, 71, 48, 49, 51, 83, 54, 51, 52}

	var sessesion_id uint32 = 2695013284  // remote id
	var msessesion_id uint32 = 0x00993554 // remote id

	w.Init(make([]byte, 0, 64))
	w.WriteUint32(sessesion_id)  // remote id
	w.WriteUint32(msessesion_id) // local id
	w.WriteBytes(lconsole_rand[:])
	w.WriteBytes(mconsole_rand[:])
	w.WriteBytes(mguid[:])
	w.WriteUint8(uint8(20))
	w.WriteUint8(uint8(13))
	w.WriteBytes([]byte("Administrator"))

	var authcode [16]byte
	copy(authcode[:], []byte("123456abc"))
	res := RAKPAlgorithmAuth_HMAC_SHA1.Gen(authcode[:], w.Bytes())

	excepted := []byte{237, 179, 62, 76, 6, 134, 66, 207, 249, 199, 97, 138, 27, 46, 107, 240, 93, 118, 15, 151}
	if !bytes.Equal(res, excepted) {
		t.Error("input    is", hex.EncodeToString(w.Bytes()))
		t.Error("authcode	 is", hex.EncodeToString(authcode[:]))

		t.Error("excepted is", hex.EncodeToString(excepted))
		t.Error("actual   is", hex.EncodeToString(res))
	}
}

func TestHmacMD5(t *testing.T) {

	var inputbuffer = []byte{0xa4, 0xa3, 0xa2, 0xa0, 0x54, 0x35, 0x99, 0x00, 0x42, 0x29, 0xd3, 0x2d, 0x30, 0xaa, 0xcf, 0xcf, 0x10, 0xf0, 0x7d, 0x35, 0xd7, 0x39, 0x71, 0xcb, 0x0f, 0xb6, 0xf4, 0x81, 0x0e, 0xe5, 0xd9, 0xea, 0xbe, 0xee, 0xf7, 0xa3, 0x06, 0xc1, 0x5b, 0xa2, 0x34, 0x39, 0x34, 0x33, 0x32, 0x39, 0x43, 0x4e, 0x47, 0x30, 0x31, 0x33, 0x53, 0x36, 0x33, 0x34, 0x14, 0x0d, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72}
	var authcode = []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x61, 0x62, 0x63, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	var actual = RAKPAlgorithmAuth_HMAC_SHA1.Gen(authcode, inputbuffer)

	excepted := []byte{0x4a, 0x29, 0xcc, 0xe0, 0x73, 0xf6, 0xe5, 0x50, 0x59, 0x09, 0xa9, 0xf1, 0x09, 0x8b, 0x86, 0xbb, 0x72, 0x97, 0x4e, 0xb3}
	if !bytes.Equal(actual, excepted) {
		t.Error("excepted is", excepted)
		t.Error("actual   is", actual)
	}
}

var ipmi_v2_rakp_message_3 = []byte{0x06, 0x00, 0xff, 0x07, 0x06, 0x14, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x54, 0x35, 0x99, 0x00, 0x32, 0x8f, 0xbb, 0x8f, 0xd6, 0x1e, 0xe1, 0x02, 0x78, 0x6d, 0x1f, 0xaa, 0x40, 0x08, 0x0c, 0x7a, 0x5e, 0x6a, 0x1e, 0xfb}

func TestIPMIV2RakpMessage3(t *testing.T) {
	var l lanPlus
	var payload RakpMessage3

	l.authenticationAlgorithm = RAKPAlgorithmAuth_HMAC_SHA1
	l.integrityAlgorithm = RAKPAlgorithmIntegrity_HMAC_SHA1_96
	l.confidentialityAlgorithm = RAKPAlgorithmEncryto_AES_CBC_128

	if err := l.FromBytes(&payload, ipmi_v2_rakp_message_3); nil != err {
		t.Error(err)
		return
	}

	assertEquals(t, "req.MessageTag", payload.MessageTag, uint8(0))
	assertEquals(t, "req.StatusCode", payload.StatusCode, uint8(0))
	assertEquals(t, "req.SessionID", payload.SessionID, uint32(0x00993554))

	assertEquals(t, "req.KeyExchange", payload.KeyExchange, []byte{50, 143, 187, 143, 214, 30, 225, 2, 120, 109, 31, 170, 64, 8, 12, 122, 94, 106, 30, 251})

	l.sequence--
	bs, err := l.ToBytes(&payload, make([]byte, 0, 1024))
	if nil != err {
		t.Error(err)
		return
	}

	if !bytes.Equal(bs, ipmi_v2_rakp_message_3) {
		t.Error("excepted is", ipmi_v2_rakp_message_3)
		t.Error("actual   is", bs)
		return
	}
}

func TestIPMIV2RakpMessage3_22(t *testing.T) {
	var l lanPlus
	var payload RakpMessage3

	l.authenticationAlgorithm = RAKPAlgorithmAuth_HMAC_SHA1
	l.integrityAlgorithm = RAKPAlgorithmIntegrity_None
	l.confidentialityAlgorithm = RAKPAlgorithmEncryto_None

	excepted_bs := []byte{0x06, 0x00, 0xff, 0x07, 0x06, 0x14, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0a, 0x00, 0x03, 0xe5, 0x73, 0x81, 0xb7, 0x05, 0xe6, 0x3f, 0xf6, 0xfc, 0x90, 0x6d, 0x40, 0x49, 0x4c, 0x18, 0x65, 0xe5, 0xca, 0x47, 0x9d}
	if err := l.FromBytes(&payload, excepted_bs); nil != err {
		t.Error(err)
		return
	}

	assertEquals(t, "req.MessageTag", payload.MessageTag, uint8(0))
	assertEquals(t, "req.StatusCode", payload.StatusCode, uint8(0))
	assertEquals(t, "req.SessionID", payload.SessionID, uint32(50334208))

	assertEquals(t, "req.KeyExchange", payload.KeyExchange, []byte{229, 115, 129, 183, 5, 230, 63, 246, 252, 144, 109, 64, 73, 76, 24, 101, 229, 202, 71, 157})

	l.sequence--
	bs, err := l.ToBytes(&payload, make([]byte, 0, 1024))
	if nil != err {
		t.Error(err)
		return
	}

	if !bytes.Equal(bs, excepted_bs) {
		t.Error("excepted is", excepted_bs)
		t.Error("actual   is", bs)
		return
	}
}

func TestRakp3HmacSha1(t *testing.T) {
	var w Writer
	mconsole_rand := [16]byte{82, 101, 117, 25, 136, 183, 60, 94, 66, 232, 148, 124, 37, 177, 56, 181}

	var session_id uint32 = 2695013284 // remote id

	w.Init(make([]byte, 0, 64))
	w.WriteBytes(mconsole_rand[:])
	w.WriteUint32(session_id) // remote id
	w.WriteUint8(uint8(20))   // privlevel

	w.WriteUint8(uint8(13))
	w.WriteBytes([]byte("Administrator"))

	var authcode [16]byte
	copy(authcode[:], []byte("123456abc"))
	res := RAKPAlgorithmAuth_HMAC_SHA1.Gen(authcode[:], w.Bytes())

	var inputbuffer = []byte{0x52, 0x65, 0x75, 0x19, 0x88, 0xb7, 0x3c, 0x5e, 0x42, 0xe8, 0x94, 0x7c, 0x25, 0xb1, 0x38, 0xb5, 0xa4, 0xa3, 0xa2, 0xa0, 0x14, 0x0d, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72}

	excepted := []byte{0x32, 0x8f, 0xbb, 0x8f, 0xd6, 0x1e, 0xe1, 0x02, 0x78, 0x6d, 0x1f, 0xaa, 0x40, 0x08, 0x0c, 0x7a, 0x5e, 0x6a, 0x1e, 0xfb}
	if !bytes.Equal(res, excepted) {
		t.Error("input       is", hex.EncodeToString(w.Bytes()))
		t.Error("excep_input is", hex.EncodeToString(inputbuffer))
		t.Error("authcode	   is", hex.EncodeToString(authcode[:]))

		t.Error("excepted    is", hex.EncodeToString(excepted))
		t.Error("actual      is", hex.EncodeToString(res))
	}
}

// rakp2 mac input buffer (71 bytes) - 1
// >> rakp2 mac input buffer (71 bytes)
//  a4 a3 a2 a0 54 35 99 00 45 7b 81 bb 81 b8 6c 94
//  68 97 11 3a b5 ff 3b 30 52 65 75 19 88 b7 3c 5e
//  42 e8 94 7c 25 b1 38 b5 34 39 34 33 32 39 43 4e
//  47 30 31 33 53 36 33 34 14 0d 41 64 6d 69 6e 69
//  73 74 72 61 74 6f 72
// >> rakp2 authcode (20 bytes)
//  31 32 33 34 35 36 61 62 63 00 00 00 00 00 00 00
//  00 00 00 00
// >> rakp2 mac as computed by the remote console (20 bytes)
//  b8 41 46 7a a6 1f 7e f4 e1 60 0c 85 76 1f 07 b2
//  74 54 33 f6
// >> rakp3 mac input buffer (35 bytes)
//  52 65 75 19 88 b7 3c 5e 42 e8 94 7c 25 b1 38 b5
//  a4 a3 a2 a0 14 0d 41 64 6d 69 6e 69 73 74 72 61
//  74 6f 72
// >> rakp3 mac key (20 bytes)
//  31 32 33 34 35 36 61 62 63 00 00 00 00 00 00 00
//  00 00 00 00
// generated rakp3 mac (20 bytes)
//  32 8f bb 8f d6 1e e1 02 78 6d 1f aa 40 08 0c 7a
//  5e 6a 1e fb
// 1
// session integrity key input (47 bytes)
//  45 7b 81 bb 81 b8 6c 94 68 97 11 3a b5 ff 3b 30
//  52 65 75 19 88 b7 3c 5e 42 e8 94 7c 25 b1 38 b5
//  14 0d 41 64 6d 69 6e 69 73 74 72 61 74 6f 72
// session authcode input (20 bytes)
//  31 32 33 34 35 36 61 62 63 00 00 00 00 00 00 00
//  00 00 00 00
// Generated session integrity key (20 bytes)
//  29 a3 aa 86 24 9a 8d 6b d5 a7 d9 59 73 39 37 5d
//  56 0f 38 06
// 1
// Generated K1 (20 bytes)
//  c5 4d dc 1b cc 14 a4 69 89 c6 f2 80 56 79 bb 58
//  3d 81 b1 c7
// Generated K2 (20 bytes)
//  35 f4 6e 58 25 02 de 15 19 81 f9 a0 52 d1 ee 2d
//  5a a7 3b d5
// >> rakp4 mac input buffer (36 bytes)
//  45 7b 81 bb 81 b8 6c 94 68 97 11 3a b5 ff 3b 30
//  54 35 99 00 34 39 34 33 32 39 43 4e 47 30 31 33
//  53 36 33 34
// >> rakp4 mac key (sik) (20 bytes)
//  29 a3 aa 86 24 9a 8d 6b d5 a7 d9 59 73 39 37 5d
//  56 0f 38 06
// >> rakp4 mac as computed by the BMC (20 bytes)
//  6b 0c 07 7a 02 03 8f 9a 65 89 16 65 25 b1 38 b5
//  34 39 34 33
// >> rakp4 mac as computed by the remote console (20 bytes)
//  6b 0c 07 7a 02 03 8f 9a 65 89 16 65 a3 9f 23 6d
//  ff da 40 cf
func TestRakp3HmacSha1Direct(t *testing.T) {

	var inputbuffer = []byte{0x52, 0x65, 0x75, 0x19, 0x88, 0xb7, 0x3c, 0x5e, 0x42, 0xe8, 0x94, 0x7c, 0x25, 0xb1, 0x38, 0xb5, 0xa4, 0xa3, 0xa2, 0xa0, 0x14, 0x0d, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72}
	var authcode = []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x61, 0x62, 0x63, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	var actual = RAKPAlgorithmAuth_HMAC_SHA1.Gen(authcode, inputbuffer)

	excepted := []byte{0x32, 0x8f, 0xbb, 0x8f, 0xd6, 0x1e, 0xe1, 0x02, 0x78, 0x6d, 0x1f, 0xaa, 0x40, 0x08, 0x0c, 0x7a, 0x5e, 0x6a, 0x1e, 0xfb}
	if !bytes.Equal(actual, excepted) {
		t.Error("excepted is", excepted)
		t.Error("actual   is", actual)
	}
}

var ipmi_v2_rakp_message_4 = []byte{0x06, 0x00, 0xff, 0x07, 0x06, 0x15, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x14, 0x00, 0x00, 0x00, 0x00, 0x00, 0xa4, 0xa3, 0xa2, 0xa0, 0x6b, 0x0c, 0x07, 0x7a, 0x02, 0x03, 0x8f, 0x9a, 0x65, 0x89, 0x16, 0x65}

func TestIPMIV2RakpMessage4(t *testing.T) {
	var l lanPlus
	var payload RakpMessage4

	l.authenticationAlgorithm = RAKPAlgorithmAuth_HMAC_SHA1
	l.integrityAlgorithm = RAKPAlgorithmIntegrity_HMAC_SHA1_96
	l.confidentialityAlgorithm = RAKPAlgorithmEncryto_AES_CBC_128

	if err := l.FromBytes(&payload, ipmi_v2_rakp_message_4); nil != err {
		t.Error(err)
		return
	}

	assertEquals(t, "req.MessageTag", payload.MessageTag, uint8(0))
	assertEquals(t, "req.StatusCode", payload.StatusCode, uint8(0))
	assertEquals(t, "req.SessionID", payload.SessionID, uint32(2695013284)) //0xA0A2A3A4

	assertEquals(t, "req.IntegrityCheck", payload.IntegrityCheck, []byte{107, 12, 7, 122, 2, 3, 143, 154, 101, 137, 22, 101})

	l.sequence--
	bs, err := l.ToBytes(&payload, make([]byte, 0, 1024))
	if nil != err {
		t.Error(err)
		return
	}

	if !bytes.Equal(bs, ipmi_v2_rakp_message_4) {
		t.Error("excepted is", ipmi_v2_rakp_message_4)
		t.Error("actual   is", bs)
		return
	}
}

func TestIPMIV2RakpMessage4_22(t *testing.T) {
	var l lanPlus
	var payload RakpMessage4

	l.authenticationAlgorithm = RAKPAlgorithmAuth_HMAC_SHA1
	l.integrityAlgorithm = RAKPAlgorithmIntegrity_HMAC_SHA1_96
	l.confidentialityAlgorithm = RAKPAlgorithmEncryto_AES_CBC_128

	excepted_bs := []byte{0x06, 0x00, 0xff, 0x07, 0x06, 0x15, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x14, 0x00, 0x00, 0x00, 0x00, 0x00, 0xa4, 0xa3, 0xa2, 0xa0, 0x48, 0xa3, 0x18, 0x44, 0x17, 0x36, 0xdc, 0x24, 0x2a, 0x3f, 0xc3, 0xe7}
	if err := l.FromBytes(&payload, excepted_bs); nil != err {
		t.Error(err)
		return
	}

	assertEquals(t, "req.MessageTag", payload.MessageTag, uint8(0))
	assertEquals(t, "req.StatusCode", payload.StatusCode, uint8(0))
	assertEquals(t, "req.SessionID", payload.SessionID, uint32(2695013284)) //0xA0A2A3A4

	assertEquals(t, "req.IntegrityCheck", payload.IntegrityCheck, []byte{72, 163, 24, 68, 23, 54, 220, 36, 42, 63, 195, 231})

	l.sequence--
	bs, err := l.ToBytes(&payload, make([]byte, 0, 1024))
	if nil != err {
		t.Error(err)
		return
	}

	if !bytes.Equal(bs, excepted_bs) {
		t.Error("excepted is", excepted_bs)
		t.Error("actual   is", bs)
		return
	}
}

// >> rakp4 mac input buffer (36 bytes)
//  45 7b 81 bb 81 b8 6c 94 68 97 11 3a b5 ff 3b 30
//  54 35 99 00 34 39 34 33 32 39 43 4e 47 30 31 33
//  53 36 33 34
// >> rakp4 mac key (sik) (20 bytes)
//  29 a3 aa 86 24 9a 8d 6b d5 a7 d9 59 73 39 37 5d
//  56 0f 38 06
// >> rakp4 mac as computed by the BMC (20 bytes)
//  6b 0c 07 7a 02 03 8f 9a 65 89 16 65 25 b1 38 b5
//  34 39 34 33
// >> rakp4 mac as computed by the remote console (20 bytes)
//  6b 0c 07 7a 02 03 8f 9a 65 89 16 65 a3 9f 23 6d
//  ff da 40 cf
// func TestRakp4HmacSha1Direct(t *testing.T) {
// 	var inputbuffer = []byte{0x45, 0x7b, 0x81, 0xbb, 0x81, 0xb8, 0x6c, 0x94, 0x68, 0x97, 0x11, 0x3a, 0xb5, 0xff, 0x3b, 0x30, 0x54, 0x35, 0x99, 0x00, 0x34, 0x39, 0x34, 0x33, 0x32, 0x39, 0x43, 0x4e, 0x47, 0x30, 0x31, 0x33, 0x53, 0x36, 0x33, 0x34}
// 	var authcode = []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x61, 0x62, 0x63, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
//
// 	var actual = RAKPAlgorithmAuth_HMAC_SHA1.Gen(authcode, inputbuffer)
//
// 	excepted := []byte{0x32, 0x8f, 0xbb, 0x8f, 0xd6, 0x1e, 0xe1, 0x02, 0x78, 0x6d, 0x1f, 0xaa, 0x40, 0x08, 0x0c, 0x7a, 0x5e, 0x6a, 0x1e, 0xfb}
// 	if !bytes.Equal(actual, excepted) {
// 		t.Error("excepted is", excepted)
// 		t.Error("actual   is", actual)
// 	}
// }

var ipmi_v2_AuthCapabilitiesRequest = []byte{0x06, 0x00, 0xff, 0x07, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x09, 0x20, 0x18, 0xc8, 0x81, 0x00, 0x38, 0x8e, 0x04, 0xb5}
var ipmi_v2_AuthCapabilitiesResponse = []byte{0x06, 0x00, 0xff, 0x07, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x81, 0x1c, 0x63, 0x20, 0x00, 0x38, 0x00, 0x02, 0x80, 0x14, 0x02, 0x00, 0x00, 0x00, 0x00, 0x10}

func handleAuthCapabilities(t *testing.T, conn *net.UDPConn, raddr net.Addr, bs []byte) {
	conn.WriteTo(ipmi_v2_AuthCapabilitiesResponse, raddr)
}
func handleOpenRequest(t *testing.T, conn *net.UDPConn, raddr net.Addr, bs []byte) {
	if !bytes.Equal(bs, ipmi_v2_open_session_request) {
		t.Error("excepted is", len(ipmi_v2_open_session_request), hex.EncodeToString(ipmi_v2_open_session_request))
		t.Error("actual   is", len(bs), hex.EncodeToString(bs))
	}
	conn.WriteTo(ipmi_v2_open_session_response, raddr)
}
func handleRAKP1(t *testing.T, conn *net.UDPConn, raddr net.Addr, bs []byte) {
	if !bytes.Equal(bs, ipmi_v2_rakp_message_1) {
		t.Error("excepted is", len(ipmi_v2_rakp_message_1), hex.EncodeToString(ipmi_v2_rakp_message_1))
		t.Error("actual   is", len(bs), hex.EncodeToString(bs))
	}

	conn.WriteTo(ipmi_v2_rakp_message_2, raddr)
}
func handleRAKP3(t *testing.T, conn *net.UDPConn, raddr net.Addr, bs []byte) {
	if !bytes.Equal(bs, ipmi_v2_rakp_message_3) {
		t.Error("excepted is", len(ipmi_v2_rakp_message_3), hex.EncodeToString(ipmi_v2_rakp_message_3))
		t.Error("actual   is", len(bs), hex.EncodeToString(bs))
	}

	conn.WriteTo(ipmi_v2_rakp_message_4, raddr)
}

var ipmi_v2_set_privilege_level_request = []byte{0x06, 0x00, 0xff, 0x07, 0x06, 0x00, 0x54, 0x35, 0x99, 0x00, 0x03, 0x00, 0x00, 0x00, 0x08, 0x00, 0x20, 0x18, 0xc8, 0x81, 0x04, 0x3b, 0x04, 0x3c}
var ipmi_v2_set_privilege_level_response = []byte{0x06, 0x00, 0xff, 0x07, 0x06, 0x00, 0xa4, 0xa3, 0xa2, 0xa0, 0x01, 0x00, 0x00, 0x00, 0x09, 0x00, 0x81, 0x1c, 0x63, 0x20, 0x04, 0x3b, 0x00, 0x04, 0x9d}

func handleSetPrivilegeLevel(t *testing.T, conn *net.UDPConn, raddr net.Addr, bs []byte) {
	if !bytes.Equal(bs, ipmi_v2_set_privilege_level_request) {
		t.Error("excepted is", len(ipmi_v2_set_privilege_level_request), hex.EncodeToString(ipmi_v2_set_privilege_level_request))
		t.Error("actual   is", len(bs), hex.EncodeToString(bs))
	}

	conn.WriteTo(ipmi_v2_set_privilege_level_response, raddr)
}

func TestOpenSession(t *testing.T) {
	conn, e := net.ListenUDP("udp", &net.UDPAddr{Port: 0})
	if nil != e {
		t.Error(e)
		return
	}

	var wait sync.WaitGroup
	var is_running int32 = 1
	go func() {
		count := 0
		defer wait.Done()
		bs := make([]byte, 1024)
		for 1 == atomic.LoadInt32(&is_running) {
			n, raddr, e := conn.ReadFrom(bs)
			if nil != e {
				fmt.Println(e)
				break
			}

			switch count {
			case 0:
				t.Log("handleAuthCapabilities")
				handleAuthCapabilities(t, conn, raddr, bs[:n])
			case 1:
				t.Log("handleOpenRequest")
				handleOpenRequest(t, conn, raddr, bs[:n])
			case 2:
				t.Log("handleRAKP1")
				handleRAKP1(t, conn, raddr, bs[:n])
			case 3:
				t.Log("handleRAKP3")
				handleRAKP3(t, conn, raddr, bs[:n])
			case 4:
				t.Log("handleSetPrivilegeLevel")
				handleSetPrivilegeLevel(t, conn, raddr, bs[:n])
			default:
				t.Error("unknow message")
			}
			count++
		}
	}()
	wait.Add(1)

	defer func() {
		atomic.StoreInt32(&is_running, 0)
		conn.Close()

		wait.Wait()
	}()

	_, sport, _ := net.SplitHostPort(conn.LocalAddr().String())
	port, _ := strconv.Atoi(sport)

	l := &lanPlus{
		lanBase: lanBase{conn_opt: &ConnectionOption{Hostname: "127.0.0.1",
			Port:      port,
			Username:  "Administrator",
			Password:  "123456abc",
			Interface: "lanplus"},
		},
		//PrivLevel:                20,
		RoleUserOnlyLookup:       true,
		authenticationAlgorithm:  RAKPAlgorithmAuth_HMAC_SHA1,
		integrityAlgorithm:       RAKPAlgorithmIntegrity_HMAC_SHA1_96,
		confidentialityAlgorithm: RAKPAlgorithmEncryto_AES_CBC_128,
	}

	copy(l.Authcode[:], []byte(l.lanBase.conn_opt.Password))
	copy(l.Username[:], []byte(l.lanBase.conn_opt.Username))
	l.UsernameLen = uint8(len(l.lanBase.conn_opt.Username))

	copy(l.LRand[:], []byte{69, 123, 129, 187, 129, 184, 108, 148, 104, 151, 17, 58, 181, 255, 59, 48})
	_session_id = 2695013284 - 1

	is_test = true
	if e := l.open(); nil != e {
		t.Error(e)
	}
}

// GetDeviceId Request
// 0000   d8 d3 85 a9 86 0e 40 8d 5c 5e e6 dc 08 00 45 00
// 0010   00 33 6d e2 00 00 80 11 47 e3 c0 a8 01 c8 c0 a8
// 0020   01 dc dd 0f 02 6f 00 1f 31 63 06 00 ff 07 06 00
// 0030   54 35 99 00 04 00 00 00 07 00 20 18 c8 81 08 01
// 0040   76
//var ipmi_v2_get_device_id_request = []byte{0x06, 0x00, 0xff, 0x07, 0x06, 0x00, 0x54, 0x35, 0x99, 0x00, 0x04, 0x00, 0x00, 0x00, 0x07, 0x00, 0x20, 0x18, 0xc8, 0x81, 0x08, 0x01, 0x76}

// func TestIPMIV2GetDeviceIdRequest(t *testing.T) {
// 	var l lanPlus
// 	var request DeviceIDRequest
// 	var payload = NewRequest(NetworkFunctionApp, CommandGetDeviceID, &request)

// 	if err := l.FromBytes(payload, ipmi_v2_get_device_id_request); nil != err {
// 		t.Error(err)
// 		return
// 	}

// 	l.sequence--
// 	bs, err := l.ToBytes(payload, make([]byte, 0, 1024))
// 	if nil != err {
// 		t.Error(err)
// 		return
// 	}

// 	if !bytes.Equal(bs, ipmi_v2_get_device_id_request) {
// 		t.Error("excepted is", ipmi_v2_get_device_id_request)
// 		t.Error("actual   is", bs)
// 		return
// 	}
// }

// GetDeviceId Response
// 0000   40 8d 5c 5e e6 dc d8 d3 85 a9 86 0e 08 00 45 00
// 0010   00 3f 33 aa 00 00 40 11 c2 0f c0 a8 01 dc c0 a8
// 0020   01 c8 02 6f dd 0f 00 2b b0 59 06 00 ff 07 06 00
// 0030   a4 a3 a2 a0 02 00 00 00 13 00 81 1c 63 20 08 01
// 0040   00 11 81 02 05 02 0f 0b 00 00 00 20 02

func TestOpenSession(t *testing.T) {

}
