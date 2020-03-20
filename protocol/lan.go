package protocol

import (
	"errors"
	"log"
	"strconv"

	"github.com/runner-mei/goipmi/protocol/commands"
)

type lan struct {
	lanBase

	active bool

	PrivLevel uint8
	AuthType  uint8
	Authcode  [16]uint8
	Username  [16]uint8
	SessionID uint32
}

func (l *lan) FromBytes(resp *Response, bs []byte) error {
	return IPMIv1FromBytes(resp, bs)
}

func IPMIv1FromBytes(resp *Response, bs []byte) error {
	var rmcpHeader RMCPHeader
	var ipmiHeader IPMIV1Header

	var r Reader
	r.Init(bs)

	rmcpHeader.ReadBytes(&r)
	ipmiHeader.ReadBytes(&r)
	resp.ReadBytes(&r)

	return r.Err()
}

func (l *lan) ToBytes(req *Request, bs []byte) ([]byte, error) {
	if req.NetFn() == commands.NetworkFunctionApp && req.Body.Cmd == commands.GetChannelAuthenticationCapabilities.Code {
		return IPMIv1ToBytes(req, bs, 0, l.nextRqSequence())
	}
	return IPMIv1ToBytes(req, bs, l.nextSequence(), l.nextRqSequence())
}

func IPMIv1ToBytes(req *Request, bs []byte, ipmiSeq uint32, rqSeq uint8) ([]byte, error) {
	rmcpHeader := RMCPHeader{
		Version:            rmcpVersion1,
		Class:              rmcpClassIPMI,
		RMCPSequenceNumber: 0xff,
	}
	ipmiHeader := IPMIV1Header{
		Sequence: ipmiSeq,
		//SessionID: l.SessionID, //l.SessionID,
	}
	req.Body.RqSeq = rqSeq

	w := Writer{}
	w.Init(bs)
	rmcpHeader.WriteBytes(&w)
	ipmiHeader.WriteBytes(&w)
	if nil != w.Err() {
		return nil, w.Err()
	}
	old_length := w.Len()
	req.WriteBytes(&w)
	if nil != w.Err() {
		return nil, w.Err()
	}
	tmp := w.Bytes()
	tmp[old_length-1] = uint8(w.Len() - old_length)

	return tmp, nil
}

func (l *lan) open() error {
	err := l.dial()
	if err != nil {
		return err
	}

	// // TODO: options
	// l.priv = PrivLevelAdmin
	// l.timeout = time.Second * 5
	// l.lun = 0
	return l.openSession()
}

func (l *lan) close() error {
	if l.active {
		err := l.closeSession()
		if err != nil {
			log.Printf("error closing session: %s", err)
		}
		l.active = false
	}

	l.disconnect()
	return nil
}

func (l *lan) send(req, resp interface{}) error {
	return l.sendExt(req.(*Request), resp.(*Response))
}

func (l *lan) sendExt(req *Request, resp *Response) error {
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

func (l *lan) exec(cmd commands.CommandCode, reqData, respData interface{}) error {
	var req Request
	var resp Response

	req.Init(cmd, reqData)
	resp.Init(cmd, respData)

	return l.sendExt(&req, &resp)
}

func (l *lan) openSession() error {
	if err := l.ping(); err != nil {
		return err
	}

	if err := l.getAuthCapabilities(false); err != nil {
		return err
	}

	res, err := l.getSessionChallenge()
	if err != nil {
		return err
	}

	if err := l.activateSession(res); err != nil {
		return err
	}

	return l.setSessionPriv()
}

func (l *lan) getAuthCapabilities(isV2 bool) error {
	var req AuthCapabilitiesRequest
	var resp AuthCapabilitiesResponse

	if isV2 {
		req.ChannelNumber = 0x8e // lanChannelE + Version compatibility: IPMI v2.0+ extended data (1)
	} else {
		req.ChannelNumber = 0x0e // lanChannelE
	}
	req.PrivLevel = l.PrivLevel

	if err := l.exec(commands.GetChannelAuthenticationCapabilities, &req, &resp); err != nil {
		return err
	}

	for _, t := range []uint8{AuthTypeMD5, AuthTypePassword} {
		if (resp.AuthTypeSupport & (1 << t)) != 0 {
			l.AuthType = t
			return nil
		}
	}

	if (resp.Reserved & (1 << 1)) != 0 { // is IPMI v2
		return errors.New("IPMI version is unsupported v1.5")
	}

	if (resp.AuthTypeSupport & 0x7f) == AuthTypeNone {
		l.AuthType = AuthTypeNone
		return nil
	}

	//log.Printf("BMC did not offer a supported AuthType")
	return errors.New("BMC did not offer a supported AuthType(" + strconv.FormatInt(int64(resp.AuthTypeSupport), 10) + ")")
}

func (l *lan) getSessionChallenge() (*SessionChallengeResponse, error) {
	var req SessionChallengeRequest
	var resp SessionChallengeResponse

	req.AuthType = l.AuthType
	req.Username = l.Username

	if err := l.exec(commands.GetSessionChallenge, req, resp); err != nil {
		return nil, err
	}

	l.SessionID = resp.TemporarySessionID
	return &resp, nil
}

func (l *lan) activateSession(sc *SessionChallengeResponse) error {
	req := &ActivateSessionRequest{
		AuthType:  l.AuthType,
		PrivLevel: l.PrivLevel,
		AuthCode:  sc.Challenge,
		InSeq:     l.inSeq(),
	}
	resp := &ActivateSessionResponse{}

	if err := l.exec(commands.ActivateSession, req, resp); err != nil {
		return err
	}

	l.active = true
	l.SessionID = resp.SessionID
	l.AuthType = resp.AuthType
	l.sequence = resp.InboundSeq
	return nil
}

func (l *lan) setSessionPriv() error {
	var req SessionPrivilegeLevelRequest
	var resp SessionPrivilegeLevelResponse
	req.PrivLevel = l.PrivLevel

	if err := l.exec(commands.SetSessionPrivilegeLevel, req, resp); err != nil {
		return err
	}

	l.PrivLevel = resp.NewPrivilegeLevel
	return nil
}

func (l *lan) closeSession() error {
	var req CloseSessionRequest
	var resp CloseSessionResponse
	req.SessionID = l.SessionID
	return l.exec(commands.SetSessionPrivilegeLevel, req, resp)
}

func newLan(opt *ConnectionOption) *lan {
	l := &lan{lanBase: lanBase{conn_opt: opt}}

	copy(l.Username[:], opt.Username[:])
	copy(l.Authcode[:], opt.Password[:])
	return l
}
