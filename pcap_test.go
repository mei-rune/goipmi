package goipmi

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/google/gopacket/pcap"
	"github.com/runner-mei/goipmi/protocol"
)

func TestDecodeMessage(t *testing.T) {
	//debug = true
	handle, e := pcap.OpenOffline("encrypto_aes_cbc.pcap")
	if nil != e {
		t.Error(e)
		return
	}

	var authenticationAlgorithm protocol.RAKPAlgorithmAuthMethod
	var integrityAlgorithm protocol.RAKPAlgorithmIntegrityMethod
	var confidentialityAlgorithm protocol.RAKPAlgorithmConfidentialityMethod

	var rakp1 *protocol.RakpMessage1
	var rakp2 *protocol.RakpMessage2
	var rakp3 *protocol.RakpMessage3
	var rakp4 *protocol.RakpMessage4

	Authcode := []byte("123456abc")
	var sik, k1, k2 []byte

	count := 0
	if e = handlePacketData(handle, func(bs []byte) error {

		var rmcpHeader protocol.RMCPHeader
		var ipmiHeader protocol.IPMIV2Header

		var r protocol.Reader
		r.Init(bs)

		count++

		rmcpHeader.ReadBytes(&r)
		if nil != r.Err() {
			return errors.New("Read IPMI header failed, " + r.Err().Error())
		}

		// if rmcpHeader.Class != rmcpClassIPMI {
		// 	return errors.New("rmcp class(" +
		// 		strconv.FormatInt(int64(rmcpHeader.Class), 10) +
		// 		") isn't ipmi.")
		// }

		ipmiHeader.ReadBytes(&r)
		if nil != r.Err() {
			return errors.New("Read IPMI header failed, " + r.Err().Error())
		}

		switch protocol.PayloadType(ipmiHeader.PayloadType).Value() {
		// case commands.PayloadIPMI:
		// 	var msg Request
		// 	msg.ReadBytes(&r)
		// 	if r.Err() != nil {
		// 		fmt.Println("======", r.Err())
		// 	}
		// 	fmt.Println(msg)
		case protocol.PayloadIPMI:
			bs := r.Bytes()
			bs = bs[ipmiHeader.Length:]
			for i := 0; ; i++ {
				if bs[i] != 0xff {
					bs = bs[i:]
					break
				}
			}
			bs = bs[1:] // pad length
			bs = bs[1:] // next header

			fmt.Printf("===== %d %x\r\n", count, bs)

			// var msg interface {
			// 	ReadBytes(r *protocol.Reader)
			// }
			// switch r.Bytes()[5] {
			// default:
			// 	msg = &protocol.Request{}
			// }

			// //bs := r.Bytes()
			// if msg == nil {
			// 	fmt.Println("===== skipped")
			// } else {
			// 	msg.ReadBytes(&r)
			// 	if r.Err() != nil {
			// 		fmt.Println("======", r.Err())
			// 	}
			// 	fmt.Println(msg)
			// 	fmt.Printf("====== %d %x\r\n", count, bs[ipmiHeader.Length:])
			// }
		case protocol.PayloadOpenSessionRequest:
			var msg protocol.OpenSessionRequest
			msg.ReadBytes(&r)
			if r.Err() != nil {
				fmt.Println("======", r.Err())
			}
			fmt.Println(&msg)

		case protocol.PayloadOpenSessionResponse:
			var msg protocol.OpenSessionResponse
			msg.ReadBytes(&r)
			if r.Err() != nil {
				fmt.Println("======", r.Err())
			}

			authenticationAlgorithm = msg.Authentication.Algorithm
			integrityAlgorithm = msg.Integrity.Algorithm
			confidentialityAlgorithm = msg.Confidentiality.Algorithm

			fmt.Println(&msg)
		case protocol.PayloadRAKPMessage1:
			var msg protocol.RakpMessage1
			msg.ReadBytes(&r)
			if r.Err() != nil {
				fmt.Println("======", r.Err())
			}
			fmt.Println(&msg)
			rakp1 = &msg
		case protocol.PayloadRAKPMessage2:
			var msg protocol.RakpMessage2
			msg.ReadBytes(&r, authenticationAlgorithm.Size())
			if r.Err() != nil {
				fmt.Println("======", r.Err())
			}
			fmt.Println(&msg)
			rakp2 = &msg

			var w protocol.Writer
			w.Init(make([]byte, 0, 64))
			//w.Reset()
			w.WriteBytes(rakp1.Rand[:])
			w.WriteBytes(rakp2.Rand[:])
			w.WriteUint8(rakp1.PrivLevel) // privlevel
			w.WriteUint8(rakp1.UserNameLength)
			w.WriteBytes(rakp1.UserName[:int(rakp1.UserNameLength)])

			//fmt.Println(hex.EncodeToString(w.Bytes()))
			sik = protocol.RAKPAlgorithmAuth_HMAC_SHA1.Gen(Authcode[:], w.Bytes())

			fmt.Printf("sik=%x\r\n", sik)
			// bytes.Repeat([]byte{1}, 20)

			k1 = protocol.RAKPAlgorithmIntegrity_HMAC_SHA1_96.Gen(sik, bytes.Repeat([]byte{1}, 20))
			fmt.Printf("k1=%x\r\n", k1)

			k2 = protocol.RAKPAlgorithmIntegrity_HMAC_SHA1_96.Gen(sik, bytes.Repeat([]byte{2}, 20))
			fmt.Printf("k2=%x\r\n", k2)

		case protocol.PayloadRAKPMessage3:
			var msg protocol.RakpMessage3
			msg.ReadBytes(&r, authenticationAlgorithm.Size())
			if r.Err() != nil {
				fmt.Println("======", r.Err())
			}
			fmt.Println(&msg)
			rakp3 = &msg
		case protocol.PayloadRAKPMessage4:
			var msg protocol.RakpMessage4
			msg.ReadBytes(&r, integrityAlgorithm.KeyExchangeSize())
			if r.Err() != nil {
				fmt.Println("======", r.Err())
			}
			fmt.Println(&msg)
			rakp4 = &msg
		default:
			fmt.Println(errors.New(fmt.Sprintf("payload(%d) is unsupported", protocol.PayloadType(ipmiHeader.PayloadType))))
		}

		return nil
	}); nil != e {
		t.Error(e)
	}

	var header protocol.IPMIV2Header
	var levelReq protocol.SessionPrivilegeLevelRequest
	var req protocol.Request
	req.Data = &levelReq
	txt := []byte{0x06, 0x40, 0x54, 0x35, 0x99, 0x00, 0x03, 0x00, 0x00, 0x00, 0x08, 0x00, 0x20, 0x18, 0xc8, 0x81,
		0x04, 0x3b, 0x04, 0x3c, 0xff, 0xff, 0x02, 0x07, 0x99, 0xe9, 0xdf, 0x35, 0x49, 0x7d, 0xe5, 0x52,
		0x0d, 0xa4, 0xc0, 0x83}

	var reader protocol.Reader
	reader.Init(txt)

	header.ReadBytes(&reader)
	reader.Read(&req)

	if reader.Err() != nil {
		t.Error(reader.Err())
		return
	}

	bs := reader.Bytes()
	fmt.Printf("%d   -    %x\r\n", header.Length, bs)

	for i := 0; ; i++ {
		if bs[i] != 0xff {
			bs = bs[i:]
			break
		}
	}

	bs = bs[1:] // pad length
	bs = bs[1:] // next header

	fmt.Printf("%x\r\n", bs)

	// fix vet
	fmt.Println(confidentialityAlgorithm)
	fmt.Println(rakp3)
	fmt.Println(rakp4)
}
