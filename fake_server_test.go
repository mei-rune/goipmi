package goipmi

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/runner-mei/goipmi/protocol"
	"github.com/runner-mei/goipmi/protocol/commands"
)

type fakeServer struct {
	datas     [][]byte
	wait      sync.WaitGroup
	isStopped int32
	LocalPort int

	closer io.Closer

	buffer bytes.Buffer
}

func (fs *fakeServer) LogBuffer() *bytes.Buffer {
	return &fs.buffer
}

func (fs *fakeServer) Close() error {
	var err error
	if atomic.CompareAndSwapInt32(&fs.isStopped, 0, 1) {

		if fs.closer != nil {
			err = fs.closer.Close()
		}

		fs.wait.Wait()
	}
	return err
}

func newFakeServer(file string) (*fakeServer, error) {
	fs := &fakeServer{}

	handle, e := pcap.OpenOffline(file)
	if nil != e {
		return nil, e
	}
	defer handle.Close()

	if e = handlePacketData(handle, func(bs []byte) error {
		copyed := make([]byte, len(bs))
		copy(copyed, bs)
		fs.datas = append(fs.datas, copyed)
		return nil
	}); e != nil {
		return nil, e
	}

	conn, e := net.ListenUDP("udp", &net.UDPAddr{Port: 0})
	if nil != e {
		return nil, e
	}

	go func() {
		count := 0
		defer fs.wait.Done()

		bs := make([]byte, 1024)
		for 0 == atomic.LoadInt32(&fs.isStopped) {
			n, raddr, e := conn.ReadFrom(bs)
			if nil != e {
				if 0 == atomic.LoadInt32(&fs.isStopped) {
					fmt.Println(e)
				}
				break
			}

			fmt.Println("===============", count)
			if !bytes.Equal(bs[:n], fs.datas[count]) {
				fmt.Fprintf(&fs.buffer, "[%d] excepted %x \r\n", count, fs.datas[count])
				fmt.Fprintf(&fs.buffer, "[%d] got      %x \r\n", count, bs[:n])
			}
			fmt.Printf("[%d] excepted %x \r\n", count, fs.datas[count])
			fmt.Printf("[%d] got      %x \r\n", count, bs[:n])
			count++

			if _, e = conn.WriteTo(fs.datas[count], raddr); e != nil {
				fmt.Println(e)
				break
			}
			count++
		}
	}()
	fs.wait.Add(1)

	_, sport, _ := net.SplitHostPort(conn.LocalAddr().String())
	port, _ := strconv.Atoi(sport)
	fs.LocalPort = port

	fmt.Println("fakeServer listen at", sport, "with data length =", len(fs.datas))

	fs.closer = conn
	return fs, nil
}

func handlePacketData(handle *pcap.Handle, handleData func(bs []byte) error) error {
	var eth layers.Ethernet
	var ip4 layers.IPv4
	var ip6 layers.IPv6
	var tcp layers.TCP
	var udp layers.UDP
	var payload gopacket.Payload

	parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &eth, &ip4, &ip6, &tcp, &udp, &payload)
	decoded := make([]gopacket.LayerType, 0, 10)

	for {
		packetData, _, e := handle.ZeroCopyReadPacketData()
		if nil != e {
			if io.EOF == e {
				return nil
			}
			return e
		}
		if e := parser.DecodeLayers(packetData, &decoded); nil != e {
			return nil
			//log.Println("[warn]", e)
			//continue
		}

		is_payload := false

		for _, layerType := range decoded {
			switch layerType {
			case layers.LayerTypeEthernetCTP:
			case layers.LayerTypeCiscoDiscovery:
			case layers.LayerTypeDot1Q:
			case layers.LayerTypeMPLS:
			case layers.LayerTypeLLC:
			case layers.LayerTypeLoopback:
			case layers.LayerTypeLinkLayerDiscovery:
			case layers.LayerTypeICMPv6:
			case layers.LayerTypeARP:
			case layers.LayerTypeDNS:
			case layers.LayerTypeEthernet:
			case layers.LayerTypeIPv6:
			case layers.LayerTypeIPv4:
			case layers.LayerTypeTCP:
			case layers.LayerTypeUDP:
				is_payload = true
			case layers.LayerTypeSCTP:
			case gopacket.LayerTypePayload:
			default:
				fmt.Println("unknow layer type -", layerType.String())
			}
		}

		if is_payload {
			// fmt.Println("=================== begin ===================")

			e := handleData(payload.Payload())
			if nil != e {
				fmt.Println(e)
				//return e
			}

			//fmt.Println("===================  end  ===================")
		}
	}
}

func TestAuthHmacsha1IntHmacsha196(t *testing.T) {
	fs, err := newFakeServer("authenticated.pcap")
	if err != nil {
		t.Error(err)
		return
	}

	defer fs.Close()

	rand, err := hex.DecodeString("67bdcd4297570bf0dccaeeeb1c3e105f")
	if err != nil {
		t.Error(err)
		return
	}
	SetFake(rand, 2695013284)

	client, err := NewClient(&ConnectionOption{
		Hostname:                "127.0.0.1",
		Port:                    fs.LocalPort,
		Username:                "Administrator",
		Password:                "123456abc",
		Interface:               "lanplus",
		PrivLevel:               commands.PrivLevelAdmin,
		AuthenticationAlgorithm: protocol.RAKPAlgorithmAuth_HMAC_SHA1,
		IntegrityAlgorithm:      protocol.RAKPAlgorithmIntegrity_HMAC_SHA1_96,
	})
	if err != nil {
		t.Error(err)
		return
	}

	if err := client.Open(); err != nil {
		t.Error(err)
		return
	}

	resp, err := client.GetDeviceID()
	if err != nil {
		t.Error(err)
		return
	}

	fs.Close()

	if fs.LogBuffer().Len() > 0 {
		t.Error("\r\n" + fs.LogBuffer().String())
	} else {
		t.Log(resp)
	}
}

func TestEncryptoAESCBC(t *testing.T) {
	fs, err := newFakeServer("encrypto_aes_cbc.pcap")
	if err != nil {
		t.Error(err)
		return
	}

	defer fs.Close()

	rand, err := hex.DecodeString("c6fc46aff2582fa5c803dd3b5d68a0c5")
	if err != nil {
		t.Error(err)
		return
	}
	SetFake(rand, 2695013284)

	client, err := NewClient(&ConnectionOption{
		Hostname:                 "127.0.0.1",
		Port:                     fs.LocalPort,
		Username:                 "Administrator",
		Password:                 "123456abc",
		Interface:                "lanplus",
		PrivLevel:                commands.PrivLevelAdmin,
		AuthenticationAlgorithm:  protocol.RAKPAlgorithmAuth_HMAC_SHA1,
		IntegrityAlgorithm:       protocol.RAKPAlgorithmIntegrity_HMAC_SHA1_96,
		ConfidentialityAlgorithm: protocol.RAKPAlgorithmEncryto_AES_CBC_128,
	})
	if err != nil {
		t.Error(err)
		return
	}

	protocol.SetInitializationVector([]byte{0x5f, 0x34, 0x43, 0xa2, 0xcb, 0x24, 0x5d, 0xd4, 0xf5, 0x6c, 0x85, 0x66, 0x8c, 0xa7, 0x57, 0x14})
	if err := client.Open(); err != nil {
		t.Error(err)
		return
	}

	protocol.SetInitializationVector([]byte{0xb4, 0xfd, 0x3a, 0x95, 0x92, 0x6a, 0xd8, 0x8f, 0x85, 0x0f, 0xf0, 0x0d, 0x1e, 0xd2, 0x5f, 0xa6})
	resp, err := client.GetDeviceID()
	if err != nil {
		t.Error(err)
		t.Error("\r\n" + fs.LogBuffer().String())
		return
	}

	fs.Close()

	if fs.LogBuffer().Len() > 0 {
		t.Error("\r\n" + fs.LogBuffer().String())
	} else {
		t.Log(resp)
	}
}
