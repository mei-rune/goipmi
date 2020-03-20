package goipmi

import (
	"testing"

	"github.com/runner-mei/goipmi/protocol"
)

// ipmi_cmd SDR[f] off=48 ilen=16 status=0 cc=0 sz=18
// GetSDR[000f] next=10 (len=64): 0f 00 51 01 3b 20 00 0d 27 01 23 c9 01 01 00 0a 0
// 0 60 30 00 80 01 00 00 01 00 00 00 00 00 00 00 00 00 7f 81 2d 29 27 00 00 00 00
// 00 00 00 00 c6 54 65 6d 70 20 31 00 00 00 00 00 00 00 00 00 00
// GetSDR[000f]: ret = 0, next=10
// got SDR (len=64):
//   0000: 0f 00 51 01 3b 20 00 0d 27 01 23 c9 01 01 00 0a
//   0010: 00 60 30 00 80 01 00 00 01 00 00 00 00 00 00 00
//   0020: 00 00 7f 81 2d 29 27 00 00 00 00 00 00 00 00 c6
//   0030: 54 65 6d 70 20 31 00 00 00 00 00 00 00 00 00 00
// ShowSDR: len=64, type=1
// raw SDR: 0f 00 51 01 3b 20 00 0d 27 01 23 c9 01 01 00 0a 00 60 30 00 80 01 00 00 01 00 00 00 00 00 00 00 00 00 7f 81 2d 29 27 00 00 00 00 00 00 00 00 c6 54 65 6d 70 20 31 00 00 00 00 00 00 00 00 00 00
// SDR[f] Full ioff=48 idTypLen=0x54 ilen=16
// entity 39.1, idlen=16 sizeof=64 idstr0=T s0=54
//
// >> Sending IPMI command payload
// >>    netfn   : 0x04
// >>    command : 0x2d
// >>    data_len: 1
// >>    data    : 0x0d
// BUILDING A v2 COMMAND
// added list entry seq=0x0c cmd=0x2d
// added list entry seq=0x0c cmd=0x34
// >> sending packet (24 bytes)
//  06 00 ff 07 06 00 26 36 99 00 4e 00 00 00 08 00 20 10 d0 81 30 2d 0d 15
// send_payload(non-SOL) type=0 data
// << received packet (28 bytes)
//  06 00 ff 07 06 00 a4 a3 a2 a0 4c 00 00 00 0c 00 81 14 6b 20 30 2d 00 18 c0 00 80 2b
// rsp session_id=a0a2a3a4 session_seq=0000004c
// << IPMI Response Session Header
// <<   Authtype                : Unknown (0x6)
// <<   Payload type            : IPMI (0)
// <<   Session ID              : 0xa0a2a3a4
// <<   Sequence                : 0x0000004c
// <<   IPMI Msg/Payload Length : 12
// << IPMI Response Message Header
// <<   Rq Addr    : 81
// <<   NetFn      : 05
// <<   Rq LUN     : 0
// <<   Rs Addr    : 20
// <<   Rq Seq     : 0c
// <<   Rs Lun     : 0
// <<   Command    : 2d
// <<   Compl Code : 0x00
// IPMI Request Match found
// removed list entry seq=0x0c cmd=0x2d
// send_payload(non-SOL) rsp dlen=4, rs_seq=76
// GetSensorReading mc=0,20,0 status=0 cc=0 sz=4 resp: 18 c0 00 80
// bitnum(00)=0 raw=18 init=c0 base/units=1/80
// units=80 base=1 mod=0 (raw=18, nom_rd=0)
// decode1: m=1 b=0 b_exp=0 rx=0, a=0 ax=0 l=0, floatval=24.000000
// get_unit_type(80,1,0,0)
// 000f SDR Full 01 01 20 a 01 snum 0d Temp 1 = 18 OK   24.00 degrees C
//
func TestFullRecord(t *testing.T) {
	raw_data := []byte{0x0f, 0x00, 0x51, 0x01, 0x3b, 0x20, 0x00, 0x0d, 0x27, 0x01, 0x23, 0xc9, 0x01, 0x01, 0x00, 0x0a, 0x00, 0x60, 0x30, 0x00, 0x80, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x7f, 0x81, 0x2d, 0x29, 0x27, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc6, 0x54, 0x65, 0x6d, 0x70, 0x20, 0x31, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	var full FullSensorRecord
	r := protocol.NewReader(raw_data)
	full.ReadBytes(r)
	if r.Err() != nil {
		t.Error(r.Err())
		return
	}

	v, e := full.Calc(24, 8)
	if e != nil {
		t.Error(e)
		return
	}

	if 24 != int(v) {
		t.Error(v)
	}
}
