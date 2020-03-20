package protocol

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/runner-mei/goipmi/protocol/commands"
)

// Connection properties for a Client
type ConnectionOption struct {
	Hostname  string
	Port      int
	Username  string
	Password  string
	Interface string

	PrivLevel                commands.PrivLevelType
	AuthenticationAlgorithm  RAKPAlgorithmAuthMethod
	IntegrityAlgorithm       RAKPAlgorithmIntegrityMethod
	ConfidentialityAlgorithm RAKPAlgorithmConfidentialityMethod
}

// RemoteIP returns the remote (bmc) IP address of the Connection
func (c *ConnectionOption) RemoteIP() string {
	if net.ParseIP(c.Hostname) == nil {
		addrs, err := net.LookupHost(c.Hostname)
		if err != nil && len(addrs) > 0 {
			return addrs[0]
		}
	}
	return c.Hostname
}

// LocalIP returns the local (client) IP address of the Connection
func (c *ConnectionOption) LocalIP() string {
	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", c.Hostname, c.Port))
	if err != nil {
		// don't bother returning an error, since this value will never
		// make it to the bmc if we can't connect to it.
		return c.Hostname
	}
	_ = conn.Close()
	host, _, _ := net.SplitHostPort(conn.LocalAddr().String())
	return host
}

type lanBase struct {
	conn_opt *ConnectionOption
	conn     *net.UDPConn
	//addr     net.Addr

	sequence  uint32
	rqSeqence uint8
}

func (l *lanBase) isConnected() bool {
	return l.conn != nil
}

func (l *lanBase) inSeq() [4]uint8 {
	seq := [4]uint8{}
	if _, err := rand.Read(seq[:]); err != nil {
		panic(err)
	}
	return seq
}

func (l *lanBase) nextSequence() uint32 {
	l.sequence++
	return l.sequence
}

func (l *lanBase) nextRqSequence() uint8 {
	l.rqSeqence++
	return l.rqSeqence << 2
}

func (l *lanBase) dial() error {
	addr, err := net.ResolveUDPAddr("udp4",
		net.JoinHostPort(l.conn_opt.Hostname, strconv.Itoa(l.conn_opt.Port)))
	if err != nil {
		return err
	}
	//l.addr = addr

	conn, err := net.DialUDP("udp4",
		nil, //net.JoinHostPort(l.conn_opt.Hostname, strconv.Itoa(l.conn_opt.Port)),
		addr)
	if err != nil {
		return err
	}
	l.conn = conn
	return nil
}

func (l *lanBase) disconnect() {
	if l.conn != nil {
		_ = l.conn.Close()
		l.conn = nil
	}
}

func (l *lanBase) sendPacket(buf []byte) error {
	//fmt.Println(l.addr)
	_, err := l.conn.Write(buf)
	return err
}

func (l *lanBase) isTargetAddr(addr net.Addr) bool {
	return true // addr.String() == l.addr.String()
}

func (l *lanBase) recvPacket() ([]byte, error) {
	buf := make([]byte, ipmiBufSize)

	for i := 0; ; i++ {
		if i >= 3 {
			return nil, ErrTimeout
		}
		err := l.conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		if err != nil {
			return nil, err
		}

		n, addr, err := l.conn.ReadFrom(buf)
		if err != nil {
			if to, ok := err.(interface {
				Timeout() bool
			}); ok {
				if to.Timeout() {
					return nil, ErrTimeout
				}
			}
			return nil, err
		}

		if l.isTargetAddr(addr) {
			return buf[:n], nil
		}
	}
}

func (l *lanBase) ping() error {
	req := &asfMessage{
		RMCP: RMCPHeader{
			Version:            rmcpVersion1,
			Class:              rmcpClassASF,
			RMCPSequenceNumber: 0xff,
		},
		ASF: asfHeader{
			IANAEnterpriseNumber: asfIANA,
			MessageType:          asfMessageTypePing,
		},
	}

	if bs, err := ToBytes(req); err != nil {
		return err
	} else if err := l.sendPacket(bs); err != nil {
		return err
	}

	buf, err := l.recvPacket()
	if err != nil {
		return err
	}

	var pong asfPong
	var resp = asfMessage{Data: &pong}

	if err := FromBytes(&resp, buf); err != nil {
		return err
	}
	if resp.ASF.MessageType != asfMessageTypePong {
		return fmt.Errorf("unsupported ASF message type: %d", resp.ASF.MessageType)
	}

	if !pong.valid() {
		return errors.New("IPMI not supported")
	}

	return nil
}

var _session_id uint32 = 100

func nextSessionID() uint32 {
	if is_test {
		return fakeSessionID
	}
	return atomic.AddUint32(&_session_id, 1)
}
