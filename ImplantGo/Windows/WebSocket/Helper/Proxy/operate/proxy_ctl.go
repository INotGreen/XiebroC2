package operate

import (
	"errors"

	"main/Helper/Proxy/option"
	"net"
	"time"

	"github.com/xtaci/smux"
)

const (
	CTL_HANDSHAKE = iota
	CTL_CONNECT_ME
	CTL_CLEANUP

	MAX_CONNECTION   = 0x800
	CLIENT_HANDSHAKE = 0xC0
	SERVER_HANDSHAKE = 0xE0
)

type Protocol struct {
	CMD byte
	N   byte

	// ACK uint16
}

var PROTO_END = []byte{0xEE, 0xFF}

func marshal(p Protocol) []byte {
	buf := make([]byte, 4)
	buf[0] = p.CMD
	buf[1] = p.N

	buf[2], buf[3] = PROTO_END[0], PROTO_END[1]
	return buf
}

func unmarshal(b []byte) Protocol {
	return Protocol{
		CMD: b[0],
		N:   b[1],
	}
}

func bytesEq(a, b []byte) bool {
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func readUntilEnd(conn net.Conn) ([]byte, error) {
	buf := make([]byte, 1)
	output := make([]byte, 0, 4)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			return nil, err
		}

		if n != 1 || len(output) > 4 {
			return nil, errors.New("Transmission error")
		}

		output = append(output, buf[0])

		if len(output) == 4 && bytesEq(PROTO_END, output[len(output)-2:]) {
			break
		}
	}

	return output[:2], nil
}

func serverHandshake(listener net.Listener) (*smux.Session, *smux.Stream, error) {
	var conn net.Conn
	var session *smux.Session
	var ctlStream *smux.Stream
	var err error

	for {
		conn, err = listener.Accept()
		if err != nil {
			continue
		}

		session, err = smux.Server(conn, &smux.Config{
			Version:           2,
			KeepAliveInterval: option.SMUX_KEEPALIVE_INTERVAL * time.Second,
			KeepAliveTimeout:  option.SMUX_KEEPALIVE_TIMEOUT * time.Second,
			MaxFrameSize:      option.SMUX_FRAMESIZE,
			MaxReceiveBuffer:  option.SMUX_RECVBUFFER,
			MaxStreamBuffer:   option.SMUX_STREAMBUFFER,
		})
		if err != nil {
			return nil, nil, err
		}

		ctlStream, err = session.AcceptStream()
		if err != nil {
			return nil, nil, err
		}

		pb, err := readUntilEnd(ctlStream)
		if err != nil {
			continue
		}

		p := unmarshal(pb)
		if p.CMD == CTL_HANDSHAKE && p.N == CLIENT_HANDSHAKE {
			ctlStream.Write(marshal(Protocol{
				CMD: CTL_HANDSHAKE,
				N:   SERVER_HANDSHAKE,
			}))
			break
		}
	}

	return session, ctlStream, nil
}

func clientHandshake(remote string) (*smux.Session, *smux.Stream, error) {
	conn, err := net.DialTimeout(
		"tcp", remote,
		time.Millisecond*time.Duration(option.TIMEOUT),
	)
	if err != nil {
		return nil, nil, err
	}

	session, err := smux.Client(conn, &smux.Config{
		Version:           2,
		KeepAliveInterval: option.SMUX_KEEPALIVE_INTERVAL * time.Second,
		KeepAliveTimeout:  option.SMUX_KEEPALIVE_TIMEOUT * time.Second,
		MaxFrameSize:      option.SMUX_FRAMESIZE,
		MaxReceiveBuffer:  option.SMUX_RECVBUFFER,
		MaxStreamBuffer:   option.SMUX_STREAMBUFFER,
	})
	if err != nil {
		return nil, nil, err
	}

	ctlStream, err := session.OpenStream()
	if err != nil {
		return nil, nil, err
	}

	ctlStream.Write(marshal(Protocol{
		CMD: CTL_HANDSHAKE,
		N:   CLIENT_HANDSHAKE,
	}))

	pb, err := readUntilEnd(ctlStream)
	if err != nil {
		return nil, nil, errors.New("Connect to remote forward server error")
	}

	p := unmarshal(pb)
	if !(p.CMD == CTL_HANDSHAKE && p.N == SERVER_HANDSHAKE) {
		return nil, nil, errors.New("Connect to remote forward server error")
	}

	return session, ctlStream, nil
}
