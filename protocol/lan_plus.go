package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/runner-mei/goipmi/protocol/commands"
)

/*
 * ipmi_lanplus_build_v2x_msg
 *
 * Encapsulates the payload data to create the IPMI v2.0 / RMCP+ packet.
 *
 *
 * IPMI v2.0 LAN Request Message Format
 * +----------------------+
 * |  rmcp.ver            | 4 bytes
 * |  rmcp.__rsvd         |
 * |  rmcp.seq            |
 * |  rmcp.class          |
 * +----------------------+
 * |  session.authtype    | 10 bytes
 * |  session.payloadtype |
 * |  session.id          |
 * |  session.seq         |
 * +----------------------+
 * |  message length      | 2 bytes
 * +----------------------+
 * | Confidentiality Hdr  | var (possibly absent)
 * +----------------------+
 * |  Payload             | var Payload
 * +----------------------+
 * | Confidentiality Trlr | var (possibly absent)
 * +----------------------+
 * | Integrity pad        | var (possibly absent)
 * +----------------------+
 * | Pad length           | 1 byte (WTF?)
 * +----------------------+
 * | Next Header          | 1 byte (WTF?)
 * +----------------------+
 * | Authcode             | var (possibly absent)
 * +----------------------+
 */
type lanPlus struct {
	lanBase

	active bool

	RoleUserOnlyLookup bool
	PrivLevel          uint8
	Authcode           [16]byte
	UsernameLen        uint8
	Username           [16]byte
	LRand              [16]byte
	RRand              [16]byte
	RGuid              [16]byte
	LSessionID         uint32
	RSessionID         uint32

	authenticationAlgorithm  RAKPAlgorithmAuthMethod
	integrityAlgorithm       RAKPAlgorithmIntegrityMethod
	confidentialityAlgorithm RAKPAlgorithmConfidentialityMethod
	sik, k1, k2              []byte
}

func (l *lanPlus) encrypto(bs []byte) ([]byte, error) {
	return l.confidentialityAlgorithm.Encrypt(l.k2)(bs, make([]byte, 32+len(bs)))
}

func (l *lanPlus) unencrypto(bs []byte) ([]byte, error) {
	return l.confidentialityAlgorithm.Decrypt(l.k2)(bs, make([]byte, 32+len(bs)))
}

func (l *lanPlus) FromBytes(payload interface{}, bs []byte) error {
	var rmcpHeader RMCPHeader
	var ipmiHeader IPMIV2Header

	var r Reader
	r.Init(bs)

	rmcpHeader.ReadBytes(&r)
	if nil != r.Err() {
		return errors.New("Read IPMI header failed, " + r.Err().Error())
	}

	if rmcpHeader.Class != rmcpClassIPMI {
		return errors.New("rmcp class(" +
			strconv.FormatInt(int64(rmcpHeader.Class), 10) +
			") isn't ipmi.")
	}

	ipmiHeader.ReadBytes(&r)
	if nil != r.Err() {
		return errors.New("Read IPMI header failed, " + r.Err().Error())
	}

	if ipmiHeader.AuthType != AuthTypeFormatIPMIV2 {
		return errors.New("ipmi version(" +
			strconv.FormatInt(int64(ipmiHeader.AuthType), 10) +
			") isn't v2.")
	}

	switch resp := payload.(type) {
	case *Request:
		if PayloadType(ipmiHeader.PayloadType) != PayloadIPMI {
			return errors.New("ipmi payload(" +
				strconv.FormatInt(int64(ipmiHeader.PayloadType), 10) +
				") isn't IPMI Request.")
		}

		if PayloadType(ipmiHeader.PayloadType).Encryption() {
			encrypted := r.ReadBytes(int(ipmiHeader.Length))
			encrypted, err := l.unencrypto(encrypted)
			if err != nil {
				return errors.New(fmt.Sprintf("unencrypto payload(%T) failed, ", payload) + err.Error())
			}
			var rd Reader
			rd.Init(encrypted)
			resp.ReadBytes(&rd)
			if rd.Err() != nil {
				return errors.New(fmt.Sprintf("read payload(%T) failed, ", payload) + rd.Err().Error())
			}
		} else {
			resp.ReadBytes(&r)
		}
	case *Response:
		if PayloadType(ipmiHeader.PayloadType).Value() != PayloadIPMI {
			return errors.New("ipmi payload(" +
				strconv.FormatInt(int64(ipmiHeader.PayloadType), 10) +
				") isn't IPMI Response.")
		}

		if PayloadType(ipmiHeader.PayloadType).Encryption() {
			encrypted := r.ReadBytes(int(ipmiHeader.Length))
			encrypted, err := l.unencrypto(encrypted)
			if err != nil {
				return errors.New(fmt.Sprintf("unencrypto payload(%T) failed, ", payload) + err.Error())
			}
			var rd Reader
			rd.Init(encrypted)
			resp.ReadBytes(&rd)
			if rd.Err() != nil {
				return errors.New(fmt.Sprintf("read payload(%T) failed, ", payload) + rd.Err().Error())
			}
		} else {
			resp.ReadBytes(&r)
		}
	case *OpenSessionRequest:
		if PayloadType(ipmiHeader.PayloadType).Value() != PayloadOpenSessionRequest {
			return errors.New("ipmi payload(" +
				strconv.FormatInt(int64(ipmiHeader.PayloadType), 10) +
				") isn't OpenSessionRequest.")
		}
		resp.ReadBytes(&r)
	case *OpenSessionResponse:
		if PayloadType(ipmiHeader.PayloadType).Value() != PayloadOpenSessionResponse {
			return errors.New("ipmi payload(" +
				strconv.FormatInt(int64(ipmiHeader.PayloadType), 10) +
				") isn't OpenSessionResponse.")
		}
		if bs := r.Bytes(); RAKP_STATUS_NO_ERRORS != StatusCode(bs[1]) {
			resp.StatusCode = bs[1]
			return nil
		}
		resp.ReadBytes(&r)
	case *RakpMessage1:
		if PayloadType(ipmiHeader.PayloadType).Value() != PayloadRAKPMessage1 {
			return errors.New("ipmi payload(" +
				strconv.FormatInt(int64(ipmiHeader.PayloadType), 10) +
				") isn't RakpMessage1.")
		}
		resp.ReadBytes(&r)
	case *RakpMessage2:
		if PayloadType(ipmiHeader.PayloadType).Value() != PayloadRAKPMessage2 {
			return errors.New("ipmi payload(" +
				strconv.FormatInt(int64(ipmiHeader.PayloadType), 10) +
				") isn't RakpMessage2.")
		}
		if bs := r.Bytes(); RAKP_STATUS_NO_ERRORS != StatusCode(bs[1]) {
			resp.StatusCode = bs[1]
			return nil
		}
		resp.ReadBytes(&r, l.authenticationAlgorithm.Size())
	case *RakpMessage3:
		if PayloadType(ipmiHeader.PayloadType).Value() != PayloadRAKPMessage3 {
			return errors.New("ipmi payload(" +
				strconv.FormatInt(int64(ipmiHeader.PayloadType), 10) +
				") isn't RakpMessage3.")
		}
		if bs := r.Bytes(); RAKP_STATUS_NO_ERRORS != StatusCode(bs[1]) {
			resp.StatusCode = bs[1]
			return nil
		}
		resp.ReadBytes(&r, l.authenticationAlgorithm.Size())
	case *RakpMessage4:
		if PayloadType(ipmiHeader.PayloadType).Value() != PayloadRAKPMessage4 {
			return errors.New("ipmi payload(" +
				strconv.FormatInt(int64(ipmiHeader.PayloadType), 10) +
				") isn't RakpMessage4.")
		}
		if bs := r.Bytes(); RAKP_STATUS_NO_ERRORS != StatusCode(bs[1]) {
			resp.StatusCode = bs[1]
			return nil
		}
		resp.ReadBytes(&r, l.integrityAlgorithm.KeyExchangeSize())
	default:
		return errors.New(fmt.Sprintf("payload(%T) is unsupported", payload))
	}

	if nil != r.Err() {
		return errors.New(fmt.Sprintf("read payload(%T) failed, ", payload) + r.Err().Error())
	}

	return nil
}

func (l *lanPlus) ToBytes(payload interface{}, bs []byte) ([]byte, error) {
	rmcpHeader := RMCPHeader{
		Version:            rmcpVersion1,
		Class:              rmcpClassIPMI,
		RMCPSequenceNumber: 0xff,
	}
	ipmiHeader := IPMIV2Header{
		//PayloadType:  uint8(payloadType),
		//OemIana:      oemIana,
		//OemPayloadId: oemPayloadId,
		SessionID: l.RSessionID, //l.SessionID,
	}

	var reqW Writable
	w := Writer{}
	w.Init(bs)
	rmcpHeader.WriteBytes(&w)

	authenticated := false
	encrypted := false
	switch req := payload.(type) {
	case *Request:
		ipmiHeader.PayloadType = uint8(PayloadIPMI)
		ipmiHeader.Sequence = l.nextSequence()
		req.Body.RqSeq = l.nextRqSequence()
		reqW = req

		if l.integrityAlgorithm != RAKPAlgorithmIntegrity_None {
			authenticated = true
			ipmiHeader.PayloadType = ipmiHeader.PayloadType | 0x40
		}

		if l.confidentialityAlgorithm != RAKPAlgorithmEncryto_None {
			encrypted = true
			ipmiHeader.PayloadType = ipmiHeader.PayloadType | 0x80
		}

	case *Response:
		ipmiHeader.PayloadType = uint8(PayloadIPMI)
		ipmiHeader.Sequence = l.nextSequence()
		req.Body.RqSeq = l.nextRqSequence()
		reqW = req

		if l.integrityAlgorithm != RAKPAlgorithmIntegrity_None {
			authenticated = true
			ipmiHeader.PayloadType = ipmiHeader.PayloadType | 0x40
		}

		if l.confidentialityAlgorithm != RAKPAlgorithmEncryto_None {
			encrypted = true
			ipmiHeader.PayloadType = ipmiHeader.PayloadType | 0x80
		}

	case *OpenSessionRequest:
		ipmiHeader.SessionID = 0
		ipmiHeader.Sequence = 0
		ipmiHeader.PayloadType = uint8(PayloadOpenSessionRequest)
		req.PrivLevel = 0
		reqW = req
	case *OpenSessionResponse:
		ipmiHeader.SessionID = 0
		ipmiHeader.Sequence = 0
		ipmiHeader.PayloadType = uint8(PayloadOpenSessionResponse)
		reqW = req
	case *RakpMessage1:
		ipmiHeader.SessionID = 0
		ipmiHeader.Sequence = 0
		ipmiHeader.PayloadType = uint8(PayloadRAKPMessage1)
		reqW = req
	case *RakpMessage2:
		ipmiHeader.SessionID = 0
		ipmiHeader.Sequence = 0
		ipmiHeader.PayloadType = uint8(PayloadRAKPMessage2)
		reqW = req
	case *RakpMessage3:
		ipmiHeader.SessionID = 0
		ipmiHeader.Sequence = 0
		ipmiHeader.PayloadType = uint8(PayloadRAKPMessage3)
		reqW = req
	case *RakpMessage4:
		ipmiHeader.SessionID = 0
		ipmiHeader.Sequence = 0
		ipmiHeader.PayloadType = uint8(PayloadRAKPMessage4)
		reqW = req
	default:
		return nil, errors.New(fmt.Sprintf("payload(%T) is unsupported", payload))
	}

	unauthenticatedLen := w.Len()
	ipmiHeader.WriteBytes(&w)
	if nil != w.Err() {
		return nil, w.Err()
	}
	oldLength := w.Len()

	reqW.WriteBytes(&w)
	if nil != w.Err() {
		return nil, w.Err()
	}
	tmp := w.Bytes()

	if encrypted {
		encryptedBytes, err := l.encrypto(tmp[oldLength:])
		if err != nil {
			return nil, errors.New(fmt.Sprintf("read payload(%T) fail, ", payload) + err.Error())
		}
		w.Truncate(oldLength)
		w.WriteBytes(encryptedBytes)
		tmp = w.Bytes()
	}
	binary.LittleEndian.PutUint16(tmp[oldLength-2:], uint16(w.Len()-oldLength))

	if authenticated {
		inputBytes := tmp[unauthenticatedLen:]
		fmt.Printf("input bytes %x\r\n", inputBytes)
		inputLen := len(inputBytes) + l.authenticationAlgorithm.Size() + 2 // 2 is 为 pad 数据长度占用一个字节和 nextheader 占用一个字节

		padLen := 4 - (inputLen % 4)
		if padLen < 4 {
			for i := 0; i < padLen; i++ {
				w.WriteUint8(0xff)
			}
			w.WriteUint8(uint8(padLen)) // pad length
		} else {
			w.WriteUint8(0) // padLen
		}
		w.WriteUint8(7) // next header

		tmp = w.Bytes()
		inputBytes = tmp[unauthenticatedLen:]

		fmt.Printf("k1 %x\r\n", l.k1)
		fmt.Printf("input bytes %x\r\n", inputBytes)
		w.WriteBytes(l.integrityAlgorithm.Gen(l.k1, inputBytes)[:l.integrityAlgorithm.KeyExchangeSize()])

		tmp = w.Bytes()
	}

	return tmp, nil
}

func (l *lanPlus) open() error {
	err := l.dial()
	if err != nil {
		return err
	}

	// // TODO: options
	// l.priv = PrivLevelAdmin
	// l.timeout = time.Second * 5
	// l.lun = 0

	if err := l.openSession(); err != nil {
		l.close()
		return err
	}
	return nil
}

func (l *lanPlus) close() error {
	if l.active {
		if err := l.closeSession(); err != nil {
			log.Println("error closing session:", err)
		}
		l.active = false
	}

	if l.conn != nil {
		_ = l.conn.Close()
		l.conn = nil
	}

	return nil
}

func (l *lanPlus) send(req, resp interface{}) error {
	bsReq, err := l.ToBytes(req, make([]byte, 0, 1024))
	if nil != err {
		return err
	}

	if err = l.sendPacket(bsReq); err != nil {
		return err
	}

	var bsResp []byte
	if bsResp, err = l.recvPacket(); err != nil {
		return err
	}

	return l.FromBytes(resp, bsResp)
}

func (l *lanPlus) exec(cmd commands.CommandCode, reqData, respData interface{}) error {
	var req Request
	var resp Response

	req.Init(cmd, reqData)
	resp.Init(cmd, respData)

	return l.send(&req, &resp)
}

func (l *lanPlus) sendv1(req *Request, res *Response) error {
	var sequence uint32
	var rqSequence uint8

	if req.NetFn() == commands.NetworkFunctionApp && req.Body.Cmd == commands.GetChannelAuthenticationCapabilities.Code {
		sequence = 0
		rqSequence = 0
	} else {
		sequence = l.nextSequence()
		rqSequence = l.nextRqSequence()
	}

	bsReq, err := IPMIv1ToBytes(req, make([]byte, 0, 1024), sequence, rqSequence)
	if nil != err {
		return err
	}

	if err = l.sendPacket(bsReq); err != nil {
		return err
	}

	var bsResp []byte
	if bsResp, err = l.recvPacket(); err != nil {
		return err
	}

	return IPMIv1FromBytes(res, bsResp)
}

func (l *lanPlus) execv1(cmd commands.CommandCode, reqData, respData interface{}) error {
	var req Request
	var resp Response

	req.Init(cmd, reqData)
	resp.Init(cmd, respData)

	return l.sendv1(&req, &resp)
}

func (l *lanPlus) openSession() error {
	// if err := l.ping(); err != nil {
	// 	return err
	// }

	l.sequence = 2
	l.rqSeqence = 0

	if err := l.getAuthCapabilities(true); err != nil {
		return err
	}
	if err := l.rmcpOpen(); err != nil {
		return err
	}
	l.active = true

	if err := l.rakp1(); err != nil {
		return err
	}
	if err := l.rakp3(); err != nil {
		return err
	}
	return l.setPrivilegeLevel()
}

func (l *lanPlus) getAuthCapabilities(isV2 bool) error {
	var req AuthCapabilitiesRequest
	var resp AuthCapabilitiesResponse

	if isV2 {
		req.ChannelNumber = 0x8e // lanChannelE + Version compatibility: IPMI v2.0+ extended data (1)
	} else {
		req.ChannelNumber = 0x0e // lanChannelE
	}
	req.PrivLevel = l.PrivLevel

	if err := l.execv1(commands.GetChannelAuthenticationCapabilities, &req, &resp); err != nil {
		return err
	}

	for _, t := range []uint8{AuthTypeMD5, AuthTypePassword} {
		if (resp.AuthTypeSupport & (1 << t)) != 0 {
			//l.AuthType = t
			return nil
		}
	}

	if (resp.Reserved & (1 << 1)) == 0 { // is IPMI v2
		return errors.New("IPMI version is unsupported v2.0")
	}

	if (resp.AuthTypeSupport & 0x7f) == AuthTypeNone {
		//l.AuthType = AuthTypeNone
		return nil
	}

	//log.Printf("BMC did not offer a supported AuthType")
	return errors.New("BMC did not offer a supported AuthType(" + strconv.FormatInt(int64(resp.AuthTypeSupport), 10) + ")")
}

func (l *lanPlus) rmcpOpen() error {
	var req = OpenSessionRequest{
		PrivLevel: l.PrivLevel,
		SessionID: nextSessionID(),
	}

	var resp OpenSessionResponse

	req.Authentication.Algorithm = l.authenticationAlgorithm
	req.Integrity.Algorithm = l.integrityAlgorithm
	req.Confidentiality.Algorithm = l.confidentialityAlgorithm
	err := l.send(&req, &resp)
	if nil != err {
		return err
	}

	if resp.StatusCode != uint8(RAKP_STATUS_NO_ERRORS) {
		return StatusCode(resp.StatusCode)
	}

	l.PrivLevel = resp.PrivLevel
	l.LSessionID = resp.SessionID
	l.RSessionID = resp.MSessionID

	// fmt.Println("AuthenticationAlgorithm =", l.authenticationAlgorithm,
	// 	", IntegrityAlgorithm =", l.integrityAlgorithm,
	// 	", ConfidentialityAlgorithm =", l.confidentialityAlgorithm)

	l.authenticationAlgorithm = resp.Authentication.Algorithm
	l.integrityAlgorithm = resp.Integrity.Algorithm
	l.confidentialityAlgorithm = resp.Confidentiality.Algorithm

	// fmt.Println("AuthenticationAlgorithm =", resp.Authentication.Algorithm,
	// 	", IntegrityAlgorithm =", resp.Integrity.Algorithm,
	// 	", ConfidentialityAlgorithm =", resp.Confidentiality.Algorithm)

	return nil
}

func (l *lanPlus) rakp1() error {
	privlevel := l.PrivLevel
	if l.RoleUserOnlyLookup {
		privlevel |= (1 << 4)
	}

	var req = RakpMessage1{SessionID: l.RSessionID,
		PrivLevel: privlevel,
	}
	var resp = RakpMessage2{}

	// var bs [4]byte
	// binary.LittleEndian.PutUint32(bs[:], l.RSessionID)
	// fmt.Println(hex.EncodeToString(bs[:]))
	// binary.LittleEndian.PutUint32(bs[:], l.LSessionID)
	// fmt.Println(hex.EncodeToString(bs[:]))

	readRand(req.Rand[:])
	copy(l.LRand[:], req.Rand[:])
	copy(req.UserName[:], l.Username[:l.UsernameLen])
	req.UserNameLength = l.UsernameLen

	err := l.send(&req, &resp)
	if nil != err {
		return err
	}
	if resp.StatusCode != uint8(RAKP_STATUS_NO_ERRORS) {
		return StatusCode(resp.StatusCode)
	}

	copy(l.RRand[:], resp.Rand[:])
	copy(l.RGuid[:], resp.Guid[:])

	var w Writer
	w.Init(make([]byte, 0, 128))
	w.WriteUint32(l.LSessionID) // local id
	w.WriteUint32(l.RSessionID) // remote id
	w.WriteBytes(l.LRand[:])
	w.WriteBytes(l.RRand[:])
	w.WriteBytes(l.RGuid[:])
	w.WriteUint8(privlevel)
	w.WriteUint8(l.UsernameLen)
	w.WriteBytes(l.Username[:int(l.UsernameLen)])

	res := l.authenticationAlgorithm.Gen(l.Authcode[:], w.Bytes())

	if !bytes.Equal(res, resp.KeyExchange) {
		return ErrPasswordNotMatch
	}

	w.Reset()
	w.WriteBytes(l.LRand[:])
	w.WriteBytes(l.RRand[:])
	w.WriteUint8(privlevel) // privlevel
	w.WriteUint8(l.UsernameLen)
	w.WriteBytes(l.Username[:int(l.UsernameLen)])

	l.sik = l.authenticationAlgorithm.Gen(l.Authcode[:], w.Bytes())
	fmt.Printf("sik=%x\r\n", l.sik)

	signSize := l.integrityAlgorithm.SignSize()
	l.k1 = l.integrityAlgorithm.Gen(l.sik, bytes.Repeat([]byte{1}, signSize))[:signSize]
	fmt.Printf("k1=%x\r\n", l.k1)

	l.k2 = l.integrityAlgorithm.Gen(l.sik, bytes.Repeat([]byte{2}, signSize))[:signSize]
	fmt.Printf("k2=%x\r\n", l.k2)

	return nil
}

func (l *lanPlus) rakp3() error {
	var req = RakpMessage3{SessionID: l.RSessionID,
		StatusCode: uint8(RAKP_STATUS_NO_ERRORS),
		//KeyExchange []byte // (9-*) N 取决于RakpMessage1 AuthAlg
	}
	var resp = RakpMessage4{}

	privlevel := l.PrivLevel
	if l.RoleUserOnlyLookup {
		privlevel |= (1 << 4)
	}

	var w Writer
	w.Init(make([]byte, 0, 64))
	w.WriteBytes(l.RRand[:])
	w.WriteUint32(l.LSessionID) // remote id
	w.WriteUint8(privlevel)     // privlevel
	w.WriteUint8(l.UsernameLen)
	w.WriteBytes(l.Username[:int(l.UsernameLen)])
	res := l.authenticationAlgorithm.Gen(l.Authcode[:], w.Bytes())
	req.KeyExchange = res

	err := l.send(&req, &resp)
	if nil != err {
		return err
	}

	if resp.StatusCode != uint8(RAKP_STATUS_NO_ERRORS) {
		return StatusCode(resp.StatusCode)
	}

	w.Reset()
	w.WriteBytes(l.LRand[:])
	w.WriteUint32(l.RSessionID)
	w.WriteBytes(l.RGuid[:])

	//fmt.Println(hex.EncodeToString(w.Bytes()))
	res = l.authenticationAlgorithm.Gen(l.sik, w.Bytes())
	//fmt.Println(hex.EncodeToString(res))

	if bytes.Equal(res, resp.IntegrityCheck) ||
		(len(resp.IntegrityCheck) < len(res) && bytes.HasPrefix(res, resp.IntegrityCheck)) {
		//l.sik = sik
		return nil
	}

	//fmt.Println("excepted is", hex.EncodeToString(res))
	//fmt.Println("actual   is", hex.EncodeToString(resp.IntegrityCheck))
	return errors.New("integrity check is failed")
}

func (l *lanPlus) setPrivilegeLevel() error {
	var req_payload = SessionPrivilegeLevelRequest{PrivLevel: l.PrivLevel}
	var resp_payload = SessionPrivilegeLevelResponse{}
	return l.exec(commands.SetSessionPrivilegeLevel, &req_payload, &resp_payload)
}

func (l *lanPlus) closeSession() error {
	var req CloseSessionRequest
	var resp CloseSessionResponse
	req.SessionID = l.RSessionID
	return l.exec(commands.CloseSession, &req, &resp)
}

func newLanPlus(opt *ConnectionOption) *lanPlus {
	l := &lanPlus{lanBase: lanBase{conn_opt: opt, sequence: 2}}

	l.RoleUserOnlyLookup = true
	l.UsernameLen = uint8(len(opt.Username))
	copy(l.Username[:], opt.Username[:])
	copy(l.Authcode[:], opt.Password[:])

	l.sequence = 2

	if opt.PrivLevel == commands.PrivLevelNone {
		l.PrivLevel = uint8(commands.PrivLevelAdmin) // Admini
	} else {
		l.PrivLevel = uint8(opt.PrivLevel)
	}

	l.authenticationAlgorithm = opt.AuthenticationAlgorithm
	l.integrityAlgorithm = opt.IntegrityAlgorithm
	l.confidentialityAlgorithm = opt.ConfidentialityAlgorithm

	if l.authenticationAlgorithm == RAKPAlgorithmAuth_None {
		l.authenticationAlgorithm = RAKPAlgorithmAuth_HMAC_SHA1
	}
	// l.authenticationAlgorithm = RAKPAlgorithmAuth_HMAC_SHA1
	// l.integrityAlgorithm = RAKPAlgorithmIntegrity_None
	// l.confidentialityAlgorithm = RAKPAlgorithmEncryto_None
	return l
}
